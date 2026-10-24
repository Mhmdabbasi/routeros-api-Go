package routeros

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Client is the main RouterOS API client
type Client struct {
	config     *Config
	conn       net.Conn
	protocol   *Protocol
	connected  atomic.Bool
	closed     atomic.Bool
	tagCounter atomic.Uint32

	// Synchronization
	mu       sync.RWMutex
	pending  map[uint32]chan *Response
	notifyCh chan *Response

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc

	// Response handlers
	responseHandler ResponseHandler
	errorHandler    ErrorHandler
}

// ResponseHandler is a callback for handling responses
type ResponseHandler func(*Response)

// ErrorHandler is a callback for handling errors
type ErrorHandler func(error)

// NewClient creates a new RouterOS API client with the given configuration
func NewClient(config *Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		config:   config,
		pending:  make(map[uint32]chan *Response),
		notifyCh: make(chan *Response, config.QueueSize),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Connect establishes a connection to the RouterOS device
func (c *Client) Connect() error {
	if c.connected.Load() {
		return ErrAlreadyConnected
	}

	// Validate configuration
	if err := c.config.Validate(); err != nil {
		return err
	}

	c.config.Logger.Debug("Connecting to %s:%d", c.config.Address, c.config.Port)

	// Build dialer options
	dialer := &net.Dialer{
		Timeout:   c.config.Timeout,
		KeepAlive: c.config.KeepAliveTime,
	}

	if c.config.KeepAlive {
		dialer.KeepAlive = c.config.KeepAliveTime
	}

	if c.config.SourceAddress != "" {
		dialer.LocalAddr = &net.TCPAddr{
			IP: net.ParseIP(c.config.SourceAddress),
		}
	}

	var err error
	var conn net.Conn

	// Establish connection with retries
	for attempt := 0; attempt <= c.config.RetryCount; attempt++ {
		if c.config.TLS {
			conn, err = c.connectTLS(dialer)
		} else {
			addr := c.config.NetworkAddress()
			conn, err = dialer.DialContext(c.ctx, "tcp", addr)
		}

		if err == nil {
			break
		}

		if attempt < c.config.RetryCount {
			c.config.Logger.Debug("Connection attempt %d failed: %v, retrying in %v", attempt+1, err, c.config.RetryDelay)
			time.Sleep(c.config.RetryDelay)
		}
	}

	if err != nil {
		return NewConnectionError(c.config.Address, c.config.Port, err)
	}

	c.conn = conn
	c.protocol = NewProtocol(conn, conn)

	// Set buffer sizes if specified
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		if c.config.BufferSize > 0 {
			tcpConn.SetReadBuffer(c.config.BufferSize)
			tcpConn.SetWriteBuffer(c.config.BufferSize)
		}
	}

	// Perform login
	if err := c.login(); err != nil {
		conn.Close()
		return err
	}

	c.connected.Store(true)
	c.closed.Store(false)

	// Start response reader goroutine
	go c.readResponses()

	c.config.Logger.Info("Connected to %s:%d", c.config.Address, c.config.Port)
	return nil
}

// connectTLS establishes a TLS connection
func (c *Client) connectTLS(dialer *net.Dialer) (net.Conn, error) {
	tlsConfig, err := c.config.BuildTLSConfig()
	if err != nil {
		return nil, err
	}

	addr := c.config.NetworkAddress()
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	if err != nil {
		return nil, err
	}

	// Verify connection state
	state := conn.ConnectionState()
	if !state.HandshakeComplete {
		return nil, ErrTLSError
	}

	return conn, nil
}

// login performs the authentication sequence
func (c *Client) login() error {
	c.config.Logger.Debug("Starting login sequence")

	// Read the first response to get the challenge
	sentence, err := c.protocol.ReadSentence()
	if err != nil {
		return fmt.Errorf("failed to read login challenge: %w", err)
	}

	response := ParseSentence(sentence)
	if response.Type != ResponseDone {
		return fmt.Errorf("unexpected response during login: %v", response.Type)
	}

	// Extract the challenge (ret field)
	ret, ok := response.Data["ret"]
	if !ok {
		// RouterOS v6.43+ uses plain text login
		return c.loginPlain()
	}

	// RouterOS pre-6.43 uses MD5 challenge-response
	return c.loginMD5(ret)
}

// loginPlain performs plain text authentication (RouterOS v6.43+)
func (c *Client) loginPlain() error {
	c.config.Logger.Debug("Using plain text authentication")

	err := c.protocol.WriteSentence("/login", fmt.Sprintf("=name=%s", c.config.Username), fmt.Sprintf("=password=%s", c.config.Password))
	if err != nil {
		return fmt.Errorf("failed to send login credentials: %w", err)
	}

	sentence, err := c.protocol.ReadSentence()
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	response := ParseSentence(sentence)
	if response.Type == ResponseTrap {
		return ErrAuthenticationFailed
	}

	if response.Type != ResponseDone {
		return fmt.Errorf("unexpected response during login: %v", response.Type)
	}

	c.config.Logger.Debug("Plain text authentication successful")
	return nil
}

// loginMD5 performs MD5 challenge-response authentication (RouterOS pre-6.43)
func (c *Client) loginMD5(challenge string) error {
	c.config.Logger.Debug("Using MD5 challenge-response authentication")

	// Decode the challenge (hex string)
	challengeBytes, err := hex.DecodeString(challenge)
	if err != nil {
		return fmt.Errorf("failed to decode challenge: %w", err)
	}

	// Calculate MD5 response: MD5(0x00 + password + challenge)
	hasher := md5.New()
	hasher.Write([]byte{0x00})
	hasher.Write([]byte(c.config.Password))
	hasher.Write(challengeBytes)
	response := hex.EncodeToString(hasher.Sum(nil))

	// Send login response
	err = c.protocol.WriteSentence("/login", fmt.Sprintf("=name=%s", c.config.Username), fmt.Sprintf("=response=00%s", response))
	if err != nil {
		return fmt.Errorf("failed to send login response: %w", err)
	}

	sentence, err := c.protocol.ReadSentence()
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	loginResponse := ParseSentence(sentence)
	if loginResponse.Type == ResponseTrap {
		return ErrAuthenticationFailed
	}

	if loginResponse.Type != ResponseDone {
		return fmt.Errorf("unexpected response during login: %v", loginResponse.Type)
	}

	c.config.Logger.Debug("MD5 authentication successful")
	return nil
}

// readResponses continuously reads responses from the connection
func (c *Client) readResponses() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			if !c.connected.Load() {
				return
			}

			sentence, err := c.protocol.ReadSentence()
			if err != nil {
				c.handleError(err)
				return
			}

			response := ParseSentence(sentence)

			// Check for tag and route to appropriate channel
			if response.Tag != "" {
				c.routeResponse(response)
			} else {
				// No tag, send to notification channel or handler
				if c.responseHandler != nil {
					c.responseHandler(response)
				} else {
					select {
					case c.notifyCh <- response:
					default:
						// Channel full, drop response
						c.config.Logger.Warn("Response channel full, dropping response")
					}
				}
			}
		}
	}
}

// routeResponse sends a response to the appropriate waiting goroutine
func (c *Client) routeResponse(response *Response) {
	var tag uint32
	fmt.Sscanf(response.Tag, "%d", &tag)

	c.mu.RLock()
	ch, ok := c.pending[tag]
	c.mu.RUnlock()

	if ok {
		ch <- response
	}
}

// handleError handles connection errors
func (c *Client) handleError(err error) {
	c.config.Logger.Error("Connection error: %v", err)

	if c.errorHandler != nil {
		c.errorHandler(err)
	}

	// Attempt reconnection if enabled
	if c.config.AutoReconnect && !c.closed.Load() {
		go c.reconnect()
	} else {
		c.Disconnect()
	}
}

// reconnect attempts to reestablish the connection
func (c *Client) reconnect() {
	c.config.Logger.Info("Attempting to reconnect...")

	c.connected.Store(false)

	if c.conn != nil {
		c.conn.Close()
	}

	for attempt := 0; attempt < c.config.RetryCount; attempt++ {
		err := c.Connect()
		if err == nil {
			c.config.Logger.Info("Reconnected successfully")
			return
		}

		c.config.Logger.Debug("Reconnection attempt %d failed: %v", attempt+1, err)
		time.Sleep(c.config.RetryDelay)
	}

	c.config.Logger.Error("Failed to reconnect after %d attempts", c.config.RetryCount)
}

// Disconnect closes the connection to the RouterOS device
func (c *Client) Disconnect() error {
	if !c.connected.Load() {
		return nil
	}

	c.closed.Store(true)
	c.connected.Store(false)

	c.cancel()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		if err != nil {
			return err
		}
	}

	c.config.Logger.Info("Disconnected from %s:%d", c.config.Address, c.config.Port)
	return nil
}

// IsConnected returns whether the client is currently connected
func (c *Client) IsConnected() bool {
	return c.connected.Load()
}

// IsClosed returns whether the client has been closed
func (c *Client) IsClosed() bool {
	return c.closed.Load()
}

// nextTag generates the next command tag
func (c *Client) nextTag() uint32 {
	return c.tagCounter.Add(1)
}

// Send sends a command and returns all responses
func (c *Client) Send(command string, args ...string) ([]*Response, error) {
	return c.SendContext(c.ctx, command, args...)
}

// SendContext sends a command with context and returns all responses
func (c *Client) SendContext(ctx context.Context, command string, args ...string) ([]*Response, error) {
	if !c.connected.Load() {
		return nil, ErrNotConnected
	}

	// Generate tag
	tag := c.nextTag()
	tagStr := fmt.Sprintf("%d", tag)

	// Create response channel
	ch := make(chan *Response, 100)

	c.mu.Lock()
	c.pending[tag] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, tag)
		c.mu.Unlock()
		close(ch)
	}()

	// Build and send command
	words := make([]string, 0, len(args)+2)
	words = append(words, command)
	words = append(words, args...)
	words = append(words, ".tag="+tagStr)

	c.config.Logger.Debug("Sending command: %s", command)

	if err := c.protocol.WriteSentence(words...); err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	// Collect responses
	responses := make([]*Response, 0)
	timeout := time.After(c.config.ReadTimeout)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case <-timeout:
			return nil, NewTimeoutError("read response", c.config.ReadTimeout.String(), nil)

		case response := <-ch:
			responses = append(responses, response)

			// Check if this is the final response
			if response.Done || response.Type == ResponseTrap || response.Type == ResponseFatal {
				if response.Type == ResponseTrap || response.Type == ResponseFatal {
					return responses, NewRouterError(response.Data)
				}
				return responses, nil
			}
		}
	}
}

// SendAsync sends a command without waiting for responses
func (c *Client) SendAsync(command string, args ...string) (uint32, error) {
	if !c.connected.Load() {
		return 0, ErrNotConnected
	}

	tag := c.nextTag()
	tagStr := fmt.Sprintf("%d", tag)

	words := make([]string, 0, len(args)+2)
	words = append(words, command)
	words = append(words, args...)
	words = append(words, ".tag="+tagStr)

	if err := c.protocol.WriteSentence(words...); err != nil {
		return 0, fmt.Errorf("failed to send command: %w", err)
	}

	return tag, nil
}

// Listen starts listening for asynchronous responses for a given tag
func (c *Client) Listen(tag uint32) (<-chan *Response, error) {
	ch := make(chan *Response, 100)

	c.mu.Lock()
	c.pending[tag] = ch
	c.mu.Unlock()

	return ch, nil
}

// Cancel cancels a pending command by tag
func (c *Client) Cancel(tag uint32) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.pending[tag]; ok {
		delete(c.pending, tag)
		close(ch)
	}

	return nil
}

// SetResponseHandler sets a handler for unsolicited responses
func (c *Client) SetResponseHandler(handler ResponseHandler) {
	c.responseHandler = handler
}

// SetErrorHandler sets a handler for connection errors
func (c *Client) SetErrorHandler(handler ErrorHandler) {
	c.errorHandler = handler
}

// Notifications returns the notification channel
func (c *Client) Notifications() <-chan *Response {
	return c.notifyCh
}

// Config returns the client configuration
func (c *Client) Config() *Config {
	return c.config
}

// Close closes the client and releases resources
func (c *Client) Close() error {
	return c.Disconnect()
}

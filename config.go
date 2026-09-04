package routeros

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"time"
)

// Default connection settings
const (
	DefaultPort         = 8728   // Default API port
	DefaultPortTLS      = 8729   // Default API TLS port
	DefaultTimeout      = 30     // Default timeout in seconds
	DefaultReadTimeout  = 30     // Default read timeout in seconds
	DefaultWriteTimeout = 30     // Default write timeout in seconds
	DefaultRetryCount   = 3      // Default number of retries
	DefaultRetryDelay   = 1      // Default delay between retries in seconds
)

// Config holds the configuration for connecting to a RouterOS device
type Config struct {
	// Connection settings
	Address      string        // Router IP address or hostname
	Port         int           // API port (default: 8728, TLS: 8729)
	Username     string        // Username for authentication
	Password     string        // Password for authentication
	Timeout      time.Duration // Connection timeout
	ReadTimeout  time.Duration // Read operation timeout
	WriteTimeout time.Duration // Write operation timeout

	// TLS/SSL settings
	TLS         bool        // Enable TLS connection
	TLSConfig   *tls.Config // TLS configuration
	SkipVerify  bool        // Skip TLS certificate verification
	CACert      string      // Path to CA certificate file
	ClientCert  string      // Path to client certificate file
	ClientKey   string      // Path to client key file

	// Connection behavior
	AutoReconnect bool          // Automatically reconnect on disconnect
	RetryCount    int           // Number of connection retries
	RetryDelay    time.Duration // Delay between retries
	KeepAlive     bool          // Enable TCP keepalive
	KeepAliveTime time.Duration // TCP keepalive interval

	// Advanced settings
	BufferSize    int           // Read/write buffer size
	QueueSize     int           // Command queue size
	LogLevel      LogLevel      // Logging level
	Logger        Logger        // Custom logger
	UserAgent     string        // User agent string
	SourceAddress string        // Source IP address for connection
}

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelNone LogLevel = iota
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
	LogLevelTrace
)

// Logger interface for custom logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// DefaultLogger is a minimal logger implementation
type DefaultLogger struct {
	level LogLevel
}

// NewDefaultLogger creates a new DefaultLogger
func NewDefaultLogger(level LogLevel) *DefaultLogger {
	return &DefaultLogger{level: level}
}

func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	if l.level >= LogLevelDebug {
		fmt.Printf("[DEBUG] "+msg+"\n", args...)
	}
}

func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	if l.level >= LogLevelInfo {
		fmt.Printf("[INFO] "+msg+"\n", args...)
	}
}

func (l *DefaultLogger) Warn(msg string, args ...interface{}) {
	if l.level >= LogLevelWarning {
		fmt.Printf("[WARN] "+msg+"\n", args...)
	}
}

func (l *DefaultLogger) Error(msg string, args ...interface{}) {
	if l.level >= LogLevelError {
		fmt.Printf("[ERROR] "+msg+"\n", args...)
	}
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		Port:          DefaultPort,
		Timeout:       time.Duration(DefaultTimeout) * time.Second,
		ReadTimeout:   time.Duration(DefaultReadTimeout) * time.Second,
		WriteTimeout:  time.Duration(DefaultWriteTimeout) * time.Second,
		AutoReconnect: false,
		RetryCount:    DefaultRetryCount,
		RetryDelay:    time.Duration(DefaultRetryDelay) * time.Second,
		KeepAlive:     true,
		KeepAliveTime: 30 * time.Second,
		BufferSize:    4096,
		QueueSize:     100,
		LogLevel:      LogLevelError,
		Logger:        NewDefaultLogger(LogLevelError),
	}
}

// DefaultConfig returns a Config with default settings
func DefaultConfig() *Config {
	return NewConfig()
}

// WithAddress sets the router address
func (c *Config) WithAddress(addr string) *Config {
	c.Address = addr
	return c
}

// WithPort sets the API port
func (c *Config) WithPort(port int) *Config {
	c.Port = port
	return c
}

// WithCredentials sets the username and password
func (c *Config) WithCredentials(username, password string) *Config {
	c.Username = username
	c.Password = password
	return c
}

// WithTimeout sets the connection timeout
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}

// WithReadTimeout sets the read timeout
func (c *Config) WithReadTimeout(timeout time.Duration) *Config {
	c.ReadTimeout = timeout
	return c
}

// WithWriteTimeout sets the write timeout
func (c *Config) WithWriteTimeout(timeout time.Duration) *Config {
	c.WriteTimeout = timeout
	return c
}

// WithTLS enables TLS connection
func (c *Config) WithTLS() *Config {
	c.TLS = true
	if c.Port == DefaultPort {
		c.Port = DefaultPortTLS
	}
	return c
}

// WithTLSConfig sets custom TLS configuration
func (c *Config) WithTLSConfig(config *tls.Config) *Config {
	c.TLS = true
	c.TLSConfig = config
	if c.Port == DefaultPort {
		c.Port = DefaultPortTLS
	}
	return c
}

// WithSkipVerify skips TLS certificate verification
func (c *Config) WithSkipVerify(skip bool) *Config {
	c.SkipVerify = skip
	return c
}

// WithCACert sets the CA certificate path
func (c *Config) WithCACert(path string) *Config {
	c.CACert = path
	return c
}

// WithClientCert sets the client certificate and key paths
func (c *Config) WithClientCert(certPath, keyPath string) *Config {
	c.ClientCert = certPath
	c.ClientKey = keyPath
	return c
}

// WithAutoReconnect enables auto-reconnect
func (c *Config) WithAutoReconnect(enabled bool) *Config {
	c.AutoReconnect = enabled
	return c
}

// WithRetry sets retry count and delay
func (c *Config) WithRetry(count int, delay time.Duration) *Config {
	c.RetryCount = count
	c.RetryDelay = delay
	return c
}

// WithKeepAlive enables TCP keepalive
func (c *Config) WithKeepAlive(enabled bool, interval time.Duration) *Config {
	c.KeepAlive = enabled
	c.KeepAliveTime = interval
	return c
}

// WithBufferSize sets the buffer size
func (c *Config) WithBufferSize(size int) *Config {
	c.BufferSize = size
	return c
}

// WithQueueSize sets the command queue size
func (c *Config) WithQueueSize(size int) *Config {
	c.QueueSize = size
	return c
}

// WithLogger sets a custom logger
func (c *Config) WithLogger(logger Logger) *Config {
	c.Logger = logger
	return c
}

// WithLogLevel sets the logging level
func (c *Config) WithLogLevel(level LogLevel) *Config {
	c.LogLevel = level
	if dl, ok := c.Logger.(*DefaultLogger); ok {
		dl.level = level
	}
	return c
}

// WithUserAgent sets the user agent string
func (c *Config) WithUserAgent(agent string) *Config {
	c.UserAgent = agent
	return c
}

// WithSourceAddress sets the source IP address
func (c *Config) WithSourceAddress(addr string) *Config {
	c.SourceAddress = addr
	return c
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}
	if c.Username == "" {
		return fmt.Errorf("username is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return ErrInvalidPort
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	return nil
}

// BuildTLSConfig builds the TLS configuration
func (c *Config) BuildTLSConfig() (*tls.Config, error) {
	if !c.TLS {
		return nil, nil
	}

	config := &tls.Config{
		InsecureSkipVerify: c.SkipVerify,
	}

	// If custom TLS config is provided, use it as base
	if c.TLSConfig != nil {
		config = c.TLSConfig.Clone()
	}

	// Load CA certificate if provided
	if c.CACert != "" {
		caCert, err := os.ReadFile(c.CACert)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		config.RootCAs = caCertPool
	}

	// Load client certificate if provided
	if c.ClientCert != "" && c.ClientKey != "" {
		cert, err := tls.LoadX509KeyPair(c.ClientCert, c.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		config.Certificates = []tls.Certificate{cert}
	}

	return config, nil
}

// NetworkAddress returns the address in host:port format
func (c *Config) NetworkAddress() string {
	return net.JoinHostPort(c.Address, fmt.Sprintf("%d", c.Port))
}

// Clone creates a copy of the configuration
func (c *Config) Clone() *Config {
	clone := &Config{
		Address:       c.Address,
		Port:          c.Port,
		Username:      c.Username,
		Password:      c.Password,
		Timeout:       c.Timeout,
		ReadTimeout:   c.ReadTimeout,
		WriteTimeout:  c.WriteTimeout,
		TLS:           c.TLS,
		TLSConfig:     c.TLSConfig,
		SkipVerify:    c.SkipVerify,
		CACert:        c.CACert,
		ClientCert:    c.ClientCert,
		ClientKey:     c.ClientKey,
		AutoReconnect: c.AutoReconnect,
		RetryCount:    c.RetryCount,
		RetryDelay:    c.RetryDelay,
		KeepAlive:     c.KeepAlive,
		KeepAliveTime: c.KeepAliveTime,
		BufferSize:    c.BufferSize,
		QueueSize:     c.QueueSize,
		LogLevel:      c.LogLevel,
		Logger:        c.Logger,
		UserAgent:     c.UserAgent,
		SourceAddress: c.SourceAddress,
	}
	return clone
}

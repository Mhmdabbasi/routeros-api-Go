package routeros

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// API is the main RouterOS API client. It exposes the transport client and
// generic resource handles for common RouterOS menus.
//
// Specialized resources that need a dedicated operation should use the Client
// with SendContext or a typed adapter. The legacy files under resources/ are
// not part of the root Go package and must not be referenced here as root types.
type API struct {
	client *Client

	System      *SystemResource
	Interface   *BaseResource
	IP          *BaseResource
	IPv6        *BaseResource
	PPP         *BaseResource
	Queue       *BaseResource
	Routing     *BaseResource
	Tools       *BaseResource
	Certificate *BaseResource
	Wireless    *BaseResource
	Log         *BaseResource
}

// New creates a new RouterOS API client with default configuration
func New(address, username, password string) (*API, error) {
	config := NewConfig().
		WithAddress(address).
		WithCredentials(username, password)

	return NewWithConfig(config)
}

// NewTLS creates a new RouterOS API client with TLS
func NewTLS(address, username, password string) (*API, error) {
	config := NewConfig().
		WithAddress(address).
		WithCredentials(username, password).
		WithTLS()

	return NewWithConfig(config)
}

// NewWithConfig creates a new RouterOS API client with custom configuration
func NewWithConfig(config *Config) (*API, error) {
	client := NewClient(config)
	api := &API{
		client: client,
	}

	// Initialize root-package resource handles.
	api.System = NewSystemResource(client)
	api.Interface = NewBaseResource(client, "/interface")
	api.IP = NewBaseResource(client, "/ip/address")
	api.IPv6 = NewBaseResource(client, "/ipv6/address")
	api.PPP = NewBaseResource(client, "/ppp/secret")
	api.Queue = NewBaseResource(client, "/queue/simple")
	api.Routing = NewBaseResource(client, "/routing/route")
	api.Tools = NewBaseResource(client, "/tool")
	api.Certificate = NewBaseResource(client, "/certificate")
	api.Wireless = NewBaseResource(client, "/interface/wireless")
	api.Log = NewBaseResource(client, "/log")

	return api, nil
}

// Connect establishes a connection to the RouterOS device
func (a *API) Connect() error {
	return a.client.Connect()
}

// Disconnect closes the connection
func (a *API) Disconnect() error {
	return a.client.Disconnect()
}

// Close closes the connection (alias for Disconnect)
func (a *API) Close() error {
	return a.client.Close()
}

// IsConnected returns whether the client is connected
func (a *API) IsConnected() bool {
	return a.client.IsConnected()
}

// Send sends a raw command to the RouterOS device
func (a *API) Send(command string, args ...string) ([]*Response, error) {
	return a.client.Send(command, args...)
}

// SendContext sends a raw command with context
func (a *API) SendContext(ctx context.Context, command string, args ...string) ([]*Response, error) {
	return a.client.SendContext(ctx, command, args...)
}

// Query creates a new query builder
func (a *API) Query(command string) *QueryBuilder {
	return NewQuery(command)
}

// Client returns the underlying client
func (a *API) Client() *Client {
	return a.client
}

// Config returns the client configuration
func (a *API) Config() *Config {
	return a.client.Config()
}

// Quick connection helpers

// QuickConnect establishes a quick connection to a RouterOS device
func QuickConnect(address, username, password string) (*API, error) {
	api, err := New(address, username, password)
	if err != nil {
		return nil, err
	}
	if err := api.Connect(); err != nil {
		return nil, err
	}
	return api, nil
}

// QuickConnectTLS establishes a quick TLS connection to a RouterOS device
func QuickConnectTLS(address, username, password string) (*API, error) {
	api, err := NewTLS(address, username, password)
	if err != nil {
		return nil, err
	}
	if err := api.Connect(); err != nil {
		return nil, err
	}
	return api, nil
}

// QuickConnectWithTimeout establishes a connection with custom timeout
func QuickConnectWithTimeout(address, username, password string, timeout time.Duration) (*API, error) {
	config := NewConfig().
		WithAddress(address).
		WithCredentials(username, password).
		WithTimeout(timeout)

	api, err := NewWithConfig(config)
	if err != nil {
		return nil, err
	}
	if err := api.Connect(); err != nil {
		return nil, err
	}
	return api, nil
}

// QuickConnectTLSWithConfig establishes a TLS connection with custom config
func QuickConnectTLSWithConfig(address, username, password string, tlsConfig *tls.Config) (*API, error) {
	config := NewConfig().
		WithAddress(address).
		WithCredentials(username, password).
		WithTLSConfig(tlsConfig)

	api, err := NewWithConfig(config)
	if err != nil {
		return nil, err
	}
	if err := api.Connect(); err != nil {
		return nil, err
	}
	return api, nil
}

// Execute executes a command and returns the parsed response
func Execute(config *Config, command string, args ...string) ([]*Response, error) {
	client := NewClient(config)
	if err := client.Connect(); err != nil {
		return nil, err
	}
	defer client.Close()

	return client.Send(command, args...)
}

// ExecuteContext executes a command with context
func ExecuteContext(ctx context.Context, config *Config, command string, args ...string) ([]*Response, error) {
	client := NewClient(config)
	if err := client.Connect(); err != nil {
		return nil, err
	}
	defer client.Close()

	return client.SendContext(ctx, command, args...)
}

// Dial establishes a connection and returns the client
func Dial(address, username, password string) (*Client, error) {
	config := NewConfig().
		WithAddress(address).
		WithCredentials(username, password)

	client := NewClient(config)
	if err := client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}

// DialTLS establishes a TLS connection and returns the client
func DialTLS(address, username, password string) (*Client, error) {
	config := NewConfig().
		WithAddress(address).
		WithCredentials(username, password).
		WithTLS()

	client := NewClient(config)
	if err := client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}

// DialContext establishes a connection with context
func DialContext(ctx context.Context, address, username, password string) (*Client, error) {
	config := NewConfig().
		WithAddress(address).
		WithCredentials(username, password)

	client := NewClient(config)
	if err := client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}

// DialTimeout establishes a connection with timeout
func DialTimeout(address, username, password string, timeout time.Duration) (*Client, error) {
	config := NewConfig().
		WithAddress(address).
		WithCredentials(username, password).
		WithTimeout(timeout)

	client := NewClient(config)
	if err := client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}

// NetworkAddress returns the network address for the connection
func NetworkAddress(host string, port int) string {
	return net.JoinHostPort(host, fmt.Sprintf("%d", port))
}

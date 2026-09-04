package routeros

import (
	"context"
)

// Tools handles RouterOS tool operations
type Tools struct {
	client *Client
}

// NewTools creates a new Tools resource handler
func NewTools(client *Client) *Tools {
	return &Tools{client: client}
}

// Ping executes a ping command
func (t *Tools) Ping(address string, count int) ([]*Response, error) {
	return t.PingContext(context.Background(), address, count)
}

// PingContext executes a ping command with context
func (t *Tools) PingContext(ctx context.Context, address string, count int) ([]*Response, error) {
	args := []string{"=address=" + address}
	if count > 0 {
		args = append(args, "=count="+string(rune(count)))
	}
	return t.client.SendContext(ctx, "/ping", args...)
}

// Traceroute executes a traceroute command
func (t *Tools) Traceroute(address string) ([]*Response, error) {
	return t.TracerouteContext(context.Background(), address)
}

// TracerouteContext executes a traceroute command with context
func (t *Tools) TracerouteContext(ctx context.Context, address string) ([]*Response, error) {
	return t.client.SendContext(ctx, "/tool/traceroute", "=address="+address)
}

// BandwidthTest executes a bandwidth test
func (t *Tools) BandwidthTest(address string, direction string) ([]*Response, error) {
	return t.BandwidthTestContext(context.Background(), address, direction)
}

// BandwidthTestContext executes a bandwidth test with context
func (t *Tools) BandwidthTestContext(ctx context.Context, address string, direction string) ([]*Response, error) {
	return t.client.SendContext(ctx, "/tool/bandwidth-test", "=address="+address, "=direction="+direction)
}

// Torch executes a torch command
func (t *Tools) Torch(interface_ string) ([]*Response, error) {
	return t.TorchContext(context.Background(), interface_)
}

// TorchContext executes a torch command with context
func (t *Tools) TorchContext(ctx context.Context, interface_ string) ([]*Response, error) {
	return t.client.SendContext(ctx, "/tool/torch", "=interface="+interface_)
}

// Sniffer handles packet sniffer operations
type Sniffer struct {
	client *Client
}

// NewSniffer creates a new Sniffer handler
func NewSniffer(client *Client) *Sniffer {
	return &Sniffer{client: client}
}

// Start starts the packet sniffer
func (s *Sniffer) Start() error {
	return s.StartContext(context.Background())
}

// StartContext starts the packet sniffer with context
func (s *Sniffer) StartContext(ctx context.Context) error {
	_, err := s.client.SendContext(ctx, "/tool/sniffer/start")
	return err
}

// Stop stops the packet sniffer
func (s *Sniffer) Stop() error {
	return s.StopContext(context.Background())
}

// StopContext stops the packet sniffer with context
func (s *Sniffer) StopContext(ctx context.Context) error {
	_, err := s.client.SendContext(ctx, "/tool/sniffer/stop")
	return err
}

// Save saves the captured packets
func (s *Sniffer) Save(filename string) error {
	return s.SaveContext(context.Background(), filename)
}

// SaveContext saves the captured packets with context
func (s *Sniffer) SaveContext(ctx context.Context, filename string) error {
	_, err := s.client.SendContext(ctx, "/tool/sniffer/save", "=file-name="+filename)
	return err
}

// Packets returns captured packets
func (s *Sniffer) Packets() (*ResourceList, error) {
	return s.PacketsContext(context.Background())
}

// PacketsContext returns captured packets with context
func (s *Sniffer) PacketsContext(ctx context.Context) (*ResourceList, error) {
	responses, err := s.client.SendContext(ctx, "/tool/sniffer/packet/print")
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// SetFilter sets sniffer filter
func (s *Sniffer) SetFilter(interface_ string, srcAddress string, dstAddress string, port string) error {
	return s.SetFilterContext(context.Background(), interface_, srcAddress, dstAddress, port)
}

// SetFilterContext sets sniffer filter with context
func (s *Sniffer) SetFilterContext(ctx context.Context, interface_ string, srcAddress string, dstAddress string, port string) error {
	args := []string{}
	if interface_ != "" {
		args = append(args, "=interface="+interface_)
	}
	if srcAddress != "" {
		args = append(args, "=src-address="+srcAddress)
	}
	if dstAddress != "" {
		args = append(args, "=dst-address="+dstAddress)
	}
	if port != "" {
		args = append(args, "=port="+port)
	}
	_, err := s.client.SendContext(ctx, "/tool/sniffer/set", args...)
	return err
}

// TrafficGenerator handles traffic generator operations
type TrafficGenerator struct {
	client *Client
}

// NewTrafficGenerator creates a new TrafficGenerator handler
func NewTrafficGenerator(client *Client) *TrafficGenerator {
	return &TrafficGenerator{client: client}
}

// Start starts the traffic generator
func (t *TrafficGenerator) Start() error {
	return t.StartContext(context.Background())
}

// StartContext starts the traffic generator with context
func (t *TrafficGenerator) StartContext(ctx context.Context) error {
	_, err := t.client.SendContext(ctx, "/tool/traffic-generator/start")
	return err
}

// Stop stops the traffic generator
func (t *TrafficGenerator) Stop() error {
	return t.StopContext(context.Background())
}

// StopContext stops the traffic generator with context
func (t *TrafficGenerator) StopContext(ctx context.Context) error {
	_, err := t.client.SendContext(ctx, "/tool/traffic-generator/stop")
	return err
}

// Graph handles graphing operations
type Graph struct {
	*BaseResource
}

// NewGraph creates a new Graph handler
func NewGraph(client *Client) *Graph {
	return &Graph{
		BaseResource: NewBaseResource(client, "/tool/graphing"),
	}
}

// GraphInterface handles interface graphing
type GraphInterface struct {
	*BaseResource
}

// NewGraphInterface creates a new GraphInterface handler
func NewGraphInterface(client *Client) *GraphInterface {
	return &GraphInterface{
		BaseResource: NewBaseResource(client, "/tool/graphing/interface"),
	}
}

// GraphQueue handles queue graphing
type GraphQueue struct {
	*BaseResource
}

// NewGraphQueue creates a new GraphQueue handler
func NewGraphQueue(client *Client) *GraphQueue {
	return &GraphQueue{
		BaseResource: NewBaseResource(client, "/tool/graphing/queue"),
	}
}

// GraphResource handles resource graphing
type GraphResource struct {
	*BaseResource
}

// NewGraphResource creates a new GraphResource handler
func NewGraphResource(client *Client) *GraphResource {
	return &GraphResource{
		BaseResource: NewBaseResource(client, "/tool/graphing/resource"),
	}
}

// Netwatch handles netwatch operations
type Netwatch struct {
	*BaseResource
}

// NewNetwatch creates a new Netwatch handler
func NewNetwatch(client *Client) *Netwatch {
	return &Netwatch{
		BaseResource: NewBaseResource(client, "/tool/netwatch"),
	}
}

// Create creates a new netwatch entry
func (n *Netwatch) Create(host string, interval string) (string, error) {
	return n.CreateContext(context.Background(), host, interval)
}

// CreateContext creates a new netwatch entry with context
func (n *Netwatch) CreateContext(ctx context.Context, host string, interval string) (string, error) {
	args := map[string]string{
		"host": host,
	}
	if interval != "" {
		args["interval"] = interval
	}
	return n.AddContext(ctx, args)
}

// MACServer handles MAC server operations
type MACServer struct {
	client *Client
}

// NewMACServer creates a new MACServer handler
func NewMACServer(client *Client) *MACServer {
	return &MACServer{client: client}
}

// NetworkBootstrap handles network bootstrap operations
type NetworkBootstrap struct {
	client *Client
}

// NewNetworkBootstrap creates a new NetworkBootstrap handler
func NewNetworkBootstrap(client *Client) *NetworkBootstrap {
	return &NetworkBootstrap{client: client}
}

// Fetch handles fetch operations
type Fetch struct {
	client *Client
}

// NewFetch creates a new Fetch handler
func NewFetch(client *Client) *Fetch {
	return &Fetch{client: client}
}

// URL fetches a URL
func (f *Fetch) URL(url, dstPath string) error {
	return f.URLContext(context.Background(), url, dstPath)
}

// URLContext fetches a URL with context
func (f *Fetch) URLContext(ctx context.Context, url, dstPath string) error {
	args := []string{"=url=" + url}
	if dstPath != "" {
		args = append(args, "=dst-path="+dstPath)
	}
	_, err := f.client.SendContext(ctx, "/tool/fetch", args...)
	return err
}

// Email handles email operations
type Email struct {
	client *Client
}

// NewEmail creates a new Email handler
func NewEmail(client *Client) *Email {
	return &Email{client: client}
}

// Send sends an email
func (e *Email) Send(to, subject, body string) error {
	return e.SendContext(context.Background(), to, subject, body)
}

// SendContext sends an email with context
func (e *Email) SendContext(ctx context.Context, to, subject, body string) error {
	_, err := e.client.SendContext(ctx, "/tool/e-mail/send",
		"=to="+to,
		"=subject="+subject,
		"=body="+body,
	)
	return err
}

// SMS handles SMS operations
type SMS struct {
	client *Client
}

// NewSMS creates a new SMS handler
func NewSMS(client *Client) *SMS {
	return &SMS{client: client}
}

// Send sends an SMS
func (s *SMS) Send(phone, message string) error {
	return s.SendContext(context.Background(), phone, message)
}

// SendContext sends an SMS with context
func (s *SMS) SendContext(ctx context.Context, phone, message string) error {
	_, err := s.client.SendContext(ctx, "/tool/sms/send",
		"=phone-number="+phone,
		"=message="+message,
	)
	return err
}

// UserManager handles user manager operations
type UserManager struct {
	client *Client
}

// NewUserManager creates a new UserManager handler
func NewUserManager(client *Client) *UserManager {
	return &UserManager{client: client}
}

// Profile handles user manager profiles
type UMProfile struct {
	*BaseResource
}

// NewUMProfile creates a new UMProfile handler
func NewUMProfile(client *Client) *UMProfile {
	return &UMProfile{
		BaseResource: NewBaseResource(client, "/user-manager/profile"),
	}
}

// User handles user manager users
type UMUser struct {
	*BaseResource
}

// NewUMUser creates a new UMUser handler
func NewUMUser(client *Client) *UMUser {
	return &UMUser{
		BaseResource: NewBaseResource(client, "/user-manager/user"),
	}
}

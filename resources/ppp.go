package routeros

import (
	"context"
)

// PPP handles RouterOS PPP operations
type PPP struct {
	client *Client
}

// NewPPP creates a new PPP resource handler
func NewPPP(client *Client) *PPP {
	return &PPP{client: client}
}

// Profile handles PPP profile operations
type PPPProfile struct {
	*BaseResource
}

// NewPPPProfile creates a new PPPProfile handler
func NewPPPProfile(client *Client) *PPPProfile {
	return &PPPProfile{
		BaseResource: NewBaseResource(client, "/ppp/profile"),
	}
}

// CreateProfile creates a new PPP profile
func (p *PPPProfile) CreateProfile(name string, localAddress, remoteAddress string) (string, error) {
	return p.CreateProfileContext(context.Background(), name, localAddress, remoteAddress)
}

// CreateProfileContext creates a new PPP profile with context
func (p *PPPProfile) CreateProfileContext(ctx context.Context, name string, localAddress, remoteAddress string) (string, error) {
	args := map[string]string{"name": name}
	if localAddress != "" {
		args["local-address"] = localAddress
	}
	if remoteAddress != "" {
		args["remote-address"] = remoteAddress
	}
	return p.AddContext(ctx, args)
}

// Secret handles PPP secret (user) operations
type PPPSecret struct {
	*BaseResource
}

// NewPPPSecret creates a new PPPSecret handler
func NewPPPSecret(client *Client) *PPPSecret {
	return &PPPSecret{
		BaseResource: NewBaseResource(client, "/ppp/secret"),
	}
}

// CreateUser creates a new PPP user
func (s *PPPSecret) CreateUser(name, password, profile, service string) (string, error) {
	return s.CreateUserContext(context.Background(), name, password, profile, service)
}

// CreateUserContext creates a new PPP user with context
func (s *PPPSecret) CreateUserContext(ctx context.Context, name, password, profile, service string) (string, error) {
	args := map[string]string{
		"name":     name,
		"password": password,
	}
	if profile != "" {
		args["profile"] = profile
	}
	if service != "" {
		args["service"] = service
	}
	return s.AddContext(ctx, args)
}

// SetPassword sets a user's password
func (s *PPPSecret) SetPassword(id, password string) error {
	return s.SetPasswordContext(context.Background(), id, password)
}

// SetPasswordContext sets a user's password with context
func (s *PPPSecret) SetPasswordContext(ctx context.Context, id, password string) error {
	return s.UpdateContext(ctx, id, map[string]string{"password": password})
}

// SetProfile sets a user's profile
func (s *PPPSecret) SetProfile(id, profile string) error {
	return s.SetProfileContext(context.Background(), id, profile)
}

// SetProfileContext sets a user's profile with context
func (s *PPPSecret) SetProfileContext(ctx context.Context, id, profile string) error {
	return s.UpdateContext(ctx, id, map[string]string{"profile": profile})
}

// Active handles active PPP sessions
type PPPActive struct {
	*BaseResource
}

// NewPPPActive creates a new PPPActive handler
func NewPPPActive(client *Client) *PPPActive {
	return &PPPActive{
		BaseResource: NewBaseResource(client, "/ppp/active"),
	}
}

// GetByUser finds an active session by username
func (a *PPPActive) GetByUser(username string) (*Resource, error) {
	return a.GetByUserContext(context.Background(), username)
}

// GetByUserContext finds an active session by username with context
func (a *PPPActive) GetByUserContext(ctx context.Context, username string) (*Resource, error) {
	return a.FindOneContext(ctx, map[string]string{"name": username})
}

// Disconnect disconnects an active session
func (a *PPPActive) Disconnect(id string) error {
	return a.DisconnectContext(context.Background(), id)
}

// DisconnectContext disconnects an active session with context
func (a *PPPActive) DisconnectContext(ctx context.Context, id string) error {
	return a.DeleteContext(ctx, id)
}

// AAA handles PPP AAA settings
type PPPAAA struct {
	client *Client
}

// NewPPPAAA creates a new PPPAAA handler
func NewPPPAAA(client *Client) *PPPAAA {
	return &PPPAAA{client: client}
}

// Get returns PPP AAA settings
func (a *PPPAAA) Get() (map[string]string, error) {
	return a.GetContext(context.Background())
}

// GetContext returns PPP AAA settings with context
func (a *PPPAAA) GetContext(ctx context.Context) (map[string]string, error) {
	responses, err := a.client.SendContext(ctx, "/ppp/aaa/print")
	if err != nil {
		return nil, err
	}
	parser := NewResponseParser(responses)
	return parser.FirstData(), nil
}

// SetUseRadius enables or disables RADIUS
func (a *PPPAAA) SetUseRadius(useRadius bool) error {
	return a.SetUseRadiusContext(context.Background(), useRadius)
}

// SetUseRadiusContext enables or disables RADIUS with context
func (a *PPPAAA) SetUseRadiusContext(ctx context.Context, useRadius bool) error {
	val := "no"
	if useRadius {
		val = "yes"
	}
	_, err := a.client.SendContext(ctx, "/ppp/aaa/set", "=use-radius="+val)
	return err
}

// PPPoEServer handles PPPoE server operations
type PPPoEServer struct {
	*BaseResource
}

// NewPPPoEServer creates a new PPPoEServer handler
func NewPPPoEServer(client *Client) *PPPoEServer {
	return &PPPoEServer{
		BaseResource: NewBaseResource(client, "/ppp/pppoe-server"),
	}
}

// PPPoEServerService handles PPPoE server service operations
type PPPoEServerService struct {
	*BaseResource
}

// NewPPPoEServerService creates a new PPPoEServerService handler
func NewPPPoEServerService(client *Client) *PPPoEServerService {
	return &PPPoEServerService{
		BaseResource: NewBaseResource(client, "/interface/pppoe-server/server"),
	}
}

// Create creates a PPPoE server service
func (p *PPPoEServerService) Create(name, interface_, service string) (string, error) {
	return p.CreateContext(context.Background(), name, interface_, service)
}

// CreateContext creates a PPPoE server service with context
func (p *PPPoEServerService) CreateContext(ctx context.Context, name, interface_, service string) (string, error) {
	return p.AddContext(ctx, map[string]string{
		"name":      name,
		"interface": interface_,
		"service":   service,
	})
}

// PPTPServer handles PPTP server operations
type PPTPServer struct {
	*BaseResource
}

// NewPPTPServer creates a new PPTPServer handler
func NewPPTPServer(client *Client) *PPTPServer {
	return &PPTPServer{
		BaseResource: NewBaseResource(client, "/interface/pptp-server/server"),
	}
}

// L2TPServer handles L2TP server operations
type L2TPServer struct {
	*BaseResource
}

// NewL2TPServer creates a new L2TPServer handler
func NewL2TPServer(client *Client) *L2TPServer {
	return &L2TPServer{
		BaseResource: NewBaseResource(client, "/interface/l2tp-server/server"),
	}
}

// SSTPServer handles SSTP server operations
type SSTPServer struct {
	*BaseResource
}

// NewSSTPServer creates a new SSTPServer handler
func NewSSTPServer(client *Client) *SSTPServer {
	return &SSTPServer{
		BaseResource: NewBaseResource(client, "/interface/sstp-server/server"),
	}
}

// OVPNServer handles OpenVPN server operations
type OVPNServer struct {
	*BaseResource
}

// NewOVPNServer creates a new OVPNServer handler
func NewOVPNServer(client *Client) *OVPNServer {
	return &OVPNServer{
		BaseResource: NewBaseResource(client, "/interface/ovpn-server/server"),
	}
}

// ISDN handles ISDN operations
type ISDN struct {
	*BaseResource
}

// NewISDN creates a new ISDN handler
func NewISDN(client *Client) *ISDN {
	return &ISDN{
		BaseResource: NewBaseResource(client, "/ppp/isdn"),
	}
}

// BCP handles BCP operations
type BCP struct {
	*BaseResource
}

// NewBCP creates a new BCP handler
func NewBCP(client *Client) *BCP {
	return &BCP{
		BaseResource: NewBaseResource(client, "/interface/bcp"),
	}
}

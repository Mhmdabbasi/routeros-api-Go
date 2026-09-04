package routeros

import (
	"context"
)

// IP handles RouterOS IP operations
type IP struct {
	client *Client
}

// NewIP creates a new IP resource handler
func NewIP(client *Client) *IP {
	return &IP{client: client}
}

// Address handles IP address operations
type IPAddress struct {
	*BaseResource
}

// NewIPAddress creates a new IPAddress handler
func NewIPAddress(client *Client) *IPAddress {
	return &IPAddress{
		BaseResource: NewBaseResource(client, "/ip/address"),
	}
}

// AddAddress adds a new IP address to an interface
func (a *IPAddress) AddAddress(address, interface_ string) (string, error) {
	return a.AddAddressContext(context.Background(), address, interface_)
}

// AddAddressContext adds a new IP address to an interface with context
func (a *IPAddress) AddAddressContext(ctx context.Context, address, interface_ string) (string, error) {
	return a.AddContext(ctx, map[string]string{
		"address":   address,
		"interface": interface_,
	})
}

// Route handles IP route operations
type Route struct {
	*BaseResource
}

// NewRoute creates a new Route handler
func NewRoute(client *Client) *Route {
	return &Route{
		BaseResource: NewBaseResource(client, "/ip/route"),
	}
}

// AddRoute adds a new static route
func (r *Route) AddRoute(dstAddress, gateway string) (string, error) {
	return r.AddRouteContext(context.Background(), dstAddress, gateway)
}

// AddRouteContext adds a new static route with context
func (r *Route) AddRouteContext(ctx context.Context, dstAddress, gateway string) (string, error) {
	return r.AddContext(ctx, map[string]string{
		"dst-address": dstAddress,
		"gateway":     gateway,
	})
}

// AddDefaultRoute adds a default route
func (r *Route) AddDefaultRoute(gateway string) (string, error) {
	return r.AddDefaultRouteContext(context.Background(), gateway)
}

// AddDefaultRouteContext adds a default route with context
func (r *Route) AddDefaultRouteContext(ctx context.Context, gateway string) (string, error) {
	return r.AddRouteContext(ctx, "0.0.0.0/0", gateway)
}

// GetByDestination finds routes by destination
func (r *Route) GetByDestination(dst string) (*ResourceList, error) {
	return r.GetByDestinationContext(context.Background(), dst)
}

// GetByDestinationContext finds routes by destination with context
func (r *Route) GetByDestinationContext(ctx context.Context, dst string) (*ResourceList, error) {
	return r.FindContext(ctx, map[string]string{"dst-address": dst})
}

// DNS handles DNS operations
type DNS struct {
	client *Client
}

// NewDNS creates a new DNS handler
func NewDNS(client *Client) *DNS {
	return &DNS{client: client}
}

// Get returns DNS configuration
func (d *DNS) Get() (map[string]string, error) {
	return d.GetContext(context.Background())
}

// GetContext returns DNS configuration with context
func (d *DNS) GetContext(ctx context.Context) (map[string]string, error) {
	responses, err := d.client.SendContext(ctx, "/ip/dns/print")
	if err != nil {
		return nil, err
	}
	parser := NewResponseParser(responses)
	return parser.FirstData(), nil
}

// SetServers sets DNS servers
func (d *DNS) SetServers(servers string) error {
	return d.SetServersContext(context.Background(), servers)
}

// SetServersContext sets DNS servers with context
func (d *DNS) SetServersContext(ctx context.Context, servers string) error {
	_, err := d.client.SendContext(ctx, "/ip/dns/set", "=servers="+servers)
	return err
}

// SetAllowRemoteRequests allows or disallows remote DNS requests
func (d *DNS) SetAllowRemoteRequests(allow bool) error {
	return d.SetAllowRemoteRequestsContext(context.Background(), allow)
}

// SetAllowRemoteRequestsContext allows or disallows remote DNS requests with context
func (d *DNS) SetAllowRemoteRequestsContext(ctx context.Context, allow bool) error {
	val := "no"
	if allow {
		val = "yes"
	}
	_, err := d.client.SendContext(ctx, "/ip/dns/set", "=allow-remote-requests="+val)
	return err
}

// Cache returns DNS cache entries
func (d *DNS) Cache() (*ResourceList, error) {
	return d.CacheContext(context.Background())
}

// CacheContext returns DNS cache entries with context
func (d *DNS) CacheContext(ctx context.Context) (*ResourceList, error) {
	responses, err := d.client.SendContext(ctx, "/ip/dns/cache/print")
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// FlushCache flushes DNS cache
func (d *DNS) FlushCache() error {
	return d.FlushCacheContext(context.Background())
}

// FlushCacheContext flushes DNS cache with context
func (d *DNS) FlushCacheContext(ctx context.Context) error {
	_, err := d.client.SendContext(ctx, "/ip/dns/cache/flush")
	return err
}

// ARP handles ARP operations
type ARP struct {
	*BaseResource
}

// NewARP creates a new ARP handler
func NewARP(client *Client) *ARP {
	return &ARP{
		BaseResource: NewBaseResource(client, "/ip/arp"),
	}
}

// AddStatic adds a static ARP entry
func (a *ARP) AddStatic(ipAddress, macAddress, interface_ string) (string, error) {
	return a.AddStaticContext(context.Background(), ipAddress, macAddress, interface_)
}

// AddStaticContext adds a static ARP entry with context
func (a *ARP) AddStaticContext(ctx context.Context, ipAddress, macAddress, interface_ string) (string, error) {
	return a.AddContext(ctx, map[string]string{
		"address":     ipAddress,
		"mac-address": macAddress,
		"interface":   interface_,
	})
}

// GetByIP finds ARP entry by IP address
func (a *ARP) GetByIP(ip string) (*Resource, error) {
	return a.GetByIPContext(context.Background(), ip)
}

// GetByIPContext finds ARP entry by IP address with context
func (a *ARP) GetByIPContext(ctx context.Context, ip string) (*Resource, error) {
	return a.FindOneContext(ctx, map[string]string{"address": ip})
}

// GetByMAC finds ARP entry by MAC address
func (a *ARP) GetByMAC(mac string) (*Resource, error) {
	return a.GetByMACContext(context.Background(), mac)
}

// GetByMACContext finds ARP entry by MAC address with context
func (a *ARP) GetByMACContext(ctx context.Context, mac string) (*Resource, error) {
	return a.FindOneContext(ctx, map[string]string{"mac-address": mac})
}

// DHCPClient handles DHCP client operations
type DHCPClient struct {
	*BaseResource
}

// NewDHCPClient creates a new DHCPClient handler
func NewDHCPClient(client *Client) *DHCPClient {
	return &DHCPClient{
		BaseResource: NewBaseResource(client, "/ip/dhcp-client"),
	}
}

// Add creates a new DHCP client
func (d *DHCPClient) Add(interface_ string) (string, error) {
	return d.AddContext(context.Background(), interface_)
}

// AddContext creates a new DHCP client with context
func (d *DHCPClient) AddContext(ctx context.Context, interface_ string) (string, error) {
	return d.BaseResource.AddContext(ctx, map[string]string{
		"interface": interface_,
		"disabled":  "no",
	})
}

// Renew renews DHCP lease
func (d *DHCPClient) Renew(id string) error {
	return d.RenewContext(context.Background(), id)
}

// RenewContext renews DHCP lease with context
func (d *DHCPClient) RenewContext(ctx context.Context, id string) error {
	_, err := d.client.SendContext(ctx, "/ip/dhcp-client/renew", ".id="+id)
	return err
}

// Release releases DHCP lease
func (d *DHCPClient) Release(id string) error {
	return d.ReleaseContext(context.Background(), id)
}

// ReleaseContext releases DHCP lease with context
func (d *DHCPClient) ReleaseContext(ctx context.Context, id string) error {
	_, err := d.client.SendContext(ctx, "/ip/dhcp-client/release", ".id="+id)
	return err
}

// DHCPServer handles DHCP server operations
type DHCPServer struct {
	*BaseResource
}

// NewDHCPServer creates a new DHCPServer handler
func NewDHCPServer(client *Client) *DHCPServer {
	return &DHCPServer{
		BaseResource: NewBaseResource(client, "/ip/dhcp-server"),
	}
}

// Network handles DHCP server network configuration
type DHCPNetwork struct {
	*BaseResource
}

// NewDHCPNetwork creates a new DHCPNetwork handler
func NewDHCPNetwork(client *Client) *DHCPNetwork {
	return &DHCPNetwork{
		BaseResource: NewBaseResource(client, "/ip/dhcp-server/network"),
	}
}

// AddNetwork adds a DHCP network
func (d *DHCPNetwork) AddNetwork(address, gateway string) (string, error) {
	return d.AddNetworkContext(context.Background(), address, gateway)
}

// AddNetworkContext adds a DHCP network with context
func (d *DHCPNetwork) AddNetworkContext(ctx context.Context, address, gateway string) (string, error) {
	return d.AddContext(ctx, map[string]string{
		"address": address,
		"gateway": gateway,
	})
}

// DHCPLease handles DHCP lease operations
type DHCPLease struct {
	*BaseResource
}

// NewDHCPLease creates a new DHCPLease handler
func NewDHCPLease(client *Client) *DHCPLease {
	return &DHCPLease{
		BaseResource: NewBaseResource(client, "/ip/dhcp-server/lease"),
	}
}

// AddStatic adds a static DHCP lease
func (d *DHCPLease) AddStatic(macAddress, address string) (string, error) {
	return d.AddStaticContext(context.Background(), macAddress, address)
}

// AddStaticContext adds a static DHCP lease with context
func (d *DHCPLease) AddStaticContext(ctx context.Context, macAddress, address string) (string, error) {
	return d.AddContext(ctx, map[string]string{
		"mac-address": macAddress,
		"address":     address,
	})
}

// Pool handles IP pool operations
type Pool struct {
	*BaseResource
}

// NewPool creates a new Pool handler
func NewPool(client *Client) *Pool {
	return &Pool{
		BaseResource: NewBaseResource(client, "/ip/pool"),
	}
}

// AddPool creates a new IP pool
func (p *Pool) AddPool(name, ranges string) (string, error) {
	return p.AddPoolContext(context.Background(), name, ranges)
}

// AddPoolContext creates a new IP pool with context
func (p *Pool) AddPoolContext(ctx context.Context, name, ranges string) (string, error) {
	return p.AddContext(ctx, map[string]string{
		"name":   name,
		"ranges": ranges,
	})
}

// Firewall handles firewall operations
type Firewall struct {
	client *Client
}

// NewFirewall creates a new Firewall handler
func NewFirewall(client *Client) *Firewall {
	return &Firewall{client: client}
}

// Filter handles firewall filter rules
type FirewallFilter struct {
	*BaseResource
}

// NewFirewallFilter creates a new FirewallFilter handler
func NewFirewallFilter(client *Client) *FirewallFilter {
	return &FirewallFilter{
		BaseResource: NewBaseResource(client, "/ip/firewall/filter"),
	}
}

// AddRule adds a firewall filter rule
func (f *FirewallFilter) AddRule(chain, action string, attrs map[string]string) (string, error) {
	return f.AddRuleContext(context.Background(), chain, action, attrs)
}

// AddRuleContext adds a firewall filter rule with context
func (f *FirewallFilter) AddRuleContext(ctx context.Context, chain, action string, attrs map[string]string) (string, error) {
	args := map[string]string{
		"chain":  chain,
		"action": action,
	}
	for k, v := range attrs {
		args[k] = v
	}
	return f.AddContext(ctx, args)
}

// NAT handles firewall NAT rules
type FirewallNAT struct {
	*BaseResource
}

// NewFirewallNAT creates a new FirewallNAT handler
func NewFirewallNAT(client *Client) *FirewallNAT {
	return &FirewallNAT{
		BaseResource: NewBaseResource(client, "/ip/firewall/nat"),
	}
}

// AddMasquerade adds a masquerade rule
func (n *FirewallNAT) AddMasquerade(chain, outInterface string) (string, error) {
	return n.AddMasqueradeContext(context.Background(), chain, outInterface)
}

// AddMasqueradeContext adds a masquerade rule with context
func (n *FirewallNAT) AddMasqueradeContext(ctx context.Context, chain, outInterface string) (string, error) {
	return n.AddContext(ctx, map[string]string{
		"chain":        chain,
		"action":       "masquerade",
		"out-interface": outInterface,
	})
}

// AddDstNat adds a destination NAT rule
func (n *FirewallNAT) AddDstNat(chain, dstAddress, toAddresses string, dstPort, toPort string) (string, error) {
	return n.AddDstNatContext(context.Background(), chain, dstAddress, toAddresses, dstPort, toPort)
}

// AddDstNatContext adds a destination NAT rule with context
func (n *FirewallNAT) AddDstNatContext(ctx context.Context, chain, dstAddress, toAddresses string, dstPort, toPort string) (string, error) {
	args := map[string]string{
		"chain":       chain,
		"action":      "dst-nat",
		"to-addresses": toAddresses,
	}
	if dstAddress != "" {
		args["dst-address"] = dstAddress
	}
	if dstPort != "" {
		args["dst-port"] = dstPort
	}
	if toPort != "" {
		args["to-ports"] = toPort
	}
	return n.AddContext(ctx, args)
}

// AddSrcNat adds a source NAT rule
func (n *FirewallNAT) AddSrcNat(chain, srcAddress, toAddresses string) (string, error) {
	return n.AddSrcNatContext(context.Background(), chain, srcAddress, toAddresses)
}

// AddSrcNatContext adds a source NAT rule with context
func (n *FirewallNAT) AddSrcNatContext(ctx context.Context, chain, srcAddress, toAddresses string) (string, error) {
	return n.AddContext(ctx, map[string]string{
		"chain":        chain,
		"action":       "src-nat",
		"src-address":  srcAddress,
		"to-addresses": toAddresses,
	})
}

// Mangle handles firewall mangle rules
type FirewallMangle struct {
	*BaseResource
}

// NewFirewallMangle creates a new FirewallMangle handler
func NewFirewallMangle(client *Client) *FirewallMangle {
	return &FirewallMangle{
		BaseResource: NewBaseResource(client, "/ip/firewall/mangle"),
	}
}

// Raw handles firewall raw rules
type FirewallRaw struct {
	*BaseResource
}

// NewFirewallRaw creates a new FirewallRaw handler
func NewFirewallRaw(client *Client) *FirewallRaw {
	return &FirewallRaw{
		BaseResource: NewBaseResource(client, "/ip/firewall/raw"),
	}
}

// AddressList handles firewall address lists
type FirewallAddressList struct {
	*BaseResource
}

// NewFirewallAddressList creates a new FirewallAddressList handler
func NewFirewallAddressList(client *Client) *FirewallAddressList {
	return &FirewallAddressList{
		BaseResource: NewBaseResource(client, "/ip/firewall/address-list"),
	}
}

// AddAddress adds an address to a list
func (a *FirewallAddressList) AddAddress(list, address string) (string, error) {
	return a.AddAddressContext(context.Background(), list, address)
}

// AddAddressContext adds an address to a list with context
func (a *FirewallAddressList) AddAddressContext(ctx context.Context, list, address string) (string, error) {
	return a.AddContext(ctx, map[string]string{
		"list":    list,
		"address": address,
	})
}

// Service handles IP service operations
type Service struct {
	*BaseResource
}

// NewIPService creates a new IP Service handler
func NewIPService(client *Client) *Service {
	return &Service{
		BaseResource: NewBaseResource(client, "/ip/service"),
	}
}

// Service handles IP service operations
func (s *Service) Enable(name string) error {
	return s.EnableContext(context.Background(), name)
}

// EnableContext enables a service with context
func (s *Service) EnableContext(ctx context.Context, name string) error {
	svc, err := s.FindByNameContext(ctx, name)
	if err != nil {
		return err
	}
	return s.BaseResource.EnableContext(ctx, svc.ID)
}

// Disable disables a service
func (s *Service) Disable(name string) error {
	return s.DisableContext(context.Background(), name)
}

// DisableContext disables a service with context
func (s *Service) DisableContext(ctx context.Context, name string) error {
	svc, err := s.FindByNameContext(ctx, name)
	if err != nil {
		return err
	}
	return s.BaseResource.DisableContext(ctx, svc.ID)
}

// SetPort sets the port for a service
func (s *Service) SetPort(name string, port int) error {
	return s.SetPortContext(context.Background(), name, port)
}

// SetPortContext sets the port for a service with context
func (s *Service) SetPortContext(ctx context.Context, name string, port int) error {
	svc, err := s.FindByNameContext(ctx, name)
	if err != nil {
		return err
	}
	return s.UpdateContext(ctx, svc.ID, map[string]string{"port": string(rune(port))})
}

// UPNP handles UPNP operations
type UPNP struct {
	*BaseResource
}

// NewUPNP creates a new UPNP handler
func NewUPNP(client *Client) *UPNP {
	return &UPNP{
		BaseResource: NewBaseResource(client, "/ip/upnp"),
	}
}

// UPNPInterfaces handles UPNP interface operations
type UPNPInterfaces struct {
	*BaseResource
}

// NewUPNPInterfaces creates a new UPNPInterfaces handler
func NewUPNPInterfaces(client *Client) *UPNPInterfaces {
	return &UPNPInterfaces{
		BaseResource: NewBaseResource(client, "/ip/upnp/interfaces"),
	}
}

// Neighbor handles neighbor discovery
type Neighbor struct {
	*BaseResource
}

// NewNeighbor creates a new Neighbor handler
func NewNeighbor(client *Client) *Neighbor {
	return &Neighbor{
		BaseResource: NewBaseResource(client, "/ip/neighbor"),
	}
}

// Hotspot handles hotspot operations
type Hotspot struct {
	*BaseResource
}

// NewHotspot creates a new Hotspot handler
func NewHotspot(client *Client) *Hotspot {
	return &Hotspot{
		BaseResource: NewBaseResource(client, "/ip/hotspot"),
	}
}

// HotspotUserProfile handles hotspot user profile operations
type HotspotUserProfile struct {
	*BaseResource
}

// NewHotspotUserProfile creates a new HotspotUserProfile handler
func NewHotspotUserProfile(client *Client) *HotspotUserProfile {
	return &HotspotUserProfile{
		BaseResource: NewBaseResource(client, "/ip/hotspot/user/profile"),
	}
}

// HotspotUser handles hotspot user operations
type HotspotUser struct {
	*BaseResource
}

// NewHotspotUser creates a new HotspotUser handler
func NewHotspotUser(client *Client) *HotspotUser {
	return &HotspotUser{
		BaseResource: NewBaseResource(client, "/ip/hotspot/user"),
	}
}

// CreateUser creates a hotspot user
func (h *HotspotUser) CreateUser(name, password, profile string) (string, error) {
	return h.CreateUserContext(context.Background(), name, password, profile)
}

// CreateUserContext creates a hotspot user with context
func (h *HotspotUser) CreateUserContext(ctx context.Context, name, password, profile string) (string, error) {
	args := map[string]string{
		"name":     name,
		"password": password,
	}
	if profile != "" {
		args["profile"] = profile
	}
	return h.AddContext(ctx, args)
}

// ICMP handles ICMP operations
type ICMP struct {
	client *Client
}

// NewICMP creates a new ICMP handler
func NewICMP(client *Client) *ICMP {
	return &ICMP{client: client}
}

// Ping executes a ping command
func (i *ICMP) Ping(address string, count int) ([]*Response, error) {
	return i.PingContext(context.Background(), address, count)
}

// PingContext executes a ping command with context
func (i *ICMP) PingContext(ctx context.Context, address string, count int) ([]*Response, error) {
	return i.client.SendContext(ctx, "/ping", "=address="+address, "=count="+string(rune(count)))
}

// Traceroute executes a traceroute command
func (i *ICMP) Traceroute(address string) ([]*Response, error) {
	return i.TracerouteContext(context.Background(), address)
}

// TracerouteContext executes a traceroute command with context
func (i *ICMP) TracerouteContext(ctx context.Context, address string) ([]*Response, error) {
	return i.client.SendContext(ctx, "/tool/traceroute", "=address="+address)
}

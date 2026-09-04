package routeros

import (
	"context"
)

// IPv6 handles RouterOS IPv6 operations
type IPv6 struct {
	client *Client
}

// NewIPv6 creates a new IPv6 resource handler
func NewIPv6(client *Client) *IPv6 {
	return &IPv6{client: client}
}

// Address handles IPv6 address operations
type IPv6Address struct {
	*BaseResource
}

// NewIPv6Address creates a new IPv6Address handler
func NewIPv6Address(client *Client) *IPv6Address {
	return &IPv6Address{
		BaseResource: NewBaseResource(client, "/ipv6/address"),
	}
}

// AddAddress adds an IPv6 address to an interface
func (a *IPv6Address) AddAddress(address, interface_ string) (string, error) {
	return a.AddAddressContext(context.Background(), address, interface_)
}

// AddAddressContext adds an IPv6 address to an interface with context
func (a *IPv6Address) AddAddressContext(ctx context.Context, address, interface_ string) (string, error) {
	return a.AddContext(ctx, map[string]string{
		"address":   address,
		"interface": interface_,
	})
}

// Route handles IPv6 route operations
type IPv6Route struct {
	*BaseResource
}

// NewIPv6Route creates a new IPv6Route handler
func NewIPv6Route(client *Client) *IPv6Route {
	return &IPv6Route{
		BaseResource: NewBaseResource(client, "/ipv6/route"),
	}
}

// AddRoute adds an IPv6 route
func (r *IPv6Route) AddRoute(dstAddress, gateway string) (string, error) {
	return r.AddRouteContext(context.Background(), dstAddress, gateway)
}

// AddRouteContext adds an IPv6 route with context
func (r *IPv6Route) AddRouteContext(ctx context.Context, dstAddress, gateway string) (string, error) {
	return r.AddContext(ctx, map[string]string{
		"dst-address": dstAddress,
		"gateway":     gateway,
	})
}

// AddDefaultRoute adds a default IPv6 route
func (r *IPv6Route) AddDefaultRoute(gateway string) (string, error) {
	return r.AddDefaultRouteContext(context.Background(), gateway)
}

// AddDefaultRouteContext adds a default IPv6 route with context
func (r *IPv6Route) AddDefaultRouteContext(ctx context.Context, gateway string) (string, error) {
	return r.AddRouteContext(ctx, "::/0", gateway)
}

// ND handles IPv6 neighbor discovery operations
type IPv6ND struct {
	*BaseResource
}

// NewIPv6ND creates a new IPv6ND handler
func NewIPv6ND(client *Client) *IPv6ND {
	return &IPv6ND{
		BaseResource: NewBaseResource(client, "/ipv6/nd"),
	}
}

// Neighbor handles IPv6 neighbor operations
type IPv6Neighbor struct {
	*BaseResource
}

// NewIPv6Neighbor creates a new IPv6Neighbor handler
func NewIPv6Neighbor(client *Client) *IPv6Neighbor {
	return &IPv6Neighbor{
		BaseResource: NewBaseResource(client, "/ipv6/neighbor"),
	}
}

// DHCPClient handles IPv6 DHCP client operations
type IPv6DHCPClient struct {
	*BaseResource
}

// NewIPv6DHCPClient creates a new IPv6DHCPClient handler
func NewIPv6DHCPClient(client *Client) *IPv6DHCPClient {
	return &IPv6DHCPClient{
		BaseResource: NewBaseResource(client, "/ipv6/dhcp-client"),
	}
}

// Add creates a new IPv6 DHCP client
func (d *IPv6DHCPClient) Add(interface_ string) (string, error) {
	return d.AddContext(context.Background(), interface_)
}

// AddContext creates a new IPv6 DHCP client with context
func (d *IPv6DHCPClient) AddContext(ctx context.Context, interface_ string) (string, error) {
	return d.BaseResource.AddContext(ctx, map[string]string{
		"interface": interface_,
		"disabled":  "no",
	})
}

// DHCPServer handles IPv6 DHCP server operations
type IPv6DHCPServer struct {
	*BaseResource
}

// NewIPv6DHCPServer creates a new IPv6DHCPServer handler
func NewIPv6DHCPServer(client *Client) *IPv6DHCPServer {
	return &IPv6DHCPServer{
		BaseResource: NewBaseResource(client, "/ipv6/dhcp-server"),
	}
}

// Pool handles IPv6 pool operations
type IPv6Pool struct {
	*BaseResource
}

// NewIPv6Pool creates a new IPv6Pool handler
func NewIPv6Pool(client *Client) *IPv6Pool {
	return &IPv6Pool{
		BaseResource: NewBaseResource(client, "/ipv6/pool"),
	}
}

// CreatePool creates an IPv6 pool
func (p *IPv6Pool) CreatePool(name, prefix string, prefixLength int) (string, error) {
	return p.CreatePoolContext(context.Background(), name, prefix, prefixLength)
}

// CreatePoolContext creates an IPv6 pool with context
func (p *IPv6Pool) CreatePoolContext(ctx context.Context, name, prefix string, prefixLength int) (string, error) {
	return p.AddContext(ctx, map[string]string{
		"name":           name,
		"prefix":         prefix,
		"prefix-length":  string(rune(prefixLength)),
	})
}

// Firewall handles IPv6 firewall operations
type IPv6Firewall struct {
	client *Client
}

// NewIPv6Firewall creates a new IPv6Firewall handler
func NewIPv6Firewall(client *Client) *IPv6Firewall {
	return &IPv6Firewall{client: client}
}

// Filter handles IPv6 firewall filter rules
type IPv6FirewallFilter struct {
	*BaseResource
}

// NewIPv6FirewallFilter creates a new IPv6FirewallFilter handler
func NewIPv6FirewallFilter(client *Client) *IPv6FirewallFilter {
	return &IPv6FirewallFilter{
		BaseResource: NewBaseResource(client, "/ipv6/firewall/filter"),
	}
}

// AddRule adds an IPv6 firewall filter rule
func (f *IPv6FirewallFilter) AddRule(chain, action string, attrs map[string]string) (string, error) {
	return f.AddRuleContext(context.Background(), chain, action, attrs)
}

// AddRuleContext adds an IPv6 firewall filter rule with context
func (f *IPv6FirewallFilter) AddRuleContext(ctx context.Context, chain, action string, attrs map[string]string) (string, error) {
	args := map[string]string{
		"chain":  chain,
		"action": action,
	}
	for k, v := range attrs {
		args[k] = v
	}
	return f.AddContext(ctx, args)
}

// NAT handles IPv6 firewall NAT rules
type IPv6FirewallNAT struct {
	*BaseResource
}

// NewIPv6FirewallNAT creates a new IPv6FirewallNAT handler
func NewIPv6FirewallNAT(client *Client) *IPv6FirewallNAT {
	return &IPv6FirewallNAT{
		BaseResource: NewBaseResource(client, "/ipv6/firewall/nat"),
	}
}

// AddressList handles IPv6 firewall address lists
type IPv6FirewallAddressList struct {
	*BaseResource
}

// NewIPv6FirewallAddressList creates a new IPv6FirewallAddressList handler
func NewIPv6FirewallAddressList(client *Client) *IPv6FirewallAddressList {
	return &IPv6FirewallAddressList{
		BaseResource: NewBaseResource(client, "/ipv6/firewall/address-list"),
	}
}

// AddAddress adds an IPv6 address to a list
func (a *IPv6FirewallAddressList) AddAddress(list, address string) (string, error) {
	return a.AddAddressContext(context.Background(), list, address)
}

// AddAddressContext adds an IPv6 address to a list with context
func (a *IPv6FirewallAddressList) AddAddressContext(ctx context.Context, list, address string) (string, error) {
	return a.AddContext(ctx, map[string]string{
		"list":    list,
		"address": address,
	})
}

// RA handles IPv6 router advertisement operations
type IPv6RA struct {
	*BaseResource
}

// NewIPv6RA creates a new IPv6RA handler
func NewIPv6RA(client *Client) *IPv6RA {
	return &IPv6RA{
		BaseResource: NewBaseResource(client, "/ipv6/nd/prefix/default"),
	}
}

// Settings handles IPv6 settings
type IPv6Settings struct {
	client *Client
}

// NewIPv6Settings creates a new IPv6Settings handler
func NewIPv6Settings(client *Client) *IPv6Settings {
	return &IPv6Settings{client: client}
}

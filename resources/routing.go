package routeros

import (
	"context"
)

// Routing handles RouterOS routing operations
type Routing struct {
	client *Client
}

// NewRouting creates a new Routing resource handler
func NewRouting(client *Client) *Routing {
	return &Routing{client: client}
}

// Table handles routing table operations
type RoutingTable struct {
	*BaseResource
}

// NewRoutingTable creates a new RoutingTable handler
func NewRoutingTable(client *Client) *RoutingTable {
	return &RoutingTable{
		BaseResource: NewBaseResource(client, "/routing/table"),
	}
}

// Create creates a new routing table
func (r *RoutingTable) Create(name string) (string, error) {
	return r.CreateContext(context.Background(), name)
}

// CreateContext creates a new routing table with context
func (r *RoutingTable) CreateContext(ctx context.Context, name string) (string, error) {
	return r.AddContext(ctx, map[string]string{"name": name})
}

// Route handles routing operations
type RoutingRoute struct {
	*BaseResource
}

// NewRoutingRoute creates a new RoutingRoute handler
func NewRoutingRoute(client *Client) *RoutingRoute {
	return &RoutingRoute{
		BaseResource: NewBaseResource(client, "/routing/route"),
	}
}

// BGP handles BGP operations
type BGP struct {
	client *Client
}

// NewBGP creates a new BGP handler
func NewBGP(client *Client) *BGP {
	return &BGP{client: client}
}

// Instance handles BGP instance operations
type BGPInstance struct {
	*BaseResource
}

// NewBGPInstance creates a new BGPInstance handler
func NewBGPInstance(client *Client) *BGPInstance {
	return &BGPInstance{
		BaseResource: NewBaseResource(client, "/routing/bgp/instance"),
	}
}

// Create creates a new BGP instance
func (b *BGPInstance) Create(name string, as int, routerID string) (string, error) {
	return b.CreateContext(context.Background(), name, as, routerID)
}

// CreateContext creates a new BGP instance with context
func (b *BGPInstance) CreateContext(ctx context.Context, name string, as int, routerID string) (string, error) {
	return b.AddContext(ctx, map[string]string{
		"name":       name,
		"as":         string(rune(as)),
		"router-id":  routerID,
	})
}

// Peer handles BGP peer operations
type BGPPeer struct {
	*BaseResource
}

// NewBGPPeer creates a new BGPPeer handler
func NewBGPPeer(client *Client) *BGPPeer {
	return &BGPPeer{
		BaseResource: NewBaseResource(client, "/routing/bgp/peer"),
	}
}

// Create creates a new BGP peer
func (b *BGPPeer) Create(name, instance, remoteAddress string, remoteAs int) (string, error) {
	return b.CreateContext(context.Background(), name, instance, remoteAddress, remoteAs)
}

// CreateContext creates a new BGP peer with context
func (b *BGPPeer) CreateContext(ctx context.Context, name, instance, remoteAddress string, remoteAs int) (string, error) {
	return b.AddContext(ctx, map[string]string{
		"name":            name,
		"instance":        instance,
		"remote-address":  remoteAddress,
		"remote-as":       string(rune(remoteAs)),
	})
}

// OSPF handles OSPF operations
type OSPF struct {
	client *Client
}

// NewOSPF creates a new OSPF handler
func NewOSPF(client *Client) *OSPF {
	return &OSPF{client: client}
}

// OSPFInstance handles OSPF instance operations
type OSPFInstance struct {
	*BaseResource
}

// NewOSPFInstance creates a new OSPFInstance handler
func NewOSPFInstance(client *Client) *OSPFInstance {
	return &OSPFInstance{
		BaseResource: NewBaseResource(client, "/routing/ospf/instance"),
	}
}

// OSPFArea handles OSPF area operations
type OSPFArea struct {
	*BaseResource
}

// NewOSPFArea creates a new OSPFArea handler
func NewOSPFArea(client *Client) *OSPFArea {
	return &OSPFArea{
		BaseResource: NewBaseResource(client, "/routing/ospf/area"),
	}
}

// Create creates a new OSPF area
func (o *OSPFArea) Create(name, instance string, areaID string) (string, error) {
	return o.CreateContext(context.Background(), name, instance, areaID)
}

// CreateContext creates a new OSPF area with context
func (o *OSPFArea) CreateContext(ctx context.Context, name, instance string, areaID string) (string, error) {
	return o.AddContext(ctx, map[string]string{
		"name":     name,
		"instance": instance,
		"area-id":  areaID,
	})
}

// OSPFInterface handles OSPF interface operations
type OSPFInterface struct {
	*BaseResource
}

// NewOSPFInterface creates a new OSPFInterface handler
func NewOSPFInterface(client *Client) *OSPFInterface {
	return &OSPFInterface{
		BaseResource: NewBaseResource(client, "/routing/ospf/interface"),
	}
}

// OSPFNetwork handles OSPF network operations
type OSPFNetwork struct {
	*BaseResource
}

// NewOSPFNetwork creates a new OSPFNetwork handler
func NewOSPFNetwork(client *Client) *OSPFNetwork {
	return &OSPFNetwork{
		BaseResource: NewBaseResource(client, "/routing/ospf/network"),
	}
}

// RIP handles RIP operations
type RIP struct {
	client *Client
}

// NewRIP creates a new RIP handler
func NewRIP(client *Client) *RIP {
	return &RIP{client: client}
}

// RIPInstance handles RIP instance operations
type RIPInstance struct {
	*BaseResource
}

// NewRIPInstance creates a new RIPInstance handler
func NewRIPInstance(client *Client) *RIPInstance {
	return &RIPInstance{
		BaseResource: NewBaseResource(client, "/routing/rip/instance"),
	}
}

// Filter handles routing filter operations
type RoutingFilter struct {
	*BaseResource
}

// NewRoutingFilter creates a new RoutingFilter handler
func NewRoutingFilter(client *Client) *RoutingFilter {
	return &RoutingFilter{
		BaseResource: NewBaseResource(client, "/routing/filter"),
	}
}

// Rules handles routing rules operations
type RoutingRule struct {
	*BaseResource
}

// NewRoutingRule creates a new RoutingRule handler
func NewRoutingRule(client *Client) *RoutingRule {
	return &RoutingRule{
		BaseResource: NewBaseResource(client, "/routing/rule"),
	}
}

// VRF handles VRF operations
type VRF struct {
	*BaseResource
}

// NewVRF creates a new VRF handler
func NewVRF(client *Client) *VRF {
	return &VRF{
		BaseResource: NewBaseResource(client, "/routing/vrf"),
	}
}

// Create creates a new VRF
func (v *VRF) Create(name string, interfaces string) (string, error) {
	return v.CreateContext(context.Background(), name, interfaces)
}

// CreateContext creates a new VRF with context
func (v *VRF) CreateContext(ctx context.Context, name string, interfaces string) (string, error) {
	args := map[string]string{"name": name}
	if interfaces != "" {
		args["interfaces"] = interfaces
	}
	return v.AddContext(ctx, args)
}

// MME handles MME routing operations
type MME struct {
	*BaseResource
}

// NewMME creates a new MME handler
func NewMME(client *Client) *MME {
	return &MME{
		BaseResource: NewBaseResource(client, "/routing/mme"),
	}
}

// BFD handles BFD operations
type BFD struct {
	*BaseResource
}

// NewBFD creates a new BFD handler
func NewBFD(client *Client) *BFD {
	return &BFD{
		BaseResource: NewBaseResource(client, "/routing/bfd"),
	}
}

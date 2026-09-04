package routeros

import (
	"context"
)

// Interface handles RouterOS interface operations
type Interface struct {
	*BaseResource
	client *Client
}

// NewInterface creates a new Interface resource handler
func NewInterface(client *Client) *Interface {
	return &Interface{
		BaseResource: NewBaseResource(client, "/interface"),
		client:       client,
	}
}

// Ethernet handles ethernet interface operations
type Ethernet struct {
	*BaseResource
}

// NewEthernet creates a new Ethernet handler
func NewEthernet(client *Client) *Ethernet {
	return &Ethernet{
		BaseResource: NewBaseResource(client, "/interface/ethernet"),
	}
}

// SetMasterPort sets the master port for an ethernet interface
func (e *Ethernet) SetMasterPort(id, masterPort string) error {
	return e.SetMasterPortContext(context.Background(), id, masterPort)
}

// SetMasterPortContext sets the master port for an ethernet interface with context
func (e *Ethernet) SetMasterPortContext(ctx context.Context, id, masterPort string) error {
	return e.UpdateContext(ctx, id, map[string]string{"master-port": masterPort})
}

// Bridge handles bridge interface operations
type Bridge struct {
	*BaseResource
}

// NewBridge creates a new Bridge handler
func NewBridge(client *Client) *Bridge {
	return &Bridge{
		BaseResource: NewBaseResource(client, "/interface/bridge"),
	}
}

// BridgePort handles bridge port operations
type BridgePort struct {
	*BaseResource
}

// NewBridgePort creates a new BridgePort handler
func NewBridgePort(client *Client) *BridgePort {
	return &BridgePort{
		BaseResource: NewBaseResource(client, "/interface/bridge/port"),
	}
}

// AddPort adds an interface to a bridge
func (bp *BridgePort) AddPort(bridge, interface_ string) (string, error) {
	return bp.AddPortContext(context.Background(), bridge, interface_)
}

// AddPortContext adds an interface to a bridge with context
func (bp *BridgePort) AddPortContext(ctx context.Context, bridge, interface_ string) (string, error) {
	return bp.AddContext(ctx, map[string]string{
		"bridge":    bridge,
		"interface": interface_,
	})
}

// BridgeVLAN handles bridge VLAN operations
type BridgeVLAN struct {
	*BaseResource
}

// NewBridgeVLAN creates a new BridgeVLAN handler
func NewBridgeVLAN(client *Client) *BridgeVLAN {
	return &BridgeVLAN{
		BaseResource: NewBaseResource(client, "/interface/bridge/vlan"),
	}
}

// VLAN handles VLAN interface operations
type VLAN struct {
	*BaseResource
}

// NewVLAN creates a new VLAN handler
func NewVLAN(client *Client) *VLAN {
	return &VLAN{
		BaseResource: NewBaseResource(client, "/interface/vlan"),
	}
}

// Create creates a new VLAN interface
func (v *VLAN) Create(name, interface_ string, vlanID int) (string, error) {
	return v.CreateContext(context.Background(), name, interface_, vlanID)
}

// CreateContext creates a new VLAN interface with context
func (v *VLAN) CreateContext(ctx context.Context, name, interface_ string, vlanID int) (string, error) {
	return v.AddContext(ctx, map[string]string{
		"name":      name,
		"interface": interface_,
		"vlan-id":   string(rune(vlanID)),
	})
}

// VRRP handles VRRP interface operations
type VRRP struct {
	*BaseResource
}

// NewVRRP creates a new VRRP handler
func NewVRRP(client *Client) *VRRP {
	return &VRRP{
		BaseResource: NewBaseResource(client, "/interface/vrrp"),
	}
}

// Bonding handles bonding interface operations
type Bonding struct {
	*BaseResource
}

// NewBonding creates a new Bonding handler
func NewBonding(client *Client) *Bonding {
	return &Bonding{
		BaseResource: NewBaseResource(client, "/interface/bonding"),
	}
}

// PPPoEClient handles PPPoE client operations
type PPPoEClient struct {
	*BaseResource
}

// NewPPPoEClient creates a new PPPoEClient handler
func NewPPPoEClient(client *Client) *PPPoEClient {
	return &PPPoEClient{
		BaseResource: NewBaseResource(client, "/interface/pppoe-client"),
	}
}

// Connect establishes a PPPoE connection
func (p *PPPoEClient) Connect(id string) error {
	return p.ConnectContext(context.Background(), id)
}

// ConnectContext establishes a PPPoE connection with context
func (p *PPPoEClient) ConnectContext(ctx context.Context, id string) error {
	return p.UpdateContext(ctx, id, map[string]string{"disabled": "no"})
}

// Disconnect disconnects a PPPoE connection
func (p *PPPoEClient) Disconnect(id string) error {
	return p.DisconnectContext(context.Background(), id)
}

// DisconnectContext disconnects a PPPoE connection with context
func (p *PPPoEClient) DisconnectContext(ctx context.Context, id string) error {
	return p.UpdateContext(ctx, id, map[string]string{"disabled": "yes"})
}

// Wireless handles wireless interface operations
type Wireless struct {
	*BaseResource
}

// NewWireless creates a new Wireless handler
func NewWireless(client *Client) *Wireless {
	return &Wireless{
		BaseResource: NewBaseResource(client, "/interface/wireless"),
	}
}

// WirelessSecurityProfile handles wireless security profiles
type WirelessSecurityProfile struct {
	*BaseResource
}

// NewWirelessSecurityProfile creates a new WirelessSecurityProfile handler
func NewWirelessSecurityProfile(client *Client) *WirelessSecurityProfile {
	return &WirelessSecurityProfile{
		BaseResource: NewBaseResource(client, "/interface/wireless/security-profiles"),
	}
}

// List returns all interfaces
func (i *Interface) List() (*ResourceList, error) {
	return i.ListContext(context.Background())
}

// ListContext returns all interfaces with context
func (i *Interface) ListContext(ctx context.Context) (*ResourceList, error) {
	return i.BaseResource.ListContext(ctx)
}

// GetByName returns an interface by name
func (i *Interface) GetByName(name string) (*Resource, error) {
	return i.GetByNameContext(context.Background(), name)
}

// GetByNameContext returns an interface by name with context
func (i *Interface) GetByNameContext(ctx context.Context, name string) (*Resource, error) {
	return i.FindByNameContext(ctx, name)
}

// Monitor monitors an interface
func (i *Interface) Monitor(id string) (map[string]string, error) {
	return i.MonitorContext(context.Background(), id)
}

// MonitorContext monitors an interface with context
func (i *Interface) MonitorContext(ctx context.Context, id string) (map[string]string, error) {
	responses, err := i.client.SendContext(ctx, "/interface/monitor-traffic", "interface="+id, "once")
	if err != nil {
		return nil, err
	}
	parser := NewResponseParser(responses)
	return parser.FirstData(), nil
}

// TrafficMonitor monitors traffic on an interface
func (i *Interface) TrafficMonitor(iface string) (map[string]string, error) {
	return i.TrafficMonitorContext(context.Background(), iface)
}

// TrafficMonitorContext monitors traffic on an interface with context
func (i *Interface) TrafficMonitorContext(ctx context.Context, iface string) (map[string]string, error) {
	responses, err := i.client.SendContext(ctx, "/interface/monitor-traffic", "=interface="+iface, "=once=")
	if err != nil {
		return nil, err
	}
	parser := NewResponseParser(responses)
	return parser.FirstData(), nil
}

// SetMTU sets the MTU for an interface
func (i *Interface) SetMTU(id string, mtu int) error {
	return i.SetMTUContext(context.Background(), id, mtu)
}

// SetMTUContext sets the MTU for an interface with context
func (i *Interface) SetMTUContext(ctx context.Context, id string, mtu int) error {
	return i.UpdateContext(ctx, id, map[string]string{"mtu": string(rune(mtu))})
}

// SetMACAddress sets the MAC address for an interface
func (i *Interface) SetMACAddress(id, mac string) error {
	return i.SetMACAddressContext(context.Background(), id, mac)
}

// SetMACAddressContext sets the MAC address for an interface with context
func (i *Interface) SetMACAddressContext(ctx context.Context, id, mac string) error {
	return i.UpdateContext(ctx, id, map[string]string{"mac-address": mac})
}

// EnableInterface enables an interface
func (i *Interface) EnableInterface(name string) error {
	return i.EnableInterfaceContext(context.Background(), name)
}

// EnableInterfaceContext enables an interface with context
func (i *Interface) EnableInterfaceContext(ctx context.Context, name string) error {
	iface, err := i.FindByNameContext(ctx, name)
	if err != nil {
		return err
	}
	return i.EnableContext(ctx, iface.ID)
}

// DisableInterface disables an interface
func (i *Interface) DisableInterface(name string) error {
	return i.DisableInterfaceContext(context.Background())
}

// DisableInterfaceContext disables an interface with context
func (i *Interface) DisableInterfaceContext(ctx context.Context, name string) error {
	iface, err := i.FindByNameContext(ctx, name)
	if err != nil {
		return err
	}
	return i.DisableContext(ctx, iface.ID)
}

// L2TPClient handles L2TP client operations
type L2TPClient struct {
	*BaseResource
}

// NewL2TPClient creates a new L2TPClient handler
func NewL2TPClient(client *Client) *L2TPClient {
	return &L2TPClient{
		BaseResource: NewBaseResource(client, "/interface/l2tp-client"),
	}
}

// L2TPServer handles L2TP server operations
type L2TPServer struct {
	*BaseResource
}

// NewL2TPServer creates a new L2TPServer handler
func NewL2TPServer(client *Client) *L2TPServer {
	return &L2TPServer{
		BaseResource: NewBaseResource(client, "/interface/l2tp-server"),
	}
}

// SSTPClient handles SSTP client operations
type SSTPClient struct {
	*BaseResource
}

// NewSSTPClient creates a new SSTPClient handler
func NewSSTPClient(client *Client) *SSTPClient {
	return &SSTPClient{
		BaseResource: NewBaseResource(client, "/interface/sstp-client"),
	}
}

// OVVPNClient handles OpenVPN client operations
type OVVPNClient struct {
	*BaseResource
}

// NewOVVPNClient creates a new OVVPNClient handler
func NewOVVPNClient(client *Client) *OVVPNClient {
	return &OVVPNClient{
		BaseResource: NewBaseResource(client, "/interface/ovpn-client"),
	}
}

// WireGuard handles WireGuard interface operations
type WireGuard struct {
	*BaseResource
}

// NewWireGuard creates a new WireGuard handler
func NewWireGuard(client *Client) *WireGuard {
	return &WireGuard{
		BaseResource: NewBaseResource(client, "/interface/wireguard"),
	}
}

// WireGuardPeer handles WireGuard peer operations
type WireGuardPeer struct {
	*BaseResource
}

// NewWireGuardPeer creates a new WireGuardPeer handler
func NewWireGuardPeer(client *Client) *WireGuardPeer {
	return &WireGuardPeer{
		BaseResource: NewBaseResource(client, "/interface/wireguard/peers"),
	}
}

// Gre handles GRE tunnel operations
type Gre struct {
	*BaseResource
}

// NewGre creates a new Gre handler
func NewGre(client *Client) *Gre {
	return &Gre{
		BaseResource: NewBaseResource(client, "/interface/gre"),
	}
}

// IPIP handles IPIP tunnel operations
type IPIP struct {
	*BaseResource
}

// NewIPIP creates a new IPIP handler
func NewIPIP(client *Client) *IPIP {
	return &IPIP{
		BaseResource: NewBaseResource(client, "/interface/ipip"),
	}
}

// EoIP handles EoIP tunnel operations
type EoIP struct {
	*BaseResource
}

// NewEoIP creates a new EoIP handler
func NewEoIP(client *Client) *EoIP {
	return &EoIP{
		BaseResource: NewBaseResource(client, "/interface/eoip"),
	}
}

// VxLAN handles VxLAN interface operations
type VxLAN struct {
	*BaseResource
}

// NewVxLAN creates a new VxLAN handler
func NewVxLAN(client *Client) *VxLAN {
	return &VxLAN{
		BaseResource: NewBaseResource(client, "/interface/vxlan"),
	}
}

// MacVLAN handles MAC VLAN operations
type MacVLAN struct {
	*BaseResource
}

// NewMacVLAN creates a new MacVLAN handler
func NewMacVLAN(client *Client) *MacVLAN {
	return &MacVLAN{
		BaseResource: NewBaseResource(client, "/interface/mac-vlan"),
	}
}

// Loopback handles loopback interface operations
type Loopback struct {
	*BaseResource
}

// NewLoopback creates a new Loopback handler
func NewLoopback(client *Client) *Loopback {
	return &Loopback{
		BaseResource: NewBaseResource(client, "/interface/loopback"),
	}
}

// VirtualEthernet handles virtual ethernet operations
type VirtualEthernet struct {
	*BaseResource
}

// NewVirtualEthernet creates a new VirtualEthernet handler
func NewVirtualEthernet(client *Client) *VirtualEthernet {
	return &VirtualEthernet{
		BaseResource: NewBaseResource(client, "/interface/veth"),
	}
}

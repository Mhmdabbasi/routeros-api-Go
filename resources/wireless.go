package routeros

import (
	"context"
)

// Wireless handles RouterOS wireless operations
type Wireless struct {
	client *Client
}

// NewWireless creates a new Wireless resource handler
func NewWireless(client *Client) *Wireless {
	return &Wireless{client: client}
}

// AccessList handles wireless access list operations
type AccessList struct {
	*BaseResource
}

// NewWirelessAccessList creates a new AccessList handler
func NewWirelessAccessList(client *Client) *AccessList {
	return &AccessList{
		BaseResource: NewBaseResource(client, "/interface/wireless/access-list"),
	}
}

// AddMAC adds a MAC address to the access list
func (a *AccessList) AddMAC(mac string, allow bool) (string, error) {
	return a.AddMACContext(context.Background(), mac, allow)
}

// AddMACContext adds a MAC address to the access list with context
func (a *AccessList) AddMACContext(ctx context.Context, mac string, allow bool) (string, error) {
	action := "reject"
	if allow {
		action = "accept"
	}
	return a.AddContext(ctx, map[string]string{
		"mac-address": mac,
		"action":      action,
	})
}

// ConnectList handles wireless connect list operations
type ConnectList struct {
	*BaseResource
}

// NewConnectList creates a new ConnectList handler
func NewConnectList(client *Client) *ConnectList {
	return &ConnectList{
		BaseResource: NewBaseResource(client, "/interface/wireless/connect-list"),
	}
}

// SecurityProfile handles wireless security profile operations
type SecurityProfile struct {
	*BaseResource
}

// NewSecurityProfile creates a new SecurityProfile handler
func NewSecurityProfile(client *Client) *SecurityProfile {
	return &SecurityProfile{
		BaseResource: NewBaseResource(client, "/interface/wireless/security-profiles"),
	}
}

// Create creates a new security profile
func (s *SecurityProfile) Create(name, mode, authenticationTypes string) (string, error) {
	return s.CreateContext(context.Background(), name, mode, authenticationTypes)
}

// CreateContext creates a new security profile with context
func (s *SecurityProfile) CreateContext(ctx context.Context, name, mode, authenticationTypes string) (string, error) {
	return s.AddContext(ctx, map[string]string{
		"name":               name,
		"mode":               mode,
		"authentication-types": authenticationTypes,
	})
}

// SetWPA2 sets WPA2 authentication
func (s *SecurityProfile) SetWPA2(id, passphrase string) error {
	return s.SetWPA2Context(context.Background(), id, passphrase)
}

// SetWPA2Context sets WPA2 authentication with context
func (s *SecurityProfile) SetWPA2Context(ctx context.Context, id, passphrase string) error {
	return s.UpdateContext(ctx, id, map[string]string{
		"mode":               "dynamic-keys",
		"authentication-types": "wpa2-psk",
		"wpa2-pre-shared-key": passphrase,
	})
}

// RegistrationTable handles wireless registration table operations
type RegistrationTable struct {
	*BaseResource
}

// NewRegistrationTable creates a new RegistrationTable handler
func NewRegistrationTable(client *Client) *RegistrationTable {
	return &RegistrationTable{
		BaseResource: NewBaseResource(client, "/interface/wireless/registration-table"),
	}
}

// Cap handles CAP operations
type WirelessCAP struct {
	*BaseResource
}

// NewWirelessCAP creates a new WirelessCAP handler
func NewWirelessCAP(client *Client) *WirelessCAP {
	return &WirelessCAP{
		BaseResource: NewBaseResource(client, "/interface/wireless/cap"),
	}
}

// Sniffer handles wireless sniffer operations
type WirelessSniffer struct {
	client *Client
}

// NewWirelessSniffer creates a new WirelessSniffer handler
func NewWirelessSniffer(client *Client) *WirelessSniffer {
	return &WirelessSniffer{client: client}
}

// Start starts the wireless sniffer
func (s *WirelessSniffer) Start(interface_ string) error {
	return s.StartContext(context.Background(), interface_)
}

// StartContext starts the wireless sniffer with context
func (s *WirelessSniffer) StartContext(ctx context.Context, interface_ string) error {
	_, err := s.client.SendContext(ctx, "/interface/wireless/sniffer/start", "=interface="+interface_)
	return err
}

// Stop stops the wireless sniffer
func (s *WirelessSniffer) Stop() error {
	return s.StopContext(context.Background())
}

// StopContext stops the wireless sniffer with context
func (s *WirelessSniffer) StopContext(ctx context.Context) error {
	_, err := s.client.SendContext(ctx, "/interface/wireless/sniffer/stop")
	return err
}

// Align handles wireless align operations
type WirelessAlign struct {
	client *Client
}

// NewWirelessAlign creates a new WirelessAlign handler
func NewWirelessAlign(client *Client) *WirelessAlign {
	return &WirelessAlign{client: client}
}

// Snooper handles wireless snooper operations
type WirelessSnooper struct {
	client *Client
}

// NewWirelessSnooper creates a new WirelessSnooper handler
func NewWirelessSnooper(client *Client) *WirelessSnooper {
	return &WirelessSnooper{client: client}
}

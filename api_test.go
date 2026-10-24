package routeros

import "testing"

func TestNetworkAddressFormatsDecimalPort(t *testing.T) {
	tests := []struct {
		host string
		port int
		want string
	}{
		{host: "192.168.88.1", port: 8728, want: "192.168.88.1:8728"},
		{host: "::1", port: 8729, want: "[::1]:8729"},
	}
	for _, tt := range tests {
		if got := NetworkAddress(tt.host, tt.port); got != tt.want {
			t.Errorf("NetworkAddress(%q, %d) = %q, want %q", tt.host, tt.port, got, tt.want)
		}
	}
}

func TestNewWithConfigInitializesRootResources(t *testing.T) {
	api, err := NewWithConfig(NewConfig().
		WithAddress("127.0.0.1").
		WithPort(DefaultPort).
		WithCredentials("test", "test"))
	if err != nil {
		t.Fatalf("NewWithConfig() error = %v", err)
	}
	if api.System == nil || api.Interface == nil || api.IP == nil || api.IPv6 == nil ||
		api.PPP == nil || api.Queue == nil || api.Routing == nil || api.Tools == nil ||
		api.Certificate == nil || api.Wireless == nil || api.Log == nil {
		t.Fatal("expected all root resource handles to be initialized")
	}
}

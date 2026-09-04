# RouterOS API Go

[![GoDoc](https://godoc.org/github.com/Mhmdabbasi/routeros-api-Go?status.svg)](https://godoc.org/github.com/Mhmdabbasi/routeros-api-Go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Mhmdabbasi/routeros-api-Go)](https://goreportcard.com/report/github.com/Mhmdabbasi/routeros-api-Go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive, production-ready Go client library for MikroTik RouterOS API. This library provides a fluent, type-safe interface for managing MikroTik routers programmatically.

## Features

- 🔌 **Easy Connection** - Simple API for connecting to RouterOS devices
- 🔐 **TLS Support** - Secure connections with TLS/SSL
- 📝 **Fluent Query Builder** - Build complex queries with a chainable interface
- 🔄 **Auto-Reconnect** - Automatic reconnection on connection loss
- ⏱️ **Context Support** - Full context support for timeout and cancellation
- 📊 **Resource Management** - Comprehensive resource handlers for all RouterOS sections
- 🎯 **Type-Safe** - Strong typing with Go structs and interfaces
- 🔧 **Configurable** - Extensive configuration options
- 📚 **Well Documented** - Comprehensive documentation and examples
- ✅ **All RouterOS Versions** - Support for all RouterOS versions from v6.x to v7.x

## Installation

```bash
go get github.com/Mhmdabbasi/routeros-api-Go
```

## Quick Start

### Basic Connection

```go
package main

import (
    "fmt"
    "log"
    
    routeros "github.com/Mhmdabbasi/routeros-api-Go"
)

func main() {
    // Quick connect
    api, err := routeros.QuickConnect("192.168.88.1", "admin", "password")
    if err != nil {
        log.Fatal(err)
    }
    defer api.Close()
    
    // Get system info
    info, err := api.System.GetInfo()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Router: %s\n", info.Identity)
    fmt.Printf("Version: %s\n", info.Version)
}
```

### TLS Connection

```go
config := routeros.NewConfig().
    WithAddress("192.168.88.1").
    WithCredentials("admin", "password").
    WithTLS().
    WithSkipVerify(true) // For self-signed certificates

api, err := routeros.NewWithConfig(config)
if err != nil {
    log.Fatal(err)
}

if err := api.Connect(); err != nil {
    log.Fatal(err)
}
defer api.Close()
```

## Configuration

The library provides extensive configuration options:

```go
config := routeros.NewConfig().
    WithAddress("192.168.88.1").
    WithPort(8728).
    WithCredentials("admin", "password").
    WithTimeout(30 * time.Second).
    WithTLS().
    WithSkipVerify(false).
    WithCACert("/path/to/ca.crt").
    WithClientCert("/path/to/client.crt", "/path/to/client.key").
    WithAutoReconnect(true).
    WithRetry(3, 5 * time.Second).
    WithKeepAlive(true, 30 * time.Second).
    WithBufferSize(4096).
    WithLogLevel(routeros.LogLevelDebug)
```

## Query Builder

Build complex queries using the fluent interface:

```go
// Simple query
query := routeros.Command("/interface/print").
    Where("name", routeros.QueryEqual, "ether1").
    PropList("name", "type", "mac-address")

responses, err := query.Execute(client)

// Complex query with multiple conditions
query := routeros.Command("/ip/firewall/filter/print").
    Where("chain", routeros.QueryEqual, "input").
    Where("action", routeros.QueryEqual, "accept").
    And()

// Range query
query := routeros.Command("/ip/address/print").
    Where("address", routeros.QueryGreater, "192.168.1.1").
    Where("address", routeros.QueryLess, "192.168.1.100").
    And()
```

## Resource Management

### Interfaces

```go
// List all interfaces
ifaces, err := api.Interface.List()
for _, iface := range ifaces.All() {
    fmt.Printf("%s: %s\n", iface.Name, iface.Get("type"))
}

// Get interface by name
iface, err := api.Interface.GetByName("ether1")

// Enable/disable interface
api.Interface.EnableInterface("ether1")
api.Interface.DisableInterface("ether1")

// Monitor interface traffic
traffic, err := api.Interface.TrafficMonitor("ether1")
```

### IP Addresses

```go
ipAddr := routeros.NewIPAddress(api.Client())

// Add IP address
id, err := ipAddr.AddAddress("192.168.90.1/24", "ether2")

// List addresses
list, err := ipAddr.List()

// Find by interface
addrs, err := ipAddr.Find(map[string]string{"interface": "ether1"})
```

### Routes

```go
route := routeros.NewRoute(api.Client())

// Add default route
id, err := route.AddDefaultRoute("192.168.88.1")

// Add static route
id, err := route.AddRoute("10.0.0.0/8", "192.168.88.1")

// List routes
routes, err := route.List()
```

### DNS

```go
dns := routeros.NewDNS(api.Client())

// Set DNS servers
err := dns.SetServers("8.8.8.8,8.8.4.4")

// Get DNS config
config, err := dns.Get()

// Flush DNS cache
err := dns.FlushCache()
```

### DHCP Server

```go
// DHCP Server
server := routeros.NewDHCPServer(api.Client())
servers, err := server.List()

// DHCP Leases
leases := routeros.NewDHCPLease(api.Client())
leaseList, err := leases.List()

// Add static lease
id, err := leases.AddStatic("AA:BB:CC:DD:EE:FF", "192.168.88.100")

// DHCP Networks
networks := routeros.NewDHCPNetwork(api.Client())
id, err := networks.AddNetwork("192.168.88.0/24", "192.168.88.1")
```

### Firewall

```go
// Filter rules
filter := routeros.NewFirewallFilter(api.Client())

// Add filter rule
id, err := filter.AddRule("input", "accept", map[string]string{
    "protocol": "tcp",
    "dst-port": "22",
    "comment":  "Allow SSH",
})

// NAT rules
nat := routeros.NewFirewallNAT(api.Client())

// Add masquerade
id, err := nat.AddMasquerade("srcnat", "ether1-wan")

// Add port forwarding
id, err := nat.AddDstNat("dstnat", "1.2.3.4", "192.168.88.10", "80", "80")

// Address lists
addrList := routeros.NewFirewallAddressList(api.Client())
id, err := addrList.AddAddress("blacklist", "10.0.0.1")
```

### PPP

```go
// PPP Secrets (Users)
secret := routeros.NewPPPSecret(api.Client())

// Create PPP user
id, err := secret.CreateUser("user1", "password123", "default", "pppoe")

// List users
users, err := secret.List()

// Active sessions
active := routeros.NewPPPActive(api.Client())
sessions, err := active.List()

// Disconnect session
err := active.Disconnect(sessionID)
```

### Queues

```go
// Simple queues
queue := routeros.NewQueueSimple(api.Client())

// Create queue
id, err := queue.Create("user-queue", "192.168.88.10/32", "10M/10M")

// Update queue
err := queue.SetMaxLimit(id, "20M/20M")

// Queue tree
tree := routeros.NewQueueTree(api.Client())
id, err := tree.Create("upload", "global", "50M")
```

### System Operations

```go
// System info
info, err := api.System.GetInfo()

// Reboot
err := api.System.Reboot()

// Shutdown
err := api.System.Shutdown()

// Identity
identity := routeros.NewIdentity(api.Client())
err := identity.Set("MyRouter")

// Clock
clock := routeros.NewClock(api.Client())
time, err := clock.Time()
err := clock.SetTimeZone("Asia/Dubai")

// Backup
backup := routeros.NewBackup(api.Client())
err := backup.Save("my-backup")
err := backup.Load("my-backup")

// Scripts
script := routeros.NewScript(api.Client())
err := script.Run("my-script")

// Scheduler
scheduler := routeros.NewScheduler(api.Client())
id, err := scheduler.Add(map[string]string{
    "name":      "daily-backup",
    "start-time": "00:00:00",
    "interval":   "1d",
    "on-event":   "/system/backup/save",
})

// Users
user := routeros.NewUser(api.Client())
id, err := user.Add(map[string]string{
    "name":     "operator",
    "password": "secure123",
    "group":    "read",
})
```

### Certificates

```go
cert := routeros.NewCertificate(api.Client())

// Create CA
id, err := cert.CreateCA("my-ca", 3650)

// Create server certificate
id, err := cert.CreateServer("server-cert", 365)

// Sign certificate
err := cert.Sign(id, "my-ca", 365)

// Export certificate
err := cert.Export(id, "export-password")
```

### Wireless

```go
// Wireless interfaces
wireless := routeros.NewWireless(api.Client())

// Security profiles
profile := routeros.NewSecurityProfile(api.Client())
id, err := profile.Create("secure-profile", "dynamic-keys", "wpa2-psk")

// Access lists
accessList := routeros.NewWirelessAccessList(api.Client())
id, err := accessList.AddMAC("AA:BB:CC:DD:EE:FF", true)

// Registration table (connected clients)
regTable := routeros.NewRegistrationTable(api.Client())
clients, err := regTable.List()
```

### IPv6

```go
// IPv6 addresses
ipv6Addr := routeros.NewIPv6Address(api.Client())
id, err := ipv6Addr.AddAddress("2001:db8::1/64", "ether1")

// IPv6 routes
ipv6Route := routeros.NewIPv6Route(api.Client())
id, err := ipv6Route.AddRoute("2001:db8::/32", "fe80::1")

// IPv6 DHCP
ipv6Dhcp := routeros.NewIPv6DHCPClient(api.Client())
id, err := ipv6Dhcp.Add("ether1")
```

### Routing

```go
// BGP
bgpInstance := routeros.NewBGPInstance(api.Client())
id, err := bgpInstance.Create("bgp-1", 65001, "192.168.88.1")

bgpPeer := routeros.NewBGPPeer(api.Client())
id, err := bgpPeer.Create("peer-1", "bgp-1", "10.0.0.1", 65002)

// OSPF
ospfInstance := routeros.NewOSPFInstance(api.Client())
ospfArea := routeros.NewOSPFArea(api.Client())
id, err := ospfArea.Create("area-0", "ospf-1", "0.0.0.0")

// VRF
vrf := routeros.NewVRF(api.Client())
id, err := vrf.Create("vrf-1", "ether2,ether3")
```

### Tools

```go
tools := routeros.NewTools(api.Client())

// Ping
responses, err := tools.Ping("8.8.8.8", 4)

// Traceroute
responses, err := tools.Traceroute("8.8.8.8")

// Bandwidth test
responses, err := tools.BandwidthTest("192.168.88.10", "both")

// Packet sniffer
sniffer := routeros.NewSniffer(api.Client())
err := sniffer.SetFilter("ether1", "", "", "80")
err = sniffer.Start()
// ... wait ...
err = sniffer.Stop()
packets, err := sniffer.Packets()

// Netwatch
netwatch := routeros.NewNetwatch(api.Client())
id, err := netwatch.Create("8.8.8.8", "1m")
```

## Context Support

All operations support context for timeout and cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

responses, err := client.SendContext(ctx, "/interface/print")

// Or with resources
list, err := iface.ListContext(ctx)
```

## Error Handling

```go
responses, err := client.Send("/interface/print")
if err != nil {
    if errors.Is(err, routeros.ErrNotConnected) {
        // Handle connection error
    } else if errors.Is(err, routeros.ErrAuthenticationFailed) {
        // Handle auth error
    } else if errors.Is(err, routeros.ErrTimeout) {
        // Handle timeout
    }
    
    // Check for router errors
    var routerErr *routeros.RouterError
    if errors.As(err, &routerErr) {
        fmt.Printf("Router error: %s (category: %d)\n", 
            routerErr.Message, routerErr.Category)
    }
}
```

## Response Parsing

```go
responses, err := client.Send("/interface/print")
if err != nil {
    log.Fatal(err)
}

parser := routeros.NewResponseParser(responses)

// Get count
count := parser.Count()

// Get first item
first := parser.First()
fmt.Printf("First: %s\n", first.Get("name"))

// Get all data
for _, data := range parser.Data() {
    fmt.Printf("%s: %s\n", data["name"], data["mac-address"])
}

// Convert to JSON
json, err := parser.ToJSONIndent()
fmt.Println(string(json))

// Get specific value
name := parser.Value("name")

// Parse typed values
mtu, err := parser.ValueInt("mtu")
enabled, err := parser.ValueBool("disabled")
```

## Raw Commands

You can also send raw commands:

```go
// Simple command
responses, err := client.Send("/system/resource/print")

// With arguments
responses, err := client.Send("/ip/address/add", 
    "=address=192.168.90.1/24", 
    "=interface=ether2",
)

// With query
responses, err := client.Send("/interface/print",
    "?name=ether1",
    "?type=ethernet",
    "?#&", // AND
)
```

## Best Practices

1. **Always close connections**: Use `defer api.Close()` to ensure proper cleanup
2. **Use context**: Set appropriate timeouts for all operations
3. **Handle errors**: Check for both connection and router errors
4. **Use TLS**: Always use TLS for connections over untrusted networks
5. **Pool connections**: Reuse API instances when possible

## Compatibility

| RouterOS Version | Support |
|-----------------|---------|
| v6.x            | ✅ Full |
| v7.x            | ✅ Full |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [go-routeros](https://github.com/go-routeros/routeros)
- [MikroTik RouterOS API Documentation](https://help.mikrotik.com/docs/display/ROS/API)
- [routeros-api-php](https://github.com/EvilFreelancer/routeros-api-php)

## Author

[Mhmdabbasi](https://github.com/Mhmdabbasi)

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	routeros "github.com/Mhmdabbasi/routeros-api-Go"
)

func main() {
	// Example 1: Simple connection
	fmt.Println("=== Example 1: Simple Connection ===")
	exampleSimpleConnection()

	// Example 2: TLS connection
	fmt.Println("\n=== Example 2: TLS Connection ===")
	exampleTLSConnection()

	// Example 3: Using query builder
	fmt.Println("\n=== Example 3: Query Builder ===")
	exampleQueryBuilder()

	// Example 4: Working with interfaces
	fmt.Println("\n=== Example 4: Interface Management ===")
	exampleInterfaceManagement()

	// Example 5: IP address management
	fmt.Println("\n=== Example 5: IP Address Management ===")
	exampleIPAddressManagement()

	// Example 6: Firewall rules
	fmt.Println("\n=== Example 6: Firewall Management ===")
	exampleFirewallManagement()

	// Example 7: PPP users
	fmt.Println("\n=== Example 7: PPP User Management ===")
	examplePPPUsers()

	// Example 8: Queue management
	fmt.Println("\n=== Example 8: Queue Management ===")
	exampleQueueManagement()

	// Example 9: System information
	fmt.Println("\n=== Example 9: System Information ===")
	exampleSystemInfo()

	// Example 10: Using context with timeout
	fmt.Println("\n=== Example 10: Context with Timeout ===")
	exampleContextTimeout()
}

// Example 1: Simple connection
func exampleSimpleConnection() {
	// Create API client
	api, err := routeros.New("192.168.88.1", "admin", "password")
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}

	// Connect to router
	if err := api.Connect(); err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer api.Close()

	fmt.Println("Connected successfully!")

	// Get router identity
	identity := routeros.NewIdentity(api.Client())
	name, err := identity.Get()
	if err != nil {
		log.Printf("Failed to get identity: %v", err)
		return
	}
	fmt.Printf("Router identity: %s\n", name)
}

// Example 2: TLS connection
func exampleTLSConnection() {
	config := routeros.NewConfig().
		WithAddress("192.168.88.1").
		WithCredentials("admin", "password").
		WithTLS().
		WithSkipVerify(true) // Skip certificate verification for self-signed certs

	api, err := routeros.NewWithConfig(config)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}

	if err := api.Connect(); err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer api.Close()

	fmt.Println("TLS connection established!")
}

// Example 3: Query builder
func exampleQueryBuilder() {
	api, err := routeros.QuickConnect("192.168.88.1", "admin", "password")
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer api.Close()

	// Build a query using fluent interface
	query := routeros.Command("/interface/print").
		Where("name", routeros.QueryEqual, "ether1").
		PropList("name", "type", "mac-address", "mtu")

	responses, err := query.Execute(api.Client())
	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}

	parser := routeros.NewResponseParser(responses)
	fmt.Printf("Found %d interfaces\n", parser.Count())
	for _, data := range parser.Data() {
		fmt.Printf("  - %s (%s)\n", data["name"], data["type"])
	}
}

// Example 4: Interface management
func exampleInterfaceManagement() {
	api, err := routeros.QuickConnect("192.168.88.1", "admin", "password")
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer api.Close()

	// Get all interfaces
	iface := routeros.NewInterface(api.Client())
	list, err := iface.List()
	if err != nil {
		log.Printf("Failed to get interfaces: %v", err)
		return
	}

	fmt.Println("Interfaces:")
	for _, item := range list.All() {
		fmt.Printf("  - %s (%s): %s\n", item.Name, item.Get("type"), item.Get("mac-address"))
	}

	// Get ethernet interfaces
	ethernet := routeros.NewEthernet(api.Client())
	ethList, err := ethernet.List()
	if err != nil {
		log.Printf("Failed to get ethernet interfaces: %v", err)
		return
	}

	fmt.Println("Ethernet interfaces:")
	for _, item := range ethList.All() {
		fmt.Printf("  - %s\n", item.Name)
	}
}

// Example 5: IP address management
func exampleIPAddressManagement() {
	api, err := routeros.QuickConnect("192.168.88.1", "admin", "password")
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer api.Close()

	// Get all IP addresses
	ipAddr := routeros.NewIPAddress(api.Client())
	list, err := ipAddr.List()
	if err != nil {
		log.Printf("Failed to get IP addresses: %v", err)
		return
	}

	fmt.Println("IP Addresses:")
	for _, item := range list.All() {
		fmt.Printf("  - %s on %s\n", item.Get("address"), item.Get("interface"))
	}

	// Add new IP address
	// id, err := ipAddr.AddAddress("192.168.90.1/24", "ether2")
	// if err != nil {
	// 	log.Printf("Failed to add IP: %v", err)
	// 	return
	// }
	// fmt.Printf("Added IP address with ID: %s\n", id)
}

// Example 6: Firewall management
func exampleFirewallManagement() {
	api, err := routeros.QuickConnect("192.168.88.1", "admin", "password")
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer api.Close()

	// Get filter rules
	filter := routeros.NewFirewallFilter(api.Client())
	list, err := filter.List()
	if err != nil {
		log.Printf("Failed to get firewall rules: %v", err)
		return
	}

	fmt.Println("Firewall Filter Rules:")
	for _, item := range list.All() {
		fmt.Printf("  - Chain: %s, Action: %s\n", item.Get("chain"), item.Get("action"))
	}

	// Add a new firewall rule
	// id, err := filter.AddRule("input", "accept", map[string]string{
	// 	"protocol": "tcp",
	// 	"dst-port": "22",
	// 	"comment":  "Allow SSH",
	// })
	// if err != nil {
	// 	log.Printf("Failed to add rule: %v", err)
	// 	return
	// }
	// fmt.Printf("Added rule with ID: %s\n", id)

	// Get NAT rules
	nat := routeros.NewFirewallNAT(api.Client())
	natList, err := nat.List()
	if err != nil {
		log.Printf("Failed to get NAT rules: %v", err)
		return
	}

	fmt.Println("NAT Rules:")
	for _, item := range natList.All() {
		fmt.Printf("  - Chain: %s, Action: %s\n", item.Get("chain"), item.Get("action"))
	}
}

// Example 7: PPP user management
func examplePPPUsers() {
	api, err := routeros.QuickConnect("192.168.88.1", "admin", "password")
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer api.Close()

	// Get PPP secrets (users)
	secret := routeros.NewPPPSecret(api.Client())
	list, err := secret.List()
	if err != nil {
		log.Printf("Failed to get PPP users: %v", err)
		return
	}

	fmt.Println("PPP Users:")
	for _, item := range list.All() {
		fmt.Printf("  - %s (Profile: %s)\n", item.Name, item.Get("profile"))
	}

	// Create a new PPP user
	// id, err := secret.CreateUser("user1", "password123", "default", "pppoe")
	// if err != nil {
	// 	log.Printf("Failed to create user: %v", err)
	// 	return
	// }
	// fmt.Printf("Created user with ID: %s\n", id)

	// Get active sessions
	active := routeros.NewPPPActive(api.Client())
	activeList, err := active.List()
	if err != nil {
		log.Printf("Failed to get active sessions: %v", err)
		return
	}

	fmt.Println("Active Sessions:")
	for _, item := range activeList.All() {
		fmt.Printf("  - %s from %s\n", item.Get("name"), item.Get("caller-id"))
	}
}

// Example 8: Queue management
func exampleQueueManagement() {
	api, err := routeros.QuickConnect("192.168.88.1", "admin", "password")
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer api.Close()

	// Get simple queues
	queue := routeros.NewQueueSimple(api.Client())
	list, err := queue.List()
	if err != nil {
		log.Printf("Failed to get queues: %v", err)
		return
	}

	fmt.Println("Simple Queues:")
	for _, item := range list.All() {
		fmt.Printf("  - %s: target=%s, max-limit=%s\n",
			item.Name, item.Get("target"), item.Get("max-limit"))
	}

	// Create a new queue
	// id, err := queue.Create("queue1", "192.168.88.10/32", "10M/10M")
	// if err != nil {
	// 	log.Printf("Failed to create queue: %v", err)
	// 	return
	// }
	// fmt.Printf("Created queue with ID: %s\n", id)
}

// Example 9: System information
func exampleSystemInfo() {
	api, err := routeros.QuickConnect("192.168.88.1", "admin", "password")
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer api.Close()

	// Get system info
	info, err := api.System.GetInfo()
	if err != nil {
		log.Printf("Failed to get system info: %v", err)
		return
	}

	fmt.Printf("Router: %s\n", info.Identity)
	fmt.Printf("Version: %s\n", info.Version)
	fmt.Printf("Board: %s\n", info.BoardName)
	fmt.Printf("Model: %s\n", info.Model)
	fmt.Printf("Architecture: %s\n", info.ArchName)
	fmt.Printf("Uptime: %s\n", info.Uptime)
	fmt.Printf("CPU: %s (%s cores)\n", info.CPU, info.CPUCount)
	fmt.Printf("Memory: %s / %s\n", info.FreeMemory, info.TotalMemory)

	// Get resource usage
	resource, err := api.System.Resource()
	if err != nil {
		log.Printf("Failed to get resource: %v", err)
		return
	}
	fmt.Printf("CPU Load: %s%%\n", resource.Get("cpu-load"))
	fmt.Printf("Free HDD: %s\n", resource.Get("free-hdd-space"))

	// Get packages
	packages, err := api.System.Packages()
	if err != nil {
		log.Printf("Failed to get packages: %v", err)
		return
	}
	fmt.Println("Installed Packages:")
	for _, pkg := range packages.All() {
		fmt.Printf("  - %s: %s\n", pkg.Get("name"), pkg.Get("version"))
	}
}

// Example 10: Using context with timeout
func exampleContextTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := routeros.NewConfig().
		WithAddress("192.168.88.1").
		WithCredentials("admin", "password").
		WithTimeout(5 * time.Second)

	client := routeros.NewClient(config)
	if err := client.Connect(); err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer client.Close()

	// Send command with context
	responses, err := client.SendContext(ctx, "/interface/print")
	if err != nil {
		log.Printf("Command failed: %v", err)
		return
	}

	parser := routeros.NewResponseParser(responses)
	fmt.Printf("Found %d interfaces\n", parser.Count())
}

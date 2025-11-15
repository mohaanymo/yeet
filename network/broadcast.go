package network

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/grandcat/zeroconf"
)

// RECEIVER MODE: Announces "I'm ready to receive files!"
func ReceiverMode() {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "Unknown"
	}

	fmt.Printf("📥 %s is waiting to receive files...\n", hostname)
	fmt.Println("Visible to senders on the network")
	fmt.Println("Press Ctrl+C to stop\n")

	// Announce that we're ready to receive
	server, err := zeroconf.Register(
		hostname,           // Computer name
		"_yeet._tcp",  // Service type
		"local.",           // Local network
		9090,               // Port where we'll receive
		[]string{"status=waiting"}, // Metadata
		nil,
	)

	if err != nil {
		log.Fatal("Failed to start receiver:", err)
	}
	defer server.Shutdown()

	// Keep running until Ctrl+C
	select {}
}

// SENDER MODE: Discovers receivers and lets user choose
func SenderMode(filename string) {
	fmt.Printf("📤 Looking for devices ready to receive '%s'...\n\n", filename)

	// Channel to collect discovered receivers
	results := make(chan *zeroconf.ServiceEntry)

	// Search for 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start searching for receivers
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = resolver.Browse(ctx, "_yeet._tcp", "local.", results)
	if err != nil {
		log.Fatal(err)
	}

	// Collect all found receivers
	receivers := []*zeroconf.ServiceEntry{}

	go func() {
		for entry := range results {
			receivers = append(receivers, entry)
			fmt.Printf("[%d] 💻 %s (%s:%d)\n", 
				len(receivers), 
				entry.Instance, 
				entry.AddrIPv4[0], 
				entry.Port)
		}
	}()

	// Wait for search to complete
	<-ctx.Done()

	if len(receivers) == 0 {
		fmt.Println("\n❌ No devices found waiting to receive files")
		fmt.Println("Make sure the receiver is running on another device")
		return
	}

	// Let user choose which receiver
	fmt.Printf("\nFound %d device(s). Which one do you want to send to?\n", len(receivers))
	fmt.Print("Enter number (or 0 to cancel): ")

	var choice int
	fmt.Scan(&choice)

	if choice < 1 || choice > len(receivers) {
		fmt.Println("Cancelled")
		return
	}

	selected := receivers[choice-1]
	fmt.Printf("\n✅ Selected: %s\n", selected.Instance)
	fmt.Printf("📍 IP: %s\n", selected.AddrIPv4[0])
	fmt.Printf("🔌 Port: %d\n", selected.Port)
	fmt.Println("\n🎉 Ready to connect and send file!")
	fmt.Printf("Connect to: %s:%d\n", selected.AddrIPv4[0], selected.Port)
	
	// TODO: Here you would implement your file transfer
	// For example: http.Post(fmt.Sprintf("http://%s:%d/upload", ip, port), ...)
}


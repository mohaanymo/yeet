package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/mohaanymo/yeet/protocol"
)

// SENDER MODE: Discovers receivers and lets user choose
func SenderMode(filePath string) {
    filename := filepath.Base(filePath)
    results := make(chan *zeroconf.ServiceEntry)
    
    // Add timeout for discovery
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    
    resolver, err := zeroconf.NewResolver(nil)
    if err != nil {
        log.Fatal(err)
        return
    }
    
    fmt.Println("Searching for devices to connect...")
    
    // Start browsing in a goroutine
    err = resolver.Browse(ctx, SERVICENAME+"._tcp", "local.", results)
    if err != nil {
        log.Fatal(err)
        return
    }
    
    receivers := []*zeroconf.ServiceEntry{}
    var target *zeroconf.ServiceEntry
    
    // collect results
    go func() {
        for entry := range results {
            receivers = append(receivers, entry)
        }
    }()
    
    // Give some time to discover devices
    time.Sleep(3 * time.Second)
    
    // Interactive selection loop
    for {
        if len(receivers) == 0 {
            fmt.Println("No devices found yet. Waiting...")
            time.Sleep(2 * time.Second)
            continue
        }

		for i, recv := range receivers {
			fmt.Printf("\n[%d] %s (%s:%d)\n",
                i+1,
                recv.Instance,
                recv.AddrIPv4[0],
                recv.Port)
		}
		
        
        fmt.Printf("\nFound %d device(s). Choose index (1-%d) or -1 to wait for more: ", 
            len(receivers), len(receivers))
        
        var choice int
        fmt.Scan(&choice)
        
        if choice == -1 {
            fmt.Println("Continuing to search...")
            time.Sleep(2 * time.Second)
            continue
        }
        
        if choice >= 1 && choice <= len(receivers) {
            target = receivers[choice-1]
            break
        }
        
        fmt.Println("Invalid choice!")
    }
    
    // Cancel discovery once target is selected
    cancel()
    
    if target == nil {
        log.Fatal("No target selected")
        return
    }
    
    addr := fmt.Sprintf("%s:%d", target.AddrIPv4[0], target.Port)
    conn, err := ConnectTo(addr)
    if err != nil {
        log.Fatalf("Cannot connect to the device: %v", err)
        return
    }
    defer conn.Close()

    fmt.Printf("Connected to %s, sending file...\n", target.HostName)
    
    file, err := protocol.NewFile(filename, filePath)
    if err != nil {
        log.Fatal(err)
        return
    }
    
    SendFile(file, conn)
    
    fmt.Println("File sent successfully!")
    
}

// RECEIVER MODE: Announces "I'm ready to receive files!"
func ReceiverMode() {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "Unknown"
	}

	fmt.Printf("%s is waiting to receive files...\n", hostname)

	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatal(err)
	}

	// Announce that we're ready to receive
	server, err := zeroconf.Register(
		hostname,                   // Computer name
		SERVICENAME+"._tcp",        // Service type
		"local.",                   // Local network
		PORT,                       // Port where we'll receive
		[]string{"status=waiting"}, // Metadata
		nil,
	)

	if err != nil {
		log.Fatal("Failed to start receiver:", err)
	}

	defer server.Shutdown()
	defer listener.Close()

	fmt.Printf("Listening on port %d...\n", PORT)
	AcceptConnections(listener)
	fmt.Println("File received successfully!")

}

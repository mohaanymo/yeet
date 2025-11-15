package peer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

const (
	PORT = 8080
	DOMAIN = "_yeet._tcp"
)


type Peer struct {
	Name	string
	IP		string
	Port	int
}

func NewPeer() *Peer {
	return &Peer{
		Name: GetUserName(),
		IP: GetLocalIP(),
		Port: PORT,
	}
}

// SENDER: This announces "Hey, I'm here and ready to send files!"
func (p *Peer) ActAsSender() {
	fmt.Println(p.Name)
	server, err := zeroconf.Register(
		p.Name,      // Your computer's name
		DOMAIN,  // Service type (like a category)
		"local.",           // Local network
		p.Port,               // Port number
		[]string{},         // Extra info (empty for now)
		nil,                // Use all network interfaces
	)
	
	if err != nil {
		log.Fatal(err)
	}
	defer server.Shutdown()

	fmt.Println("I'm now visible on the network!")
	fmt.Println("Others can find me at port 8080")
	fmt.Println("Press Ctrl+C to stop")
	
	// Keep running forever
	select {}
}

// RECEIVER: This searches for anyone broadcasting
func (p *Peer) ActAsReceiver() {
	fmt.Println("Looking for people sharing files...")
	
	// Channel to receive discovered services
	results := make(chan *zeroconf.ServiceEntry)
	
	// Create a context with timeout (5 seconds to search)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Start searching for "_fileshare._tcp" services
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
	}
	
	err = resolver.Browse(ctx, DOMAIN, "local.", results)
	if err != nil {
		log.Fatal(err)
	}
	
	// Collect results until timeout
	for {
		select {
		case entry := <-results:
			// Found someone!
			fmt.Println("\nFound:", entry.Instance)
			if len(entry.AddrIPv4) > 0 {
				fmt.Println("   IP:", entry.AddrIPv4[0])
			}
			fmt.Println("   Port:", entry.Port)
			// Now you can connect to this IP:Port to transfer files
			
		case <-ctx.Done():
			fmt.Println("\nSearch finished")
			return
		}
	}
}

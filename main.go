package main

import (
	"fmt"
	"os"

	"github.com/mohaanymo/yeet/network"
)



func main() {

	// if len(os.Args) < 2 {
	// 	network.ReceiverMode()
	// }

	mode := os.Args[1]
	
	switch mode {
	
	case "receive":
		network.ReceiverMode()
	case "send":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please specify a filepath")
			fmt.Println("Example: yeet send file.zip")
			os.Exit(1)
		}
		filename := os.Args[2]
		network.SenderMode(filename)

	default:
		fmt.Printf("Unknown command: %s\n", mode)
		os.Exit(1)
	}
}
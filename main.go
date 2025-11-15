package main

import (
	"fmt"
	"os"

	"yeet/network"
)



func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  To receive files:  yeet receive")
		fmt.Println("  To send files:     yeet send <filename>")
		fmt.Println("\nExample:")
		fmt.Println("  Receiver: yeet receive")
		fmt.Println("  Sender:   yeet send file.zip")
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "receive":
		network.ReceiverMode()

	case "send":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please specify a filename")
			fmt.Println("Example: yeet send file.zip")
			os.Exit(1)
		}
		filename := os.Args[2]
		network.SenderMode(filename)

	default:
		fmt.Printf("Unknown command: %s\n", mode)
		fmt.Println("Use 'receive' or 'send'")
		os.Exit(1)
	}
}
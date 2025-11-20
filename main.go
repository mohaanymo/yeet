package main

import (
	"os"

	"github.com/mohaanymo/yeet/network"
)


func main() {
	// simpler way to start
	// ./yeet file.zip E:\\MyFolder ... -> then run the sender mode
	// ./yeet -> then run the receiver mode
	if len(os.Args) >= 2 {
		filesAndDirs := os.Args[1:]
		network.SenderMode(filesAndDirs)
	}else if len(os.Args) < 2 {
		network.ReceiverMode()
	}

}
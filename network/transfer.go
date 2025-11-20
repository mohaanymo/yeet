package network

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohaanymo/yeet/protocol"
)

// AcceptConnection accepts only one connection then retunr it
func AcceptConnection(listener net.Listener) *net.Conn {

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error in accepting connection")
			continue
		}
		return &conn
	}
}

// SendFiles takes an array of files and start sending them one by one
func SendFiles(files []string, conn net.Conn) error {

	for _, filePath := range files {

		// Fet the filename and create a new file struct
		filename := filepath.Base(filePath)
		file, err := protocol.NewFile(filename, filePath)
		if err != nil {
			log.Fatal(err)
			return err
		}

		// Sending metadata message
		metadataMsg := protocol.NewFileMetadataMessage(file)
		_, err = conn.Write(metadataMsg.Encode())
		if err != nil {
			return fmt.Errorf("error sending metadata msg: %v", err)
		}

		// Printing the sending filename
		fmt.Println(file.Metadata.Filename)

		// Initialize some values to start the transfer process
		buf := make([]byte, CHUNKSIZE)
		totalSent := int64(0)
		pb := NewProgressBar(file.Metadata.FileSize)

		for {
			// reading bytes from file into the buffer
			n, err := file.Reader.Read(buf)
			if err == io.EOF || n == 0 {
				break
			}
			if err != nil {
				err = fmt.Errorf("failed reading into buffer: %v", err)
				errMsg := protocol.NewErrorMessage(err.Error())
				conn.Write(errMsg.Encode())
				return err
			}

			chunkMsg := protocol.NewFileDataMessage(buf[:n])
			_, err = conn.Write(chunkMsg.Encode())
			if err != nil {
				fmt.Printf("error writing data: %v\n", err)
			}

			// Updating progress bar
			totalSent += int64(n)
			pb.Update(totalSent)
		}

		// Finish sending
		pb.Finish()
		completeMsg := protocol.NewFileCompleteMessage()
		_, err = conn.Write(completeMsg.Encode())
		if err != nil {
			return fmt.Errorf("failed to send completion message: %w", err)
		}
		file.Reader.Close()

	}
	
	return nil
}

func ReceiveFiles(conn net.Conn) error {
	defer conn.Close()

	// Initialize some values preparing for receving files
	var metadata *protocol.Metadata
	var writer *bufio.Writer
	var file *os.File
	var receivedBytes int64 = 0
	var pb *ProgressBar


	// just for safety
	defer func() {
		if writer != nil {
			writer.Flush()
		}
		if file != nil {
			file.Close()
		}
	}()

	for {

		msg, err := protocol.Decode(conn)

		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("\n[*] %s disconnected\n", conn.RemoteAddr())
				return nil
			}
			if strings.Contains(err.Error(), "failed to read header") {
				fmt.Printf("\n[*] %s closed the connection cleanly\n", conn.RemoteAddr())
				return nil
			}
			err = fmt.Errorf("\n[-] Error reading message from %s: %v", conn.RemoteAddr(), err)
			fmt.Println(err.Error())
			return err
		}

		switch msg.Type {

		case protocol.MsgTypeError:
			fmt.Printf("message type error: %s\n", string(msg.Payload))

		case protocol.MsgTypeFileMetadata:
			metadata, err = protocol.ParseFileMetadata(msg.Payload)
			if err != nil {
				return fmt.Errorf("parsing metadata error: %v", err)
			}
			file, err = os.Create(metadata.Filename)
			if err != nil {
				return fmt.Errorf("failed to create output file: %v", err)
			}
			writer = bufio.NewWriterSize(file, CHUNKSIZE)
			pb = NewProgressBar(metadata.FileSize)
			
			// Printing sending filename
			fmt.Println(metadata.Filename)

		case protocol.MsgTypeFileData:
			if writer == nil {
				return fmt.Errorf("received data before metadata")
			}

			n, err := writer.Write(msg.Payload)
			if err != nil {
				return fmt.Errorf("error writing to file: %v", err)
			}

			// Updating progress bar
			receivedBytes += int64(n)
			pb.Update(receivedBytes)

		case protocol.MsgTypeFileComplete:
			if writer != nil {
				writer.Flush()
			}
			if file != nil {
				file.Close()
			}

			pb.Finish()
			
			// zeroing vars to send another file
			metadata, pb = nil, nil
			receivedBytes = 0

		default:
			return fmt.Errorf("unexpected message type: %d", msg.Type)
		}

	}

}

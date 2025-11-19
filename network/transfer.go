package network

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"os"

	"github.com/mohaanymo/yeet/protocol"
)

func AcceptConnections(listener net.Listener) {
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error in accepting connection")
			continue
		}
		defer conn.Close()
		ReceiveFile(conn)
		break
	}
}

func ConnectTo(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("connected to %s\n", addr)

	return conn, nil
}

func SendFile(f *protocol.File, conn net.Conn) error {
	defer f.Reader.Close()

	// Sending metadata message
	metadataMsg := protocol.NewFileMetadataMessage(f)
	_, err := conn.Write(metadataMsg.Encode())
	if err != nil {
		return fmt.Errorf("error sending metadata msg: %v", err)
	}

	buf := make([]byte, CHUNKSIZE)
	totalSent := int64(0)
	pb := NewProgressBar(f.Metadata.FileSize)

	for {
		n, err := f.Reader.Read(buf)

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
	return nil
}

func ReceiveFile(conn net.Conn) error {
	defer conn.Close()

	var metadata *protocol.Metadata
	var writer *bufio.Writer
	var file *os.File
	var receivedBytes int64 = 0
	var pb *ProgressBar

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
                fmt.Printf("[*] %s disconnected\n", conn.RemoteAddr())
                return nil
            }
            if strings.Contains(err.Error(), "failed to read header") {
                fmt.Printf("[*] %s closed the connection cleanly\n", conn.RemoteAddr())
                return nil
            }
			err = fmt.Errorf("[-] Error reading message from %s: %v", conn.RemoteAddr(), err)
            fmt.Println(err.Error())
            return err
        }

		switch msg.Type {

		case protocol.MsgTypeError:
			fmt.Printf("message type error: %s", string(msg.Payload))

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

		case protocol.MsgTypeFileData:
			if writer == nil {
				return fmt.Errorf("received data before metadata")
			}

			n, err := writer.Write(msg.Payload)
			if err != nil {
				return fmt.Errorf("error writing to file: %v", err)
			}

			receivedBytes += int64(n)
			pb.Update(receivedBytes)

		case protocol.MsgTypeFileComplete:
			if writer != nil {
				writer.Flush()
			}
			pb.Finish()
			return nil

		default:
			return fmt.Errorf("unexpected message type: %d", msg.Type)

		}

	}

}

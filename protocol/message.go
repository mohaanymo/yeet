package protocol

import (
	"encoding/binary"
	"io"
	"fmt"
)

const (
	MsgTypeFileMetadata = 0x01
	MsgTypeFileData     = 0x02
	MsgTypeFileComplete = 0x03
	MsgTypeError        = 0x04
)

type Message struct {
	Type    byte
	Payload []byte
}

func NewMessage(msgType byte, payload []byte) *Message {
	return &Message{
		Type:    msgType,
		Payload: payload,
	}
}

// Format: [Type:1Byte,PaylodeLength:4Bytes,Payload:N Bytes]
func (m *Message) Encode() []byte {
	payloadLen := len(m.Payload)
	buf := make([]byte, payloadLen+5)
	buf[0] = m.Type
	binary.BigEndian.PutUint32(buf[1:5], uint32(payloadLen))
	copy(buf[5:], m.Payload)
	
	return buf
}

func Decode(r io.Reader) (*Message, error) {
	header := make([]byte, 5)
	_, err := io.ReadFull(r, header)
	if err!=nil {
		return nil, err
	}
	msgType := header[0]
	payloadLen := binary.BigEndian.Uint32(header[1:5])

	payload := make([]byte, payloadLen)
	
	_, err = io.ReadFull(r, payload)
	if err!=nil {
		return nil, err
	}

	msg := NewMessage(byte(msgType), payload)
	return msg, nil
}

// NewFileMetadataMessage creates a file metadata message
func NewFileMetadataMessage(f *File) *Message {
    // Format: [FilenameLength:4bytes][Filename][FileSize:8bytes]
    filenameBytes := []byte(f.Metadata.Filename)
    filenameLen := len(filenameBytes)
    
    payload := make([]byte, 4+filenameLen+8)
    
    // Filename length
    binary.BigEndian.PutUint32(payload[0:4], uint32(filenameLen))
    
    // Filename
    copy(payload[4:4+filenameLen], filenameBytes)
    
    // File size
    binary.BigEndian.PutUint64(payload[4+filenameLen:], uint64(f.Metadata.FileSize))
    
    return NewMessage(MsgTypeFileMetadata, payload)
}

// ParseFileMetadata extracts filename and size from metadata message
func ParseFileMetadata(payload []byte) (*Metadata, error) {
    if len(payload) < 12 { // At least 4 (len) + 8 (size)
        return nil, fmt.Errorf("invalid metadata payload")
    }
    
    // Read filename length
    filenameLen := binary.BigEndian.Uint32(payload[0:4])
    
    if len(payload) < int(4+filenameLen+8) {
        return nil, fmt.Errorf("payload too short")
    }
    
    // Read filename
    filename := string(payload[4 : 4+filenameLen])
    
    // Read file size
    fileSize := int64(binary.BigEndian.Uint64(payload[4+filenameLen:]))
    
    return &Metadata{
		Filename: filename,
		FilenameLen: len(filename),
		FileSize: int64(fileSize),
	}, nil
}

// NewFileDataMessage creates a file data message
func NewFileDataMessage(data []byte) *Message {
    return NewMessage(MsgTypeFileData, data)
}

// NewFileCompleteMessage creates a file complete message
func NewFileCompleteMessage() *Message {
    return NewMessage(MsgTypeFileComplete, nil)
}

// NewErrorMessage creates an error message
func NewErrorMessage(errMsg string) *Message {
    return NewMessage(MsgTypeError, []byte(errMsg))
}
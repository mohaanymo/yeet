![YEET Banner](./assets/yeet-banner.svg)

A fast, simple, and decentralized file transfer tool that uses mDNS (Zeroconf) for automatic device discovery on your local network. No configuration needed—just send and receive files seamlessly between devices.

## Features

- **🔍 Automatic Discovery** - Devices are automatically discovered on your local network using mDNS
- **🛡️ Peer-to-Peer** - Direct TCP connection between sender and receiver, no intermediary servers
- **⚡ Fast** - Simple binary protocol optimized for speed
- **🎯 Interactive Selection** - Choose which device to send files to from a list of available receivers
- **💻 Cross-Platform** - Works on Windows, macOS, and Linux
- **🔧 Zero Configuration** - Just run and go, no setup required

## Installation

### Prerequisites
- Go 1.23.5 or higher

### Build from Source

```bash
git clone https://github.com/mohaanymo/yeet.git
cd yeet
go mod download
go build -o yeet
```

## Usage

### Receiver Mode
Start listening for incoming files on your device:

```bash
./yeet
```

This will:
- Announce your device on the local network
- Listen for incoming file transfers on port 9090
- Wait for a sender to connect and transfer a file
- Automatically save received files

### Sender Mode
Send a file to an available receiver:

```bash
./yeet <filepath> <dir> ...
```

This will:
1. Collect all files you want to send
2. Search for devices running in receiver mode (30 second timeout)
3. Display a numbered list of available receivers
4. Let you choose which device to send to (or wait for more devices)
5. Connect and transfer the file

**Example:**
```bash
./yeet myfile.zip
./yeet "/path/to/my folder"
./yeet /home/user/photos/vacation.tar.gz
```

## How It Works

1. **Receiver Advertises**: When you run `yeet`, the device registers itself on the mDNS network
2. **Sender Discovers**: When you run `yeet "somefiles"`, the tool searches for all devices advertising the yeet service
3. **User Selection**: You interactively choose which receiver to send files to
4. **Direct Transfer**: A TCP connection is established and the files are sent directly between devices

### Network Protocol

Yeet uses a custom binary protocol for efficient file transfer:
- Protocol version negotiation
- File metadata transmission (filename, size, etc.)
- Chunked file data transfer
- Integrity verification

## Configuration

By default, Yeet uses:
- **Port**: 9090 (TCP)
- **Service Name**: `_yeet._tcp`
- **Network**: Local network only (mDNS)
- **Discovery Timeout**: 30 seconds for sender mode

## Architecture

```
yeet/
├── main.go                 # Entry point and CLI handler
├── network/
│   ├── broadcast.go        # Sender/Receiver modes and discovery
│   ├── config.go           # Network configuration constants
│   ├── transfer.go         # File transfer logic
│   └── progress.go         # Transfer progress tracking
└── protocol/
    ├── file.go             # File metadata handling
    ├── message.go          # Protocol message definitions
    └── utils.go            # Protocol utilities
```

## Development

### Requirements
- Go 1.23.5+
- `github.com/grandcat/zeroconf` - mDNS/Zeroconf implementation

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests to improve Yeet.

## Project Status

Yeet is actively under development. Current features are functional, with ongoing improvements to transfer speed and reliability.

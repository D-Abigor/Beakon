# Beakon

A lightweight UDP broadcast tool written in Go for discovering a home lab server's IP address on a network without static IP assignment.

## Problem

Running a home lab on a college campus network where the server receives a different DHCP-assigned IP on each boot. Without a static IP or dynamic DNS, there's no reliable way to SSH into the lab from another device on the same network.

## Solution

Beakon solves this with a two-component broadcast/listen pattern:

1. **Broadcaster** (`Beakon.go`) - Runs on the lab server. Enumerates local network interfaces, computes the subnet broadcast address by bitwise-ORing the host IP against the inverted subnet mask, and sends the server's IP via UDP broadcast.

2. **Listener** (`UDP String Listener.go`) - Runs on any device on the same network. Binds to a configurable UDP port (default `52121`), receives the broadcast, and prints the server's IP address so you can connect.

## How the Broadcast Address Is Calculated

```go
for i := 0; i < len(ip); i++ {
    broadcast[i] = ip[i] + ^ipmask[i]
}
```

The broadcast IP is derived by taking each byte of the host IP and ORing it with the bitwise complement of the subnet mask, standard network math, implemented manually rather than relying on a library.

## Listener Features

- Configurable listen port and initial buffer size via interactive prompts
- Automatic buffer growth: if a received datagram fills the buffer, the buffer doubles (up to the 65535-byte UDP maximum) and logs a warning
- Color-coded terminal output for received transmissions and errors

## Usage

**On the lab server (broadcaster):**
```bash
go run Beakon.go
```

- although ideally u would have project compiled and run as a daemon persistently.

**On your laptop (listener):**
```bash
go run "UDP String Listener.go"
```

The listener will prompt for a port (default `52121`) and buffer size (default `1024`), then wait for incoming broadcasts.

## Project Structure

```
Beakon/
├── Beakon.go                # Broadcaster - computes and sends UDP broadcast
├── UDP String Listener.go   # Listener - receives and displays broadcasts
└── go.mod
```

## Tech Stack

- **Language:** Go
- **Networking:** UDP broadcast via `net.ListenUDP` / `net.UDPAddr`
- **Platform:** Built for Linux (home lab) and macOS (client)

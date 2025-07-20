
# TSL v5 Go Module

A reliable, production-grade Go module for the [TSL UMD v5 protocol](https://www.tslproducts.com/umd/). Supports both receiver and sender operation, with robust logging and flexible deployment.

---

## Features

- TSL v5 packet encoding/decoding
- TCP and UDP transport support (receiver and sender)
- Modular logging with file rotation and color-coded console output
- Idiomatic Go API and extensible types

---

## Quick Start

### Install

```
go get github.com/sevanjam/tsl/v5
```

### Usage Example: Receiver

```go
package main

import (
    "github.com/sevanjam/tsl/v5"
)

func main() {
    v5.InitLogger(true, true, "logs")
    v5.StartReceiver("udp", 5729, func(pkt *v5.Packet, addr string) {
        v5.LogInfof("Packet received from %s: %+v", addr, pkt)
    })
    select {}
}
```

### Usage Example: Sender

```go
package main

import (
    "github.com/sevanjam/tsl/v5"
    "time"
)

func main() {
    v5.InitLogger(true, false, "")
    sender, err := v5.NewSender("udp", "127.0.0.1", 5729)
    if err != nil {
        v5.LogErrorf("Failed to create sender: %v", err)
        return
    }
    defer sender.Close()

    pkt := &v5.Packet{ /* ...fill fields... */ }
    data, err := v5.MarshalPacket(pkt)
    if err != nil {
        v5.LogErrorf("Marshal failed: %v", err)
        return
    }

    for i := 0; i < 5; i++ {
        sender.Send(data)
        time.Sleep(time.Second)
    }
}
```

---

## API Reference

### Types

#### Packet

```go
type Packet struct {
    Version  byte
    Flags    byte
    Screen   uint16
    Messages []DisplayMessage
}
```
Represents a TSL v5 protocol packet containing all display messages.

---

#### DisplayMessage

```go
type DisplayMessage struct {
    Index       uint16
    ControlWord uint16
    Tally       byte
    Text        string
}
```
A single message for a UMD/tally display in TSL v5 format.

---

#### Sender

```go
type Sender struct {
    // Fields are internal; use constructor and methods.
}
```
A sender for transmitting TSL v5 packets over TCP or UDP.

---

### Functions

#### StartReceiver

```go
func StartReceiver(proto string, port int, handler PacketHandler) error
```
Starts a TSL v5 receiver on the specified protocol (`"udp"` or `"tcp"`) and port.  
Invokes the provided handler function for each received and decoded packet.

---

#### NewSender

```go
func NewSender(proto, ip string, port int) (*Sender, error)
```
Creates a new sender for the given protocol, IP address, and port.

---

#### (Sender) Send

```go
func (s *Sender) Send(data []byte) (int, error)
```
Sends a raw TSL v5 packet to the configured address.

---

#### (Sender) Close

```go
func (s *Sender) Close()
```
Closes the sender's network connection.

---

#### (DisplayMessage) MarshalBinary

```go
func (d DisplayMessage) MarshalBinary() ([]byte, error)
```
Serializes a display message into TSL v5 binary format.

---

#### (Packet) MarshalBinary

```go
func (p Packet) MarshalBinary() ([]byte, error)
```
Serializes the packet (and all messages) into TSL v5 binary format.

---

#### UnmarshalPacket

```go
func UnmarshalPacket(data []byte) (*Packet, error)
```
Parses a raw TSL v5 byte slice into a `Packet` structure, including all display messages.

---

### Logging

#### InitLogger

```go
func InitLogger(enableConsole bool, enableFile bool, dir string) error
```
Initializes logging to console, file, or both.  
Specify a directory for file logging.

---

#### LogInfof, LogWarnf, LogErrorf

```go
func LogInfof(format string, args ...interface{})
func LogWarnf(format string, args ...interface{})
func LogErrorf(format string, args ...interface{})
```
Log info, warning, or error messages.

---

#### CloseLogger

```go
func CloseLogger()
```
Flushes and closes the log file if open.

---

### Type Aliases

#### PacketHandler

```go
type PacketHandler func(pkt *Packet, addr string)
```
A handler function for incoming packets, used with `StartReceiver`.

---

### Example

```go
v5.StartReceiver("udp", 5729, func(pkt *v5.Packet, addr string) {
    v5.LogInfof("Received packet: %+v", pkt)
})
```

---

_For full documentation, see Go source comments or [pkg.go.dev](https://pkg.go.dev/github.com/sevanjam/tsl/v5)._

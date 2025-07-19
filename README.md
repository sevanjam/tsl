# TSL UMD Protocol v5 — Go Package

This Go module implements the [TSL UMD Protocol Version 5](https://tslproducts.com/media/5426/tsl-umd-protocol-specification-v5.pdf), which is widely used in professional broadcast environments for controlling under-monitor displays (UMDs), tally lights, and source labels.

## Features

- TCP and UDP transport (inbound and outbound)
- Full support for TSL v5 framing and byte-stuffing
- Structured decoding and encoding of packets and display messages
- Tally and brightness control
- ASCII text display support (UTF-16 planned for a future release)
- Modular, testable, and reusable code structure

API Reference

Types
Packet

type Packet struct {
    Version  byte
    Flags    byte
    Screen   uint16
    Messages []DisplayMessage
}
Represents a TSL v5 protocol packet containing all display messages.

DisplayMessage

type DisplayMessage struct {
    Index       uint16
    ControlWord uint16
    Tally       byte
    Text        string
}
A single message for a UMD/tally display in TSL v5 format.

Sender

type Sender struct {
    // Fields are internal; use constructor and methods.
}
A sender for transmitting TSL v5 packets over TCP or UDP.

Functions
StartReceiver

func StartReceiver(proto string, port int, handler PacketHandler) error
Starts a TSL v5 receiver on the specified protocol ("udp" or "tcp") and port.
Invokes the provided handler function for each received and decoded packet.

NewSender

func NewSender(proto, ip string, port int) (*Sender, error)
Creates a new sender for the given protocol, IP address, and port.

(*Sender) Send

func (s *Sender) Send(data []byte) (int, error)
Sends a raw TSL v5 packet to the configured address.

(*Sender) Close

func (s *Sender) Close()
Closes the sender's network connection.

(DisplayMessage) MarshalBinary

func (d DisplayMessage) MarshalBinary() ([]byte, error)
Serializes a display message into TSL v5 binary format.

(Packet) MarshalBinary

func (p Packet) MarshalBinary() ([]byte, error)
Serializes the packet (and all messages) into TSL v5 binary format.

UnmarshalPacket

func UnmarshalPacket(data []byte) (*Packet, error)
Parses a raw TSL v5 byte slice into a Packet structure, including all display messages.

Logging
InitLogger

func InitLogger(enableConsole bool, enableFile bool, dir string) error
Initializes logging to console, file, or both.
Specify a directory for file logging.

LogInfof, LogWarnf, LogErrorf

func LogInfof(format string, args ...interface{})
func LogWarnf(format string, args ...interface{})
func LogErrorf(format string, args ...interface{})
Log info, warning, or error messages.

CloseLogger

func CloseLogger()
Flushes and closes the log file if open.

Type Aliases
PacketHandler

type PacketHandler func(pkt *Packet, addr string)
A handler function for incoming packets, used with StartReceiver.

Example
v5.StartReceiver("udp", 5729, func(pkt *v5.Packet, addr string) {
    v5.LogInfof("Received packet: %+v", pkt)
})
For full documentation, see Go source comments or pkg.go.dev.


## Getting Started

Example: Sending a TSL Display Message over TCP

package main

import (
	"github.com/yourname/tslv5/v5"
)

func main() {
	sender, err := v5.NewSender("tcp", "192.168.100.65", 5727)
	if err != nil {
		panic(err)
	}
	defer sender.Close()

	// LH = GREEN, TEXT = RED, RH = OFF, Brightness = FULL
	msg := v5.NewDMSG(3, "CAM 3", 2, 1, 0, 3)

	packet := v5.Packet{
		Version:  0,
		Flags:    0x00, // ASCII mode
		Screen:   0,
		Messages: []v5.DMSG{msg},
	}

	data, err := packet.MarshalBinary()
	if err != nil {
		panic(err)
	}

	sender.Send(data)
}


## Testing

You can use main.go as a standalone test application. It allows you to:

Listen for inbound TSL packets over TCP or UDP
Send outbound packets for testing UMD receivers
Notes

TCP framing with DLE/STX and byte-stuffing is handled automatically
ASCII text is currently supported; UTF-16LE encoding will be added in a future version
Display control fields such as tally state and brightness are configurable using helper methods
License




## Author

Created by Sevan Malatjalian
Broadcast Systems Engineer

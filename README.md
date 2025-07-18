# TSL UMD Protocol v5 — Go Package

This Go module implements the [TSL UMD Protocol Version 5](https://tslproducts.com/media/5426/tsl-umd-protocol-specification-v5.pdf), which is widely used in professional broadcast environments for controlling under-monitor displays (UMDs), tally lights, and source labels.

## Features

- TCP and UDP transport (inbound and outbound)
- Full support for TSL v5 framing and byte-stuffing
- Structured decoding and encoding of packets and display messages
- Tally and brightness control
- ASCII text display support (UTF-16 planned for a future release)
- Modular, testable, and reusable code structure

## Package Structure

tslv5/
├── go.mod
├── main.go         # Optional test
└── v5/
├── receiver.go     # Handles TCP/UDP inbound transport
├── sender.go       # Handles outbound transport
├── types.go        # Core protocol types (Packet, DMSG, etc.)
├── encoder.go      # Encoding logic for packets and messages
├── decoder.go      # Decoding and parsing logic
└── doc.go          # Package documentation


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
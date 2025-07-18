package main

import (
	"encoding/binary"
	"fmt"
	"time"
	v5 "tslv5/v5" // use your correct module path
)

func main() {
	sender, err := v5.NewSender("tcp", "10.211.55.3", 5729) // adjust IP/port
	if err != nil {
		panic(err)
	}
	defer sender.Close()

	for {
		msg := v5.NewDMSG(3, "CAM 55", 2, 1, 0, 3)
		pkt := v5.Packet{
			Version:  0,
			Flags:    0x00, // ASCII
			Screen:   0,
			Messages: []v5.DMSG{msg},
		}

		data, err := pkt.MarshalBinary()
		if err != nil {
			fmt.Println("Marshal error:", err)
			return
		}

		err = sender.Send(data)

		time.Sleep(2 * time.Second)
	}
}

func buildExampleTSLPacket() ([]byte, error) {
	text := "CAM 3"
	textBytes := []byte(text)
	textLen := uint16(len(textBytes))

	// Control Word:
	// RH Tally = 0 (OFF)
	// Text Tally = 1 (RED)
	// LH Tally = 2 (GREEN)
	// Brightness = 3 (FULL)
	// Control bit 15 = 0 (display data)
	controlWord := uint16(0)
	controlWord |= 0b00      // RH Tally
	controlWord |= 0b01 << 2 // Text Tally
	controlWord |= 0b10 << 4 // LH Tally
	controlWord |= 0b11 << 6 // Brightness (FULL)

	dmsg := make([]byte, 0)
	buf := make([]byte, 2)

	// Index (2 bytes)
	binary.LittleEndian.PutUint16(buf, 3)
	dmsg = append(dmsg, buf...)

	// Control Word (2 bytes)
	binary.LittleEndian.PutUint16(buf, controlWord)
	dmsg = append(dmsg, buf...)

	// Length (2 bytes)
	binary.LittleEndian.PutUint16(buf, textLen)
	dmsg = append(dmsg, buf...)

	// Text
	dmsg = append(dmsg, textBytes...)

	// Top-Level Packet
	header := make([]byte, 6)
	binary.LittleEndian.PutUint16(header[0:2], uint16(len(dmsg)+4)) // PBC
	header[2] = 0x00                                                // Version
	header[3] = 0x00                                                // Flags (ASCII)
	binary.LittleEndian.PutUint16(header[4:6], 0x0000)              // Screen 0

	return append(header, dmsg...), nil
}

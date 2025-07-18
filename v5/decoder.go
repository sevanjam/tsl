package v5

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// UnmarshalPacket parses a raw TSL V5 payload (after framing/unstuffing)
func UnmarshalPacket(data []byte) (*Packet, error) {
	if len(data) < 6 {
		return nil, errors.New("packet too short")
	}

	p := &Packet{}
	pbc := binary.LittleEndian.Uint16(data[0:2])
	if pbc != uint16(len(data)-2) {
		return nil, fmt.Errorf("byte count mismatch: header says %d, got %d", pbc, len(data)-2)
	}

	p.Version = data[2]
	p.Flags = data[3]
	p.Screen = binary.LittleEndian.Uint16(data[4:6])

	// Start parsing messages
	cursor := 6
	for cursor+4 <= len(data) {
		msg := DMSG{}
		msg.Index = binary.LittleEndian.Uint16(data[cursor : cursor+2])
		msg.ControlWord = binary.LittleEndian.Uint16(data[cursor+2 : cursor+4])
		cursor += 4

		if msg.ControlWord&0x8000 != 0 {
			// Control data – not implemented yet
			return nil, errors.New("control data parsing not supported yet")
		}

		if cursor+2 > len(data) {
			return nil, errors.New("unexpected end of data while reading text length")
		}
		textLen := int(binary.LittleEndian.Uint16(data[cursor : cursor+2]))
		cursor += 2

		if cursor+textLen > len(data) {
			return nil, errors.New("unexpected end of data while reading text")
		}

		msg.Text = string(data[cursor : cursor+textLen])
		cursor += textLen

		p.Messages = append(p.Messages, msg)
	}

	return p, nil
}

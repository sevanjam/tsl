package v5

import (
	"encoding/binary"
	"errors"
)

func (d DMSG) MarshalBinary() ([]byte, error) {
	textBytes := []byte(d.Text)

	if len(textBytes) > 2048 {
		return nil, errors.New("text too long")
	}

	out := make([]byte, 6+len(textBytes)) // 2 + 2 + 2 + text
	binary.LittleEndian.PutUint16(out[0:2], d.Index)
	binary.LittleEndian.PutUint16(out[2:4], d.ControlWord)
	binary.LittleEndian.PutUint16(out[4:6], uint16(len(textBytes)))
	copy(out[6:], textBytes)

	return out, nil
}

func (p Packet) MarshalBinary() ([]byte, error) {
	payload := make([]byte, 6) // PBC, Version, Flags, Screen
	payload[2] = p.Version
	payload[3] = p.Flags
	binary.LittleEndian.PutUint16(payload[4:6], p.Screen)

	// Marshal all messages
	for _, msg := range p.Messages {
		encoded, err := msg.MarshalBinary()
		if err != nil {
			return nil, err
		}
		payload = append(payload, encoded...)
	}

	// Now calculate and insert PBC = total payload after PBC field (i.e. len - 2)
	pbc := uint16(len(payload) - 2)
	binary.LittleEndian.PutUint16(payload[0:2], pbc)

	return payload, nil
}

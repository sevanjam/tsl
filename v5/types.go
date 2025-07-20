package v5

// Packet represents a TSL v5 protocol packet containing version, flags, screen, and a list of display messages.
type Packet struct {
	Version  byte             // TSL protocol version
	Flags    byte             // Flags for the packet
	Screen   uint16           // Screen or address identifier
	Messages []DisplayMessage // List of display messages included in the packet
}

// DisplayMessage represents a single message in a TSL v5 packet, including index, control word, tally, and display text.
type DisplayMessage struct {
	Index       uint16 // Index of the message
	ControlWord uint16 // Control word (bitmask for tallies, attributes, etc.)
	Text        string // Display text for the UMD/tally
}

func (d DisplayMessage) IsControlData() bool {
	return d.ControlWord&(1<<15) != 0
}

func (d DisplayMessage) Brightness() uint8 {
	return uint8((d.ControlWord >> 6) & 0x03)
}

func (d DisplayMessage) LHTally() uint8 {
	return uint8((d.ControlWord >> 4) & 0x03)
}

func (d DisplayMessage) TextTally() uint8 {
	return uint8((d.ControlWord >> 2) & 0x03)
}

func (d DisplayMessage) RHTally() uint8 {
	return uint8(d.ControlWord & 0x03)
}

var TallyColor = map[uint8]string{
	0: "OFF",
	1: "RED",
	2: "GREEN",
	3: "AMBER",
}

func NewDMSG(index uint16, text string, lhTally, txtTally, rhTally, brightness uint8) DisplayMessage {
	control := uint16(0)
	control |= uint16(rhTally&0x03) << 0
	control |= uint16(txtTally&0x03) << 2
	control |= uint16(lhTally&0x03) << 4
	control |= uint16(brightness&0x03) << 6
	// Bit 15 = 0 (display data)

	return DisplayMessage{
		Index:       index,
		ControlWord: control,
		Text:        text,
	}
}

package v5

type Packet struct {
	Version  byte
	Flags    byte
	Screen   uint16
	Messages []DMSG
}

type DMSG struct {
	Index       uint16
	ControlWord uint16
	Text        string // or Unicode later
	// parsed fields coming later: LH Tally, RH Tally, etc.
}

func (d DMSG) IsControlData() bool {
	return d.ControlWord&(1<<15) != 0
}

func (d DMSG) Brightness() uint8 {
	return uint8((d.ControlWord >> 6) & 0x03)
}

func (d DMSG) LHTally() uint8 {
	return uint8((d.ControlWord >> 4) & 0x03)
}

func (d DMSG) TextTally() uint8 {
	return uint8((d.ControlWord >> 2) & 0x03)
}

func (d DMSG) RHTally() uint8 {
	return uint8(d.ControlWord & 0x03)
}

var TallyColor = map[uint8]string{
	0: "OFF",
	1: "RED",
	2: "GREEN",
	3: "AMBER",
}

func NewDMSG(index uint16, text string, lhTally, txtTally, rhTally, brightness uint8) DMSG {
	control := uint16(0)
	control |= uint16(rhTally&0x03) << 0
	control |= uint16(txtTally&0x03) << 2
	control |= uint16(lhTally&0x03) << 4
	control |= uint16(brightness&0x03) << 6
	// Bit 15 = 0 (display data)

	return DMSG{
		Index:       index,
		ControlWord: control,
		Text:        text,
	}
}

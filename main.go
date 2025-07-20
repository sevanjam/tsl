package main

import (
	"net"
	"os"
	"time"
	v5 "tslv5/v5"
)

func main() {
	// Initialize logger: console enabled, file disabled, no log dir needed for console-only
	v5.InitLogger(false, false, "Logs")

	// Uncomment one of the following lines in main() to test TCP or UDP receive mode:
	//
	// tcpReceiverTest()
	// udpReceiverTest()

	// Uncomment one of the following lines in main() to test TCP or UDP sending:
	// tcpSenderTest()
	// udpSenderTest()
}

func tcpReceiverTest() {
	// Start the receiver on TCP port 5728
	err := v5.StartReceiver("tcp", 5728, func(srcAddr net.Addr, data []byte) {
		v5.LogInfof("Received from %s: % X", srcAddr.String(), data)
		pkt, err := v5.UnmarshalPacket(data)
		if err != nil {
			v5.LogErrorf("Decode error: %v", err)
			return
		}
		v5.LogInfof("Decoded TSL v5 Packet: Version=%d Flags=0x%02X Screen=%d", pkt.Version, pkt.Flags, pkt.Screen)
		for i, msg := range pkt.Messages {
			v5.LogInfof("  [%d] - Index:%d     - Ctrl:0x%04X     - Brightness:%02X    - LH Tally:%02X    - Text Tally:%02X    - RH Tally:%02X    - Text:%q", i, msg.Index, msg.ControlWord, msg.Brightness(), msg.LHTally(), msg.TextTally(), msg.RHTally(), msg.Text)
		}
	})
	if err != nil {
		v5.LogErrorf("Failed to start TSL receiver: %v", err)
		os.Exit(1)
	}
	select {} // Block forever
}
func udpReceiverTest() {
	// Start the receiver on TCP port 5728
	err := v5.StartReceiver("udp", 5728, func(srcAddr net.Addr, data []byte) {
		v5.LogInfof("Received from %s: % X", srcAddr.String(), data)
		pkt, err := v5.UnmarshalPacket(data)
		if err != nil {
			v5.LogErrorf("Decode error: %v", err)
			return
		}
		v5.LogInfof("Decoded TSL v5 Packet: Version=%d Flags=0x%02X Screen=%d", pkt.Version, pkt.Flags, pkt.Screen)
		for i, msg := range pkt.Messages {
			v5.LogInfof("  [%d] - Index:%d     - Ctrl:0x%04X     - Brightness:%02X    - LH Tally:%02X    - Text Tally:%02X    - RH Tally:%02X    - Text:%q", i, msg.Index, msg.ControlWord, msg.Brightness(), msg.LHTally(), msg.TextTally(), msg.RHTally(), msg.Text)
		}
	})
	if err != nil {
		v5.LogErrorf("Failed to start TSL receiver: %v", err)
		os.Exit(1)
	}
	select {} // Block forever
}
func tcpSenderTest() {
	sender, err := v5.NewSender("tcp", "10.211.55.3", 5729)
	if err != nil {
		v5.LogErrorf("Failed to create TCP sender: %v", err)
		os.Exit(1)
	}
	defer sender.Close()

	for {
		pkt := v5.Packet{
			Version: 5,
			Flags:   0,
			Screen:  1,
			Messages: []v5.DisplayMessage{
				{
					Index:       1,
					ControlWord: 0x0001,
					Text:        "Hello from TCP sender",
				},
			},
		}
		data, err := pkt.MarshalBinary()
		if err != nil {
			v5.LogErrorf("Marshal error: %v", err)
			continue
		}
		_, err = sender.Send(data)
		if err != nil {
			v5.LogErrorf("Send error: %v", err)
			continue
		}
		v5.LogInfof("Sent TSL packet over TCP: % X", data)
		time.Sleep(2 * time.Second)
	}
}

func udpSenderTest() {
	sender, err := v5.NewSender("udp", "10.211.55.3", 5729)
	if err != nil {
		v5.LogErrorf("Failed to create UDP sender: %v", err)
		os.Exit(1)
	}
	defer sender.Close()

	for {
		pkt := v5.Packet{
			Version: 5,
			Flags:   0,
			Screen:  1,
			Messages: []v5.DisplayMessage{
				{
					Index:       1,
					ControlWord: 0x0001,
					Text:        "Hello from UDP sender",
				},
			},
		}
		data, err := pkt.MarshalBinary()
		if err != nil {
			v5.LogErrorf("Marshal error: %v", err)
			continue
		}
		_, err = sender.Send(data)
		if err != nil {
			v5.LogErrorf("Send error: %v", err)
			continue
		}
		v5.LogInfof("Sent TSL packet over UDP: % X", data)
		time.Sleep(2 * time.Second)
	}
}

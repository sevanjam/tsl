package main

import (
	"net"
	v5 "tslv5/v5"
)

func handlePacket(src net.Addr, data []byte) {
	pkt, err := v5.UnmarshalPacket(data)
	if err != nil {
		v5.Logf("Failed to parse packet from %s: %v", src.String(), err)
		return
	}

	v5.Logf("Parsed packet from %s: %+v", src.String(), pkt)
}

func main() {
	err := v5.InitLogger(true, "tsl.log")
	if err != nil {
		panic(err)
	}
	defer v5.CloseLogger()

	v5.Logf("TSL v5 receiver starting on UDP port 5729")

	err = v5.StartReceiver("tcp", 5728, handlePacket)
	if err != nil {
		v5.Logf("Receiver error: %v", err)
		return
	}

	// Prevent exit — block forever
	select {}
}

package v5

import (
	"bytes"
	"fmt"
	"net"
)

const (
	DLE = 0xFE
	STX = 0x02
)

// PacketHandler is the callback you implement to process incoming TSL V5 packets.
// It provides the sender address and the raw unwrapped data.
type PacketHandler func(src net.Addr, data []byte)

// StartReceiver starts a TCP or UDP listener on the given port.
// For TCP, it handles DLE/STX framing and DLE byte unstuffing internally.
func StartReceiver(proto string, port int, handler PacketHandler) error {
	switch proto {
	case "udp":
		return startUDPReceiver(port, handler)
	case "tcp":
		return startTCPReceiver(port, handler)
	default:
		return fmt.Errorf("unsupported protocol: %s", proto)
	}
}

// --- UDP Receiver Implementation ---

func startUDPReceiver(port int, handler PacketHandler) error {
	addr := net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return fmt.Errorf("UDP listen failed: %w", err)
	}

	go func() {
		defer conn.Close()
		buf := make([]byte, 2048)
		for {
			n, remote, err := conn.ReadFrom(buf)
			if err != nil {
				fmt.Println("UDP read error:", err)
				continue
			}
			handler(remote, buf[:n])
		}
	}()

	fmt.Printf("TSL V5 UDP receiver listening on port %d\n", port)
	return nil
}

// --- TCP Receiver Implementation ---

func startTCPReceiver(port int, handler PacketHandler) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return fmt.Errorf("TCP listen failed: %w", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("TCP accept error:", err)
				continue
			}
			go handleTCPConnection(conn, handler)
		}
	}()

	fmt.Printf("TSL V5 TCP receiver listening on port %d\n", port)
	return nil
}

func handleTCPConnection(conn net.Conn, handler PacketHandler) {
	defer conn.Close()

	remote := conn.RemoteAddr()
	buffer := make([]byte, 0, 4096)
	tmp := make([]byte, 1024)

	fmt.Println("TCP client connected from", remote)

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			fmt.Println("TCP read error:", err)
			return
		}

		buffer = append(buffer, tmp[:n]...)

		for {
			packet, rest, ok := extractFramedPacket(buffer)
			if !ok {
				break
			}
			handler(remote, packet)
			buffer = rest
		}
	}
}

// --- TCP Framing: Extract one full framed and unstuffed packet from the buffer ---

func extractFramedPacket(data []byte) ([]byte, []byte, bool) {
	start := bytes.Index(data, []byte{DLE, STX})
	if start == -1 {
		return nil, data, false
	}

	// Look for the next start-of-frame after this one
	nextStart := bytes.Index(data[start+2:], []byte{DLE, STX})
	var end int
	if nextStart != -1 {
		end = start + 2 + nextStart
	} else {
		end = len(data)
	}

	// Extract framed payload (between start+2 and end)
	raw := data[start+2 : end]
	unstuffed := make([]byte, 0, len(raw))

	for i := 0; i < len(raw); i++ {
		if raw[i] == DLE {
			if i+1 < len(raw) && raw[i+1] == DLE {
				unstuffed = append(unstuffed, DLE)
				i++ // skip the second DLE
			} else {
				// DLE followed by non-DLE — treat as raw
				unstuffed = append(unstuffed, raw[i])
			}
		} else {
			unstuffed = append(unstuffed, raw[i])
		}
	}

	return unstuffed, data[end:], true
}

package v5

import (
	"fmt"
	"net"
)

// Sender represents a network sender for TSL v5 packets over TCP or UDP.
type Sender struct {
	protocol string   // Protocol ("tcp" or "udp")
	addr     string   // Destination address (ip:port)
	udp      net.Conn // UDP connection (if proto is "udp")
	tcp      net.Conn // TCP connection (if proto is "tcp")
}

// NewSender creates a new Sender for the given protocol ("tcp" or "udp"), IP address, and port.
// Returns a pointer to Sender or an error if the connection could not be established.
func NewSender(protocol, ip string, port int) (*Sender, error) {
	LogInfof("Initializing sender to %s:%d over %s", ip, port, protocol)

	addr := fmt.Sprintf("%s:%d", ip, port)
	s := &Sender{protocol: protocol, addr: addr}

	switch protocol {
	case "udp":
		conn, err := net.Dial("udp", addr)
		if err != nil {
			LogErrorf("Failed to connect UDP to %s: %v", addr, err)
			return nil, err
		}
		s.udp = conn
		LogInfof("UDP sender connected to %s", addr)
	case "tcp":
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			LogErrorf("Failed to connect TCP to %s: %v", addr, err)
			return nil, err
		}
		s.tcp = conn
		LogInfof("TCP sender connected to %s", addr)
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
	return s, nil
}

// Send writes a raw byte payload to the established connection (TCP or UDP).
// Returns the number of bytes sent and any error encountered.
func (s *Sender) Send(data []byte) (int, error) {
	if s.protocol == "udp" && s.udp != nil {
		n, err := s.udp.Write(data)
		LogInfof("Sent %d bytes via UDP to %s", n, s.addr)
		return n, err
	}
	if s.protocol == "tcp" && s.tcp != nil {
		data = wrapTSLV5Packet(data)
		n, err := s.tcp.Write(data)
		LogInfof("Sent %d bytes via TCP to %s", n, s.addr)
		return n, err
	}
	return 0, fmt.Errorf("connection not established")
}

// Close shuts down the sender's network connection.
func (s *Sender) Close() {
	if s.udp != nil {
		s.udp.Close()
	}
	if s.tcp != nil {
		s.tcp.Close()
	}
}
func wrapTSLV5Packet(payload []byte) []byte {
	const DLE = 0xFE
	const STX = 0x02

	// Byte-stuffing: any occurrence of DLE becomes DLE DLE
	stuffed := make([]byte, 0, len(payload)+8)
	for _, b := range payload {
		stuffed = append(stuffed, b)
		if b == DLE {
			stuffed = append(stuffed, DLE)
		}
	}
	// Wrap with DLE/STX at the start (NO ETX required per your protocol version)
	wrapped := []byte{DLE, STX}
	wrapped = append(wrapped, stuffed...)
	return wrapped
}

package v5

import (
	"bytes"
	"fmt"
	"net"
)

type Sender struct {
	proto string
	addr  string
	tcp   net.Conn
	udp   *net.UDPConn
}

// NewSender connects to a remote TSL V5 device
func NewSender(proto, ip string, port int) (*Sender, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	s := &Sender{
		proto: proto,
		addr:  addr,
	}

	switch proto {
	case "udp":
		remoteAddr, err := net.ResolveUDPAddr("udp4", addr)
		if err != nil {
			return nil, fmt.Errorf("resolve UDP addr: %w", err)
		}
		conn, err := net.DialUDP("udp", nil, remoteAddr)
		if err != nil {
			return nil, fmt.Errorf("dial UDP: %w", err)
		}
		s.udp = conn

	case "tcp":
		conn, err := net.Dial("tcp4", addr)
		if err != nil {
			return nil, fmt.Errorf("dial TCP: %w", err)
		}
		s.tcp = conn

	default:
		return nil, fmt.Errorf("unsupported protocol: %s", proto)
	}

	return s, nil
}

// Send transmits raw TSL V5 packet data
func (s *Sender) Send(data []byte) error {
	switch s.proto {
	case "udp":
		_, err := s.udp.Write(data)
		return err

	case "tcp":
		framed := frameTCP(data)
		_, err := s.tcp.Write(framed)
		return err

	default:
		return fmt.Errorf("unknown protocol: %s", s.proto)
	}
}

// Close shuts down the connection
func (s *Sender) Close() error {
	if s.tcp != nil {
		return s.tcp.Close()
	}
	if s.udp != nil {
		return s.udp.Close()
	}
	return nil
}

// frameTCP adds DLE/STX at the beginning and escapes all DLE bytes
func frameTCP(data []byte) []byte {
	var buf bytes.Buffer
	buf.WriteByte(DLE)
	buf.WriteByte(STX)

	for _, b := range data {
		buf.WriteByte(b)
		if b == DLE {
			buf.WriteByte(DLE) // Byte-stuffing for DLE
		}
	}

	return buf.Bytes()
}

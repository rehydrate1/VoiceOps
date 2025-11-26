package service

import (
	"fmt"
	"net"
)

func WakeOnLan(macAddr, broadcastAddr string) error {
	macBytes, err := net.ParseMAC(macAddr)
	if err != nil {
		return fmt.Errorf("invalid mac address: %w", err)
	}

	packet := []byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}

	for range 16 {
		packet = append(packet, macBytes...)
	}

	localAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	remoteAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		return fmt.Errorf("invalid broadcast address: %w", err)
	}

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		return fmt.Errorf("udp connection failed: %w", err)
	}
	defer conn.Close()

	n, err := conn.Write(packet)
	if err != nil {
		return fmt.Errorf("failed to write packet: %w", err)
	}

	if n != 102 {
		return fmt.Errorf("wrote %d bytes, expected 102", n)
	}

	return nil
}
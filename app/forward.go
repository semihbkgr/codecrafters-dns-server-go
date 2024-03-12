package main

import (
	"errors"
	"net"
)

type Forwarder struct {
	addr *net.UDPAddr
}

func NewForwarder(resolverAddr string) (*Forwarder, error) {
	addr, err := net.ResolveUDPAddr("udp", resolverAddr)
	if err != nil {
		return nil, err
	}

	return &Forwarder{
		addr: addr,
	}, nil
}

func (f *Forwarder) Forward(reqMsg *Message) (*Message, error) {
	conn, err := net.DialUDP("udp", nil, f.addr)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(reqMsg.Bytes())
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4096)
	size, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}

	response := ParseMessage(buf[:size])
	if response.Headers.ID != reqMsg.Headers.ID {
		return nil, errors.New("request and response message id mismatched")
	}

	return response, nil
}

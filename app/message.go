package main

import "encoding/binary"

/*
                                   1  1  1  1  1  1
     0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                      ID                       |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    QDCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    ANCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    NSCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    ARCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/

type MessageHeaders struct {
	//Packet Identifier
	ID uint16
	//Query/Response Indicator
	//Operation Code
	//Authoritative Answer
	//Truncation
	//Recursion Desired
	//Recursion Available
	//Reserved
	//Response Code
	CODE uint16
	//Question Count
	QDCOUNT uint16
	//Answer Record Count
	ANCOUNT uint16
	//Authority Record Count
	NSCOUNT uint16
	//Additional Record Count
	ARCOUNT uint16
}

func (h *MessageHeaders) QR() bool {
	return (h.CODE >> 15 & 0x01) == 1
}

func (h *MessageHeaders) OPCODE() uint8 {
	return uint8(h.CODE >> 11 & 0x0f)
}

func (h *MessageHeaders) AA() bool {
	return (h.CODE >> 10 & 0x01) == 1
}

func (h *MessageHeaders) TC() bool {
	return (h.CODE >> 9 & 0x01) == 1
}

func (h *MessageHeaders) RD() bool {
	return (h.CODE >> 8 & 0x01) == 1
}

func (h *MessageHeaders) RA() bool {
	return (h.CODE >> 7 & 0x01) == 1
}

func (h *MessageHeaders) Z() uint8 {
	return uint8(h.CODE >> 4 & 0x07)
}

func (h *MessageHeaders) RCODE() uint8 {
	return uint8(h.CODE >> 0 & 0x0F)
}

func ParseMessageHeaders(b []byte) MessageHeaders {
	return MessageHeaders{
		ID:      binary.BigEndian.Uint16(b[0:2]),
		CODE:    binary.BigEndian.Uint16(b[2:4]),
		QDCOUNT: binary.BigEndian.Uint16(b[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(b[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(b[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(b[10:12]),
	}
}

func (headers MessageHeaders) WriteMessageHeaders(b []byte) {
	binary.BigEndian.AppendUint16(b, headers.ID)
	binary.BigEndian.AppendUint16(b, headers.CODE)
	binary.BigEndian.AppendUint16(b, headers.QDCOUNT)
	binary.BigEndian.AppendUint16(b, headers.ANCOUNT)
	binary.BigEndian.AppendUint16(b, headers.NSCOUNT)
	binary.BigEndian.AppendUint16(b, headers.ARCOUNT)
}

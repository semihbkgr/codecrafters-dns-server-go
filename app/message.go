package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type Message struct {
	Headers   *MessageHeaders
	Questions []*MessageQuestion
	Answers   []*MessageResourceRecord
}

func ParseMessage(b []byte) *Message {
	headers := ParseMessageHeaders(b[0:12])

	offset := 12
	questions := make([]*MessageQuestion, 0, headers.QDCOUNT)
	for i := uint16(0); i < headers.QDCOUNT; i++ {
		question, questionOffset := ParseMessageQuestion(b[offset:])
		questions = append(questions, question)
		offset += questionOffset
	}

	return &Message{
		Headers:   headers,
		Questions: questions,
	}
}

func (m *Message) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(m.Headers.Bytes())
	for _, question := range m.Questions {
		buf.Write(question.Bytes())
	}
	for _, answer := range m.Answers {
		buf.Write(answer.Bytes())
	}
	return buf.Bytes()
}

func (m *Message) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%+v\n", *m.Headers))
	for _, question := range m.Questions {
		b.WriteString(fmt.Sprintf("%+v\n", *question))
	}
	for _, rr := range m.Answers {
		b.WriteString(fmt.Sprintf("%+v\n", *rr))
	}
	return b.String()
}

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

func (h *MessageHeaders) SetQR(b bool) {
	h.CODE |= 1 << 15
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

func ParseMessageHeaders(b []byte) *MessageHeaders {
	return &MessageHeaders{
		ID:      binary.BigEndian.Uint16(b[0:2]),
		CODE:    binary.BigEndian.Uint16(b[2:4]),
		QDCOUNT: binary.BigEndian.Uint16(b[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(b[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(b[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(b[10:12]),
	}
}

func (headers *MessageHeaders) Bytes() []byte {
	b := make([]byte, 12)
	binary.BigEndian.PutUint16(b[0:2], headers.ID)
	binary.BigEndian.PutUint16(b[2:4], headers.CODE)
	binary.BigEndian.PutUint16(b[4:6], headers.QDCOUNT)
	binary.BigEndian.PutUint16(b[6:8], headers.ANCOUNT)
	binary.BigEndian.PutUint16(b[8:10], headers.NSCOUNT)
	binary.BigEndian.PutUint16(b[10:12], headers.ARCOUNT)
	return b
}

type MessageQuestion struct {
	QNAME  Labels
	QTYPE  Type
	QCLASS Class
}

func ParseMessageQuestion(b []byte) (*MessageQuestion, int) {
	qname, offset := ParseLabels(b)
	qtype := Type(binary.BigEndian.Uint16(b[offset : offset+2]))
	offset += 2
	qclass := Class(binary.BigEndian.Uint16(b[offset : offset+2]))
	offset += 2
	return &MessageQuestion{
		QNAME:  qname,
		QTYPE:  qtype,
		QCLASS: qclass,
	}, offset
}

func (question *MessageQuestion) Bytes() []byte {
	b := question.QNAME.Bytes()
	b = binary.BigEndian.AppendUint16(b, uint16(question.QTYPE))
	b = binary.BigEndian.AppendUint16(b, uint16(question.QCLASS))
	return b
}

type Labels []string

func ParseLabels(b []byte) (Labels, int) {
	offset := 0
	name := make(Labels, 0)
	for l := b[offset]; l > 0; l = b[offset] {
		name = append(name, string(b[offset+1:offset+int(l)]))
		offset += int(l) + 1
	}
	return name, offset
}

func (labels Labels) Bytes() []byte {
	b := make([]byte, 0)
	for _, label := range labels {
		b = append(b, byte(len(label)))
		b = append(b, []byte(label)...)
	}
	return append(b, 0)
}

type Type uint16

const (
	A Type = 1 + iota
	NS
	MD
	MF
	CNAME
	SOA
	MB
	MG
	MR
	NULL
	WKS
	PTR
	HINFO
	MINFO
	MX
	TXT
)

type Class uint16

const (
	IN Class = 1 + iota
	CS
	CH
	HS
)

type MessageResourceRecord struct {
	NAME     Labels
	TYPE     Type
	CLASS    Class
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

func NewMessageResourceRecord(question *MessageQuestion) *MessageResourceRecord {
	return &MessageResourceRecord{
		NAME:  question.QNAME,
		TYPE:  question.QTYPE,
		CLASS: question.QCLASS,
	}
}

func (rr *MessageResourceRecord) SetData(b []byte) {
	rr.RDLENGTH = uint16(len(b))
	rr.RDATA = b
}

func (rr *MessageResourceRecord) Bytes() []byte {
	b := rr.NAME.Bytes()
	b = binary.BigEndian.AppendUint16(b, uint16(rr.TYPE))
	b = binary.BigEndian.AppendUint16(b, uint16(rr.CLASS))
	b = binary.BigEndian.AppendUint32(b, rr.TTL)
	b = binary.BigEndian.AppendUint16(b, rr.RDLENGTH)
	buf := bytes.NewBuffer(b)
	binary.Write(buf, binary.BigEndian, rr.RDATA)
	return buf.Bytes()
}

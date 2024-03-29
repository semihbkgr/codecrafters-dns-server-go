package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type Message struct {
	Headers   *Headers
	Questions []*Question
	Answers   []*ResourceRecord
}

func ParseMessage(b []byte) *Message {
	headers := ParseHeaders(b[0:12])

	offset := 12
	questions := make([]*Question, 0, headers.QDCOUNT)
	for i := uint16(0); i < headers.QDCOUNT; i++ {
		var question *Question
		question, offset = ParseQuestion(b, offset)
		questions = append(questions, question)
	}

	answers := make([]*ResourceRecord, 0, headers.ANCOUNT)
	for i := uint16(0); i < headers.ANCOUNT; i++ {
		var resourceRecord *ResourceRecord
		resourceRecord, offset = ParseResourceRecord(b, offset)
		answers = append(answers, resourceRecord)
	}

	return &Message{
		Headers:   headers,
		Questions: questions,
		Answers:   answers,
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

type Headers struct {
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

func NewHeaders(id uint16) *Headers {
	return &Headers{
		ID: id,
	}
}

func (h *Headers) QR() bool {
	return (h.CODE >> 15 & 0b1) == 1
}

func (h *Headers) SetQR(b bool) {
	h.CODE &^= 0b1 << 15
	if b {
		h.CODE |= 0b1 << 15
	}
}

func (h *Headers) OPCODE() uint8 {
	return uint8(h.CODE >> 11 & 0b1111)
}

func (h *Headers) SetOPCODE(opcode uint8) {
	h.CODE &^= uint16(0b1111) << 11
	h.CODE |= uint16(opcode) << 11
}

func (h *Headers) AA() bool {
	return (h.CODE >> 10 & 0b1) == 1
}

func (h *Headers) SetAA(b bool) {
	h.CODE &^= 0b1 << 10
	if b {
		h.CODE |= 0b1 << 10
	}
}

func (h *Headers) TC() bool {
	return (h.CODE >> 9 & 0b1) == 1
}

func (h *Headers) SetTC(b bool) {
	h.CODE &^= 0b1 << 9
	if b {
		h.CODE |= 0b1 << 9
	}
}

func (h *Headers) RD() bool {
	return (h.CODE >> 8 & 0b1) == 1
}

func (h *Headers) SetRD(b bool) {
	h.CODE &^= 0b1 << 8
	if b {
		h.CODE |= 0b1 << 8
	}
}

func (h *Headers) RA() bool {
	return (h.CODE >> 7 & 0b1) == 1
}

func (h *Headers) SetRA(b bool) {
	h.CODE &^= 0b1 << 7
	if b {
		h.CODE |= 0b1 << 7
	}
}

func (h *Headers) Z() uint8 {
	return uint8(h.CODE >> 4 & 0b111)
}

func (h *Headers) SetZ(z uint8) {
	h.CODE &^= uint16(0b111) << 4
	h.CODE |= uint16(z) << 4
}

func (h *Headers) RCODE() uint8 {
	return uint8(h.CODE >> 0 & 0b1111)
}

func (h *Headers) SetRCODE(rcode uint8) {
	h.CODE &^= uint16(0b1111) << 0
	h.CODE |= uint16(rcode) << 0
}

func ParseHeaders(b []byte) *Headers {
	return &Headers{
		ID:      binary.BigEndian.Uint16(b[0:2]),
		CODE:    binary.BigEndian.Uint16(b[2:4]),
		QDCOUNT: binary.BigEndian.Uint16(b[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(b[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(b[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(b[10:12]),
	}
}

func (h *Headers) Bytes() []byte {
	b := make([]byte, 12)
	binary.BigEndian.PutUint16(b[0:2], h.ID)
	binary.BigEndian.PutUint16(b[2:4], h.CODE)
	binary.BigEndian.PutUint16(b[4:6], h.QDCOUNT)
	binary.BigEndian.PutUint16(b[6:8], h.ANCOUNT)
	binary.BigEndian.PutUint16(b[8:10], h.NSCOUNT)
	binary.BigEndian.PutUint16(b[10:12], h.ARCOUNT)
	return b
}

type Question struct {
	QNAME  Labels
	QTYPE  Type
	QCLASS Class
}

func ParseQuestion(b []byte, offset int) (*Question, int) {
	qname, offset := ParseLabels(b, offset)
	qtype := Type(binary.BigEndian.Uint16(b[offset : offset+2]))
	offset += 2
	qclass := Class(binary.BigEndian.Uint16(b[offset : offset+2]))
	offset += 2
	return &Question{
		QNAME:  qname,
		QTYPE:  qtype,
		QCLASS: qclass,
	}, offset
}

func (question *Question) Bytes() []byte {
	b := question.QNAME.Bytes()
	b = binary.BigEndian.AppendUint16(b, uint16(question.QTYPE))
	b = binary.BigEndian.AppendUint16(b, uint16(question.QCLASS))
	return b
}

type Labels []string

func ParseLabels(b []byte, offset int) (Labels, int) {
	labels := make(Labels, 0)
	for l := b[offset]; l > 0; l = b[offset] {
		if l&0xC0 == 0xC0 { //pointer label
			pointerOffset := int(binary.BigEndian.Uint16([]byte{l, b[offset+1]})) & 0x3FFF
			referencedLabels, _ := ParseLabels(b, pointerOffset)
			return append(labels, referencedLabels...), offset + 2
		} else {
			labels = append(labels, string(b[offset+1:offset+1+int(l)]))
			offset += int(l) + 1
		}
	}
	return labels, offset + 1
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

type ResourceRecord struct {
	NAME     Labels
	TYPE     Type
	CLASS    Class
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

func NewResourceRecord(question *Question) *ResourceRecord {
	return &ResourceRecord{
		NAME:  question.QNAME,
		TYPE:  question.QTYPE,
		CLASS: question.QCLASS,
	}
}

func ParseResourceRecord(b []byte, offset int) (*ResourceRecord, int) {
	name, offset := ParseLabels(b, offset)
	typ := Type(binary.BigEndian.Uint16(b[offset : offset+2]))
	offset += 2
	class := Class(binary.BigEndian.Uint16(b[offset : offset+2]))
	offset += 2
	ttl := binary.BigEndian.Uint32(b[offset : offset+4])
	offset += 4
	rdlength := binary.BigEndian.Uint16(b[offset : offset+2])
	offset += 2
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, b[offset:offset+int(rdlength)])
	offset += int(rdlength)
	return &ResourceRecord{
		NAME:     name,
		TYPE:     typ,
		CLASS:    class,
		TTL:      ttl,
		RDLENGTH: rdlength,
		RDATA:    buf.Bytes(),
	}, offset
}

func (rr *ResourceRecord) SetData(b []byte) {
	rr.RDLENGTH = uint16(len(b))
	rr.RDATA = b
}

func (rr *ResourceRecord) Bytes() []byte {
	b := rr.NAME.Bytes()
	b = binary.BigEndian.AppendUint16(b, uint16(rr.TYPE))
	b = binary.BigEndian.AppendUint16(b, uint16(rr.CLASS))
	b = binary.BigEndian.AppendUint32(b, rr.TTL)
	b = binary.BigEndian.AppendUint16(b, rr.RDLENGTH)
	buf := bytes.NewBuffer(b)
	binary.Write(buf, binary.BigEndian, rr.RDATA)
	return buf.Bytes()
}

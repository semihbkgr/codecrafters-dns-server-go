package main

import (
	"fmt"
	"net"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		message := ParseMessage(buf[:size])
		fmt.Printf("request\n%v\n", message)

		// response
		headers := NewHeaders(message.Headers.ID)
		headers.SetQR(true)
		headers.SetOPCODE(message.Headers.OPCODE())
		headers.SetRD(message.Headers.RD())
		if headers.OPCODE() == 0 {
			headers.SetRCODE(0) // no error
		} else {
			headers.SetRCODE(4) // not implemented
		}

		questions := message.Questions

		answers := make([]*ResourceRecord, 0, len(questions))
		for _, question := range questions {
			resourceRecord := NewResourceRecord(question)
			resourceRecord.TTL = 60
			resourceRecord.SetData([]byte{8, 8, 8, 8})
			answers = append(answers, resourceRecord)
		}

		responseMessage := &Message{
			Headers:   headers,
			Questions: questions,
			Answers:   answers,
		}
		responseMessage.Headers.QDCOUNT = uint16(len(responseMessage.Questions))
		responseMessage.Headers.ANCOUNT = uint16(len(responseMessage.Answers))

		fmt.Printf("response\n%v\n", responseMessage)

		response := responseMessage.Bytes()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

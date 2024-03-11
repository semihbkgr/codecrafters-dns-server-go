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

		responseHeaders := &MessageHeaders{
			ID:      1234,
			QDCOUNT: 1,
		}
		responseHeaders.SetQR(true)
		responseQuestions := []*MessageQuestion{
			{
				QNAME:  []string{"codecrafters", "io"},
				QTYPE:  A,
				QCLASS: IN,
			},
		}

		resourceRecord := NewMessageResourceRecord(responseQuestions[0])
		resourceRecord.TTL = 60
		resourceRecord.SetData([]byte{8, 8, 8, 8})

		responseAnswers := []*MessageResourceRecord{resourceRecord}

		responseMessage := &Message{
			Headers:   responseHeaders,
			Questions: responseQuestions,
			Answers:   responseAnswers,
		}

		fmt.Printf("response\n%v\n", responseMessage)

		response := responseMessage.Bytes()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

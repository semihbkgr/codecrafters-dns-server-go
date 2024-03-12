package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	resolver := flag.String("resolver", "", "resolver address where to forward DNS requests")
	flag.Parse()

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

	forwarder, err := NewForwarder(*resolver)
	if err != nil {
		fmt.Println("Failed to create forwarder:", err)
		return
	}

	buf := make([]byte, 512)
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		message := ParseMessage(buf[:size])
		fmt.Printf("request\n%v\n", message)

		responses := make([]*Message, 0, message.Headers.QDCOUNT)
		for _, question := range message.Questions {
			forwardHeaders := *message.Headers
			forwardHeaders.QDCOUNT = 1
			forwardQuestion := *question
			forwardMessage := &Message{
				Headers:   &forwardHeaders,
				Questions: []*Question{&forwardQuestion},
			}
			response, err := forwarder.Forward(forwardMessage)
			if err != nil {
				fmt.Println("Failed to forward request:", err)
				continue
			}
			responses = append(responses, response)

		}

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

		answers := make([]*ResourceRecord, 0, len(responses))
		for _, resp := range responses {
			answers = append(answers, resp.Answers...)
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

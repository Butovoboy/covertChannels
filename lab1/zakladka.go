package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func strToBits(sendCh chan string) {
	fmt.Printf("Converting message to a string of bits\n")

	bits := ""
	for _, char := range os.Args[1] {
		bits += fmt.Sprintf("%08b", char)
	}
	fmt.Println(bits)
	sendCh <- bits
}

func sendPackets(sendCh chan string, destAddr *net.IPAddr) {
	str := <-sendCh

	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening to ICMP traffic: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Listening to ICMP traffic...")

	// Create a buffer to hold incoming packets.
	receivePacket := make([]byte, 1024)

	// Set up the buffer to hold the incoming data
	buf := bytes.NewBuffer(nil)

	for _, bit := range str {
		if bit == '1' {
			fmt.Print("a")
		} else {
			fmt.Print("b")
		}
	}
	fmt.Println("i111")

	for {
		fmt.Println("qwert")
		n, _, err := conn.ReadFrom(receivePacket)
		if err != nil {
			log.Println(err)
			continue
		}
		buf.Write(receivePacket[:n])

		// buffer is full
		if buf.Len() >= 1024 {
			data := buf.Bytes() // Buffer is full => send data to a new ICMP packet
			buf.Reset()         // Clear the buffer to next data

			msg := icmp.Message{
				Type: ipv4.ICMPTypeEcho,
				Code: 0,
				Body: &icmp.Echo{
					ID:   os.Getpid() & 0xffff,
					Seq:  1,
					Data: data,
				},
			}

			msgBytes, err := msg.Marshal(nil)
			if err != nil {
				log.Fatal(err)
			}

			_, err = conn.WriteTo(msgBytes, destAddr)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	// IPv4 address of destination point (tbh the best variant would be to inspect ICMP packets to get DestIP, but may be later)
	destAddr, err := net.ResolveIPAddr("ip4", "3.72.181.255")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s Write a string to send\n", os.Args[0])
		os.Exit(1)
	}

	// Create a channel to signal when to start sending packets.
	sendCh := make(chan string)

	// Start a goroutine to convert input string to string of bits
	go strToBits(sendCh)

	// Start a goroutine to handle incoming packets.
	go sendPackets(sendCh, destAddr)

	// Wait for signal from user to start sending packets.
	<-sendCh
	fmt.Println("Starting to send packets...")
}

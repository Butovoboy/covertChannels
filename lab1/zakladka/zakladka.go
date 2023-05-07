package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func strToBits(sendCh chan string, message string) {
	fmt.Printf("Converting message to a string of bits\n")

	bits := ""
	for _, char := range message {
		bits += fmt.Sprintf("%08b", char)
	}
	sendCh <- bits
}

func getPackets(sendCh chan string, destAddr *net.IPAddr, delay int) {
	str := <-sendCh
	fmt.Println(str)
	// Wait for signal from user to start sending packets.
	fmt.Println("Starting to send packets...")

	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening to ICMP traffic: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Listening to ICMP traffic...")

	// Create a buffer to hold incoming packets.
	receivePacket := make([]byte, 1024)

	// Set up the buffer to hold the incoming data
	buf := bytes.NewBuffer(nil)

	for {
		n, _, err := conn.ReadFrom(receivePacket)
		if err != nil {
			log.Println(err)
			continue
		}
		buf.Write(receivePacket[:n])

		// buffer is full
		if buf.Len() >= 64 {
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

			// says that covert channel starts working
			go sendPackets(msgBytes, conn, destAddr, '1')

			start := time.Now()

			// sending covert message
			for i := range str {
				for {
					elapsed := time.Since(start)
					if elapsed >= time.Duration(delay)*time.Millisecond { // 3 - is a time delay between each packet sending
						go sendPackets(msgBytes, conn, destAddr, str[i])
						fmt.Printf("DELAY: %v, byte - %v\n", time.Since(start), str[i])
						start = time.Now()
						break
					}
				}
			}
		}
		break
	}
	time.Sleep(3 * time.Second)
}

func sendPackets(msgBytes []byte, conn net.PacketConn, destAddr *net.IPAddr, b byte) {
	if b == '1' {
		_, err := conn.WriteTo(msgBytes, destAddr)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// IPv4 address of destination point (tbh the best variant would be to inspect ICMP packets to get DestIP, but may be later)
	destAddr, err := net.ResolveIPAddr("ip4", "192.168.3.17")
	if err != nil {
		log.Fatal(err)
	}

	//	if len(os.Args) != 2 {
	//		fmt.Fprintf(os.Stderr, "Usage: %s Write a string to send\n", os.Args[0])
	//		os.Exit(1)
	//	}

	// Create a channel to signal when to start sending packets.
	sendCh := make(chan string)

	//message := os.Args[1]
	message := "Hello from Moscow!"
	delay := 1000 // milliseconds
	// Start a goroutine to convert input string to string of bits
	go strToBits(sendCh, message)

	// Start a goroutine to handle incoming packets.
	go getPackets(sendCh, destAddr, delay)
	time.Sleep(300 * time.Second) // need to run goroutines
}

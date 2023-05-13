package main

import (
	"bytes"
	"crypto/rand"
	"flag"
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

func getPackets(sendCh chan string, destAddr *net.IPAddr, delay int, traffic bool) {
	var conn net.PacketConn
	var err error

	str := <-sendCh

	conn, err = net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening to ICMP traffic: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Starting sending packets by covert channel ...")
	defer conn.Close()

	// Create a buffer to hold incoming packets.
	receivePacket := make([]byte, 4096)

	// Set up the buffer to hold the incoming data
	buf := bytes.NewBuffer(nil)

	for {
		if traffic {
			n, _, err := conn.ReadFrom(receivePacket)
			if err != nil {
				log.Println(err)
				continue
			}
			buf.Write(receivePacket[:n])
		} else {
			b := make([]byte, 128)
			// Fill the byte slice with random bytes
			_, err := rand.Read(b)
			if err != nil {
				panic(err)
			}
			buf.Write(b)
		}

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
					if elapsed >= time.Duration(delay)*time.Millisecond {
						go sendPackets(msgBytes, conn, destAddr, str[i])
						start = time.Now()
						break
					}
				}
			}
		}
		break
	}

	// need to send the last 1 to say that the covert channel is closing and these 0 at the end were last
	time.Sleep(time.Duration(delay) * time.Millisecond)
	go sendPackets(make([]byte, 1024), conn, destAddr, '1')
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

	traffic := flag.Bool("traffic", false, "A string to decide to get ICMP trafic or to generate your own")

	flag.Parse()

	// Create a channel to signal when to start sending packets.
	sendCh := make(chan string)

	var message string

	if *traffic {
		message = os.Args[2]
	} else {
		message = os.Args[1]
	}
	fmt.Printf("MESSAGE: %v\n", message)
	delay := 1000 // milliseconds

	// Start a goroutine to convert input string to string of bits
	go strToBits(sendCh, message)

	// Start a goroutine to handle incoming packets.
	go getPackets(sendCh, destAddr, delay, *traffic)

	time.Sleep(300 * time.Second) // need to run goroutines
}

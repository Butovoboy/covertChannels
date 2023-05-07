package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {

	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening to ICMP traffic: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Listening to ICMP traffic...")

	// Create a buffer to hold incoming packets.
	receivePacket := make([]byte, 1024)

	start := time.Now()
	for {
		n, _, err := conn.ReadFrom(receivePacket)

		elapsed := time.Since(start)
		if elapsed >= 2920*time.Millisecond { // 3 - is a time delay between each packet sending
			if err != nil {
				log.Println(err)
				continue
			}

			fmt.Printf("data - %v\n", n)
			if n > 0 {
				fmt.Printf("It was 1\n")
				start = time.Now()
				continue
			}
			fmt.Printf("It was 0\n")
			start = time.Now()
		} else {
			continue
		}
	}
}

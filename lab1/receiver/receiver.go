package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func binaryToASCII(binary string) error {
	var ascii string

	for i := 0; i < len(binary); i += 8 {
		if i+8 > len(binary) {
			return fmt.Errorf("binary string length must be a multiple of 8")
		}

		value, err := strconv.ParseInt(binary[i:i+8], 2, 64)
		if err != nil {
			return err
		}

		ascii += fmt.Sprintf("%c", value)
	}
	fmt.Printf(ascii)
	return nil
}

func cutByte(s string) (string, string) {
	if len(s) < 8 {
		return "", ""
	}
	return s[8:], s[:8]
}

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

	binary := ""

	_, _, err = conn.ReadFrom(receivePacket)

	start := time.Now()

	for {
		_, _, err = conn.ReadFrom(receivePacket)

		elapsed := time.Since(start)

		if elapsed >= 2990*time.Millisecond { // 3 - is a time delay between each packet sending

			zeros := int(time.Since(start).Round(time.Second).Seconds()) / 3
			for i := 1; i < zeros; i++ {
				binary = binary + "0"
				fmt.Printf("Now binary is %v\n", binary)
			}

			if err != nil {
				log.Println(err)
				continue
			}
			binary = binary + "1"
			fmt.Printf("Now binary is %v\n", binary)
			if len(binary) >= 8 {
				rest, character := cutByte(binary)
				binary = rest
				fmt.Printf("Now binary is %v\n", binary)
				go binaryToASCII(character)
			}
			start = time.Now()
		}
	}
}

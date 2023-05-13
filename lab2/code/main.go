package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcapgo"
)

const (
	PacketID     = 100
	TimeInterval = 1 * time.Second
)

func decodeMessage(dumpFile string) (string, error) {
	f, err := os.Open(dumpFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r, err := pcapgo.NewReader(f)
	if err != nil {
		return "", err
	}

	blocks := []gopacket.Packet{}
	packetSource := gopacket.NewPacketSource(r, r.LinkType())
	for packet := range packetSource.Packets() {
		blocks = append(blocks, packet)
	}

	err = utils.show_gaps(blocks)
	if err != nil {
		return "", err
	}

	message := ""
	for counter := PacketID; counter < len(blocks); counter++ {
		difference := blocks[counter].Metadata().Timestamp.Sub(blocks[counter-1].Metadata().Timestamp) // counts the difference in timestamps between neighbour blocks
		message += strconv.Itoa(int(difference.Round(time.Second)) / 1000000000)                       // reading message in bits
	}

	res, err := binaryToASCII(message)
	if err != nil {
		return "", err
	}

	return res, nil
}

func binaryToASCII(binary string) (string, error) {
	var ascii string

	for i := 0; i < len(binary); i += 8 {
		if i+8 > len(binary) {
			return ascii, fmt.Errorf("binary string length must be a multiple of 8")
		}

		value, err := strconv.ParseInt(binary[i:i+8], 2, 64)
		if err != nil {
			return ascii, err
		}

		ascii += fmt.Sprintf("%c", value)
	}
	return ascii, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: program <pcapng file>")
		os.Exit(1)
	}

	message, err := decodeMessage(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Message:", message)
}

package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcapgo"

	"decoder/utils"
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

	if err != nil {
		return "", err
	}

	message := ""
	for counter := PacketID; counter < len(blocks); counter++ {
		message += strconv.Itoa(int(countDifference(blocks, counter).Round(time.Second)) / 1000000000) // reading message in bits
	}

	createGraphic(blocks, 1, "all_packets")
	createGraphic(blocks, PacketID, "covert_packets")
	res, err := binaryToASCII(message)
	if err != nil {
		return "", err
	}

	return res, nil
}

// counts the difference in timestamps between neighbour blocks
func countDifference(blocks []gopacket.Packet, counter int) time.Duration {
	difference := blocks[counter].Metadata().Timestamp.Sub(blocks[counter-1].Metadata().Timestamp)
	return difference
}

// sorts keys in increasing order
func sortKeys(intervalsMap map[int]int) []int {
	// Extract the keys into a slice
	keys := make([]int, 0, len(intervalsMap))
	for key := range intervalsMap {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	return keys
}

// get a map with numbers of packats with each interval
func createGraphic(blocks []gopacket.Packet, strt int, name string) error {
	intervalsMap := make(map[int]int)
	for counter := strt; counter < len(blocks); counter++ {
		interval := int(countDifference(blocks, counter).Round(250*time.Millisecond) / 1000000)
		_, exists := intervalsMap[interval]
		if exists {
			intervalsMap[interval] += 1
		} else {
			intervalsMap[interval] = 1
		}
	}

	err := utils.Show_gaps(intervalsMap, sortKeys(intervalsMap), name)
	if err != nil {
		return err
	}

	return nil
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

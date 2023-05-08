package main

import (
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func changeSpeed(destAddr *net.IPAddr, delay int) {
	// Open network interface for listening to packets
	handle, err := pcap.OpenLive("wlp0s20f3", 65535, true, 1*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set BPF filter for ICMP packets
	filter := "icmp"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	// Start packet processing loop
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Extract ICMP layer from packet
		icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
		if icmpLayer == nil {
			continue
		}
		icmp, ok := icmpLayer.(*layers.ICMPv4)
		if !ok {
			continue
		}

		time.Sleep(time.Duration(delay) * time.Second)

		// Create new IP packet with the same payload
		newIP := layers.IPv4{
			SrcIP:    net.IPv4(127, 0, 0, 1),
			DstIP:    net.ParseIP(destAddr.IP.String()),
			Protocol: layers.IPProtocolICMPv4,
		}
		newPacket := gopacket.NewSerializeBuffer()
		err = gopacket.SerializeLayers(newPacket, gopacket.SerializeOptions{},
			&newIP, icmp)
		if err != nil {
			log.Fatal(err)
		}

		// Send the new packet to the destination IP address
		conn, err := net.Dial("ip4:icmp", destAddr.IP.String())
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write(newPacket.Bytes())
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

	delay := 2 // seconds

	go changeSpeed(destAddr, delay)

	time.Sleep(300 * time.Second) // need to run goroutines
}

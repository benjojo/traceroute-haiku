package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"

	"github.com/songgao/water"
)

func main() {
	log.Printf("Haiku Tun")
	ifacename := flag.String("iface", "haiku0", "the name of the interface")
	flag.Parse()

	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = *ifacename

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}

	for {
		packet := make([]byte, 1500)
		plen, err := ifce.Read(packet)
		if err != nil {
			log.Fatalf("Tun spluttered, Threw back %s on the tun file", err.Error())
		}

		if plen < 40 {
			log.Printf("dropped packet becase it was too small %d bytes", plen)
			continue // packet too small to be real, drop it
		}

		if packet[39] != 0x04 {
			log.Printf("dropped packet because it didnt match the last byte policy")
			continue // drop and ignore, we don't care
		}

		TTL := uint8(packet[7])
		if TTL < 4 {
			log.Printf("debug: TTL is less than 4")
			// Handle ICMP Hop limit expired
			returnpacket := make([]byte, plen+8+40)
			returnpacket[0] = 0x60 // IP packet version, Thus it is 6
			// <---> Flow labels and crap like that here, leaving this as zeros
			uintbuf := new(bytes.Buffer)
			binary.Write(uintbuf, binary.BigEndian, uint16(plen+8))
			returnpacket[4] = uintbuf.Bytes()[0] // Packetlength 1/2
			returnpacket[5] = uintbuf.Bytes()[1] // Packetlength 2/2
			returnpacket[6] = 0x3a               // Next header (aka packet content protocol), 0x3a == 58 == ICMPv6
			returnpacket[7] = 0x40               // Hop Limit of the outgoing packet, 0x40 == 64
			// Set the source to the destination of the packet, but morph it based on the TTL
			for i := 0; i < 16; i++ {
				returnpacket[7+i] = packet[23+i]
			}
			returnpacket[7+16] = 0x00 + TTL
			// Copy the source address from the incoming packet, and use it as the destination.
			for i := 0; i < 16; i++ {
				returnpacket[23+i] = packet[7+i]
			}
			// Now, We should have a IPv6 packet that "works", now we need to make the ICMPv6 chunk
			// and checksum it.

			returnpacket[40] = 0x03 // Time Exceeded
			returnpacket[41] = 0x00 // Hop limit exceeded in transit
			// two bytes here are used for CRC, we will do that in a sec.
			returnpacket[44] = 0x00 // "Reserved"
			returnpacket[45] = 0x00 // "Reserved"
			returnpacket[46] = 0x00 // "Reserved"
			returnpacket[47] = 0x00 // "Reserved"
			for i := 0; i < plen; i++ {
				returnpacket[48+i] = packet[i]
			}

			// Oh GOD now here comes a strange CRC dance
			crc, err := computeChecksum(returnpacket)
			if err != nil {
				log.Printf("Failed to CRC: %s", err.Error())
			}
			crcbuf := new(bytes.Buffer)
			binary.Write(crcbuf, binary.BigEndian, uint16(crc))
			// crcbuf := make([]byte, binary.MaxVarintLen16)
			// binary.PutUvarint(crcbuf, uint64(crc))
			returnpacket[42] = crcbuf.Bytes()[0]
			returnpacket[43] = crcbuf.Bytes()[1]

			// Aaaaaaaaaaaand that is all folks, Send it out. Ship it.
			log.Printf("debug: all done sending.")

			ifce.Write(returnpacket)
		} else {
			// Handle generic responce (UDP and ICMP)
		}

	}
}

func computeChecksum(headerAndPayload []byte) (uint16, error) {
	length := uint32(len(headerAndPayload))
	csum, err := pseudoheaderChecksum(headerAndPayload[7:7+16], headerAndPayload[23:23+16])
	if err != nil {
		return 0, err
	}
	csum += uint32(58)
	csum += length & 0xffff
	csum += length >> 16
	return tcpipChecksum(headerAndPayload[39:], csum), nil
}

func pseudoheaderChecksum(SrcIP, DstIP []byte) (csum uint32, err error) {
	for i := 0; i < 16; i += 2 {
		csum += uint32(SrcIP[i]) << 8
		csum += uint32(SrcIP[i+1])
		csum += uint32(DstIP[i]) << 8
		csum += uint32(DstIP[i+1])
	}
	return csum, nil
}

func tcpipChecksum(data []byte, csum uint32) uint16 {
	// to handle odd lengths, we loop to length - 1, incrementing by 2, then
	// handle the last byte specifically by checking against the original
	// length.
	length := len(data) - 1
	for i := 0; i < length; i += 2 {
		// For our test packet, doing this manually is about 25% faster
		// (740 ns vs. 1000ns) than doing it by calling binary.BigEndian.Uint16.
		csum += uint32(data[i]) << 8
		csum += uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		csum += uint32(data[length]) << 8
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum)
}

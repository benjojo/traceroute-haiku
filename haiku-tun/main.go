package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"net"

	"github.com/songgao/water"
	"golang.org/x/net/icmp"
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
		if TTL < 5 {
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
			for i := 0; i < 17; i++ {
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
			// crc := v6checksumICMP(returnpacket)
			src := net.IP(returnpacket[8 : 8+16])
			dst := net.IP(returnpacket[24 : 24+16])
			log.Printf("Returning fire from %s to %s", src.String(), dst.String())
			crcb := Checksum(returnpacket[40:], src, dst)
			// crcbuf := new(bytes.Buffer)
			// binary.Write(crcbuf, binary.BigEndian, uint16(crc))
			returnpacket[42] = crcb[0]
			returnpacket[43] = crcb[1]

			// Aaaaaaaaaaaand that is all folks, Send it out. Ship it.
			log.Printf("debug: all done sending.")

			ifce.Write(returnpacket)
		} else {
			// Handle generic responce (UDP and ICMP)
		}

	}
}

func Checksum(body []byte, srcIP, dstIP net.IP) (crc []byte) {
	out := make([]byte, 2)
	// from golang.org/x/net/icmp/message.go
	checksum := func(b []byte) uint16 {
		csumcv := len(b) - 1 // checksum coverage
		s := uint32(0)
		for i := 0; i < csumcv; i += 2 {
			s += uint32(b[i+1])<<8 | uint32(b[i])
		}
		if csumcv&1 == 0 {
			s += uint32(b[csumcv])
		}
		s = s>>16 + s&0xffff
		s = s + s>>16
		return ^uint16(s)
	}

	b := body

	// remember origin length
	l := len(b)
	// generate pseudo header
	psh := icmp.IPv6PseudoHeader(srcIP, dstIP)
	// concat psh with b
	b = append(psh, b...)
	// set length of total packet
	off := 2 * net.IPv6len
	binary.BigEndian.PutUint32(b[off:off+4], uint32(l))
	// calculate checksum
	s := checksum(b)
	// set checksum in bytes and return original Body
	out[0] ^= byte(s)
	out[1] ^= byte(s >> 8)

	return out
}

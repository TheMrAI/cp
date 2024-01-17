package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func HexStringToBinary(hexString string) []byte {
	binarySequence := make([]byte, 0, len(hexString)*4)
	for _, r := range hexString {
		binary, err := HexToBinary(r)
		if err != nil {
			log.Fatal(err)
		}
		binarySequence = append(binarySequence, binary...)
	}
	return binarySequence
}

// A little hacky, but I see no other simple way of parsing
// the values quickly, without having to implement additional methods.
func HexToBinary(r rune) ([]byte, error) {
	hexToBinary := map[rune][]byte{
		'0': {'0', '0', '0', '0'},
		'1': {'0', '0', '0', '1'},
		'2': {'0', '0', '1', '0'},
		'3': {'0', '0', '1', '1'},
		'4': {'0', '1', '0', '0'},
		'5': {'0', '1', '0', '1'},
		'6': {'0', '1', '1', '0'},
		'7': {'0', '1', '1', '1'},
		'8': {'1', '0', '0', '0'},
		'9': {'1', '0', '0', '1'},
		'A': {'1', '0', '1', '0'},
		'B': {'1', '0', '1', '1'},
		'C': {'1', '1', '0', '0'},
		'D': {'1', '1', '0', '1'},
		'E': {'1', '1', '1', '0'},
		'F': {'1', '1', '1', '1'},
	}
	binarySequence, ok := hexToBinary[r]
	if !ok {
		return []byte{0, 0, 0, 0}, fmt.Errorf("Invalid rune received in hex conversion: '%c'", r)
	}
	return binarySequence, nil
}

type Packet struct {
	version    uint8
	packetType uint8
	value      uint64
	subPackets []Packet
}

func ParseInput(packetStr string) Packet {
	bytes := HexStringToBinary(packetStr)
	bytesRead, packet, err := ParsePacket(bytes, 0)
	if err != nil {
		log.Fatalf("Failed to parse packets after %v bit(s), with error \n\t%v", bytesRead, err)
	}
	return packet
}

func ParsePacket(bytes []byte, index int) (int, Packet, error) {
	version, err := strconv.ParseUint(string(bytes[index:index+3]), 2, 8)
	if err != nil {
		return index, Packet{}, fmt.Errorf("version parse failure")
	}
	index += 3

	packetType, err := strconv.ParseUint(string(bytes[index:index+3]), 2, 8)
	if err != nil {
		return index, Packet{}, fmt.Errorf("packetType parse failure")
	}
	index += 3

	packet := Packet{uint8(version), uint8(packetType), 0, []Packet{}}

	if packetType == 4 {
		index, packet.value, err = ParseLiteral(bytes, index)
		if err != nil {
			return index, packet, err
		}
		return index, packet, err
	} else {
		lengthType := bytes[index]
		index++
		if lengthType == '0' {
			subPacketBitLength, err := strconv.ParseUint(string(bytes[index:index+15]), 2, 64)
			index += 15
			if err != nil {
				return index, packet, fmt.Errorf("sub packet length parse failure")
			}
			end := index + int(subPacketBitLength)
			for index < end {
				nextIndex, subPacket, err := ParsePacket(bytes, index)
				index = nextIndex
				if err != nil {
					return index, packet, fmt.Errorf("packet parse failure")
				}
				packet.subPackets = append(packet.subPackets, subPacket)
			}
			return end, packet, nil
		} else {
			subPacketCount, err := strconv.ParseUint(string(bytes[index:index+11]), 2, 64)
			index += 11

			if err != nil {
				return index, packet, fmt.Errorf("sub packet count parse failure")
			}
			for count := 0; count < int(subPacketCount); count++ {
				nextIndex, subPacket, err := ParsePacket(bytes, index)
				if err != nil {
					return nextIndex, packet, fmt.Errorf("packet parse failure")
				}
				packet.subPackets = append(packet.subPackets, subPacket)
				index = nextIndex
			}
			return index, packet, nil
		}
	}
}

func ParseLiteral(bytes []byte, index int) (int, uint64, error) {
	numberSequence := []byte{}
	for ; index <= len(bytes)-5; index += 5 {
		numberSequence = append(numberSequence, bytes[index+1:index+5]...)
		if bytes[index] == '0' {
			index += 5
			break
		}
	}

	value, err := strconv.ParseUint(string(numberSequence), 2, 64)
	if err != nil {
		return index, value, err
	}
	return index, value, nil
}

func PartOne(packet Packet) int {
	return sumAllVersionNumbers(packet)
}

func sumAllVersionNumbers(packet Packet) int {
	versionNumberSum := 0
	for i := range packet.subPackets {
		versionNumberSum += sumAllVersionNumbers(packet.subPackets[i])
	}
	return versionNumberSum + int(packet.version)
}

func PartTwo(packet Packet) uint64 {
	return evaluate(packet)
}

func evaluate(packet Packet) uint64 {
	switch packet.packetType {
	case 0:
		{
			sum := evaluate(packet.subPackets[0])
			for i := 1; i < len(packet.subPackets); i++ {
				sum += evaluate(packet.subPackets[i])
			}
			return sum
		}
	case 1:
		{
			sum := evaluate(packet.subPackets[0])
			for i := 1; i < len(packet.subPackets); i++ {
				sum *= evaluate(packet.subPackets[i])
			}
			return sum
		}
	case 2:
		{
			minimum := evaluate(packet.subPackets[0])
			for i := 1; i < len(packet.subPackets); i++ {
				val := evaluate(packet.subPackets[i])
				if val < minimum {
					minimum = val
				}
			}
			return minimum
		}
	case 3:
		{
			maximum := evaluate(packet.subPackets[0])
			for i := 1; i < len(packet.subPackets); i++ {
				val := evaluate(packet.subPackets[i])
				if val > maximum {
					maximum = val
				}
			}
			return maximum
		}
	case 4:
		{
			return packet.value
		}
	case 5:
		{
			if evaluate(packet.subPackets[0]) > evaluate(packet.subPackets[1]) {
				return 1
			}
			return 0
		}
	case 6:
		{
			if evaluate(packet.subPackets[0]) < evaluate(packet.subPackets[1]) {
				return 1
			}
			return 0
		}
	case 7:
		{
			if evaluate(packet.subPackets[0]) == evaluate(packet.subPackets[1]) {
				return 1
			}
			return 0
		}
	}
	return 0
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	for i := range lines {
		packet := ParseInput(lines[i])
		fmt.Println("Part one")
		fmt.Printf("Sum of all version numbers in the packet: %v\n", PartOne(packet))
		fmt.Println("Part two")
		fmt.Printf("Value of the packet: %v\n", PartTwo(packet))
	}
}

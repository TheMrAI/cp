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
		log.Fatalf("Failed to parse packets after %v bytes, with error \n\t%v", bytesRead, err)
	}
	return packet
}

func ParsePacket(bytes []byte, index int) (int, Packet, error) {
	version, err := strconv.ParseUint(string(bytes[index:index+3]), 2, 8)
	if err != nil {
		return index, Packet{}, fmt.Errorf("Failed to parse version starting at index: %v", index)
	}
	index += 3

	packetType, err := strconv.ParseUint(string(bytes[index:index+3]), 2, 8)
	if err != nil {
		return index, Packet{}, fmt.Errorf("Failed to parse packetType starting at index: %v", index)
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
				return index, packet, err
			}
			end := index + int(subPacketBitLength)
			for index < end {
				prevIndex := index
				nextIndex, subPacket, err := ParsePacket(bytes, index)
				index = nextIndex
				if err != nil {
					return index, packet, err
				}
				if prevIndex == index {
					break
				}
				packet.subPackets = append(packet.subPackets, subPacket)
			}
			return end, packet, nil
		} else {
			subPacketCount, err := strconv.ParseUint(string(bytes[index:index+11]), 2, 64)
			index += 11
			if err != nil {
				return index, packet, err
			}
			for count := 0; count < int(subPacketCount); count++ {
				nextIndex, subPacket, err := ParsePacket(bytes, index)
				if err != nil {
					return nextIndex, packet, err
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
	for ; index < len(bytes)-5; index += 5 {
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
	}
}

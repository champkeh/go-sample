package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"bytes"
)

// ICMP represent icmp protocol datagram
type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host")
		os.Exit(1)
	}

	addr, err := net.ResolveIPAddr("ip", os.Args[1])
	if err != nil {
		fmt.Println("Resolution error", err.Error())
		os.Exit(1)

	}

	conn, err := net.DialIP("ip4:icmp", nil, addr)
	checkError(err)

	icmp := ICMP{
		Type:        8,
		Code:        0,
		Checksum:    0,
		Identifier:  14,
		SequenceNum: 237,
	}

	var buffer bytes.Buffer

	// 将icmp数据报按照大端序写入buffer,进行校验和的计算
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = checkSum(buffer.Bytes())

	// 然后清空buffer，将计算完校验和的icmp数据报重新写入buffer
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)

	_, err = conn.Write(buffer.Bytes())
	checkError(err)

	msg := make([]byte, 512)
	n, err := conn.Read(msg)
	checkError(err)

	fmt.Printf("Got response %d bytes\n", n)
	// 这里会返回28个字节，不知道是什么原因
	// 这里进行截断，取最后8个字节
	if n >= 8 {
		n -= 8
	}
	if msg[n+5] == 14 {
		fmt.Println("identifier matches")
	}
	if msg[n+7] == 237 {
		fmt.Println("Sequence matches")
	}

	os.Exit(0)
}

// 计算checksum
func checkSum(data []byte) uint16 {
	var (
		sum    uint32
		length = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)
	return uint16(^sum)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

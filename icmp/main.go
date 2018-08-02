package main

import (
	"fmt"
	"net"
	"os"
)

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

	var msg [512]byte
	msg[0] = 8    // type
	msg[1] = 0    // code
	msg[2] = 0    // checksum, fix later
	msg[3] = 0    // checksume, fix later
	msg[4] = 0    // identifier[0]
	msg[5] = 14   // identifier[1]
	msg[6] = 0    // sequence[0]
	msg[7] = 237  // sequence[1]
	len := 8

	sum := checkSum(msg[0:len])
	msg[2] = byte(sum >> 8)
	msg[3] = byte(sum & 255)

	_, err = conn.Write(msg[0:len])
	checkError(err)

	n, err := conn.Read(msg[0:])
	checkError(err)

	fmt.Printf("Got response %d bytes\n", n)
	// 在windows上测试，会返回28个字节
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

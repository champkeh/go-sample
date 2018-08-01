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
	msg[2] = 0xf7 // checksum, fix later, 此处直接使用正确的checksum
	msg[3] = 0xcd // checksume, fix later, 此处直接使用正确的checksum
	msg[4] = 0    // identifier[0]
	msg[5] = 13   // identifier[1]
	msg[6] = 0    // sequence[0]
	msg[7] = 37   // sequence[1]
	len := 8

	//sum := checkSum(msg[0:len])
	//msg[2] = byte(sum >> 8)
	//msg[3] = byte(sum & 255)

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
	if msg[n+5] == 13 {
		fmt.Println("identifier matches")
	}
	if msg[n+7] == 37 {
		fmt.Println("Sequence matches")
	}

	os.Exit(0)
}

// 计算checksum
// note: 算法好像有问题，计算的结果不正确
func checkSum(msg []byte) uint16 {
	sum := 0

	for n := 1; n < len(msg)-1; n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer = uint16(^sum)
	return answer
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

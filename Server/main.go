package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
)

// 使用给定的码片序列编码信息
func encode(encoded *([]int), chip []int, msg *([]byte)) {
	chipLen := len(chip)
	for i, v := range *msg {
		for j := 0; j < 8; j++ {
			bit := (v >> uint(j)) & 1
			if bit == 1 {
				(*encoded)[8*chipLen*i+chipLen*j] += chip[0]
				(*encoded)[8*chipLen*i+chipLen*j+1] += chip[1]
			} else {
				(*encoded)[8*chipLen*i+chipLen*j] += -chip[0]
				(*encoded)[8*chipLen*i+chipLen*j+1] += -chip[1]
			}
		}
	}
}

// 使用给定的码片序列解码信息
// func decode(decoded *([]byte), chip []int, encoded *([]int)) {
// 	chipLen := len(chip)
// 	if len(*encoded)/chipLen/8 != len(*decoded) {
// 		fmt.Println("Insufficient space to decode message!")
// 		return
// 	}
// 	for i := 0; i < len(*encoded)/chipLen/8; i++ {
// 		for j := 0; j < 8; j++ {
// 			temp := 0
// 			for k := 0; k < chipLen; k++ {
// 				temp += (*encoded)[i*8*chipLen+j*chipLen+k] * chip[k]
// 			}
// 			temp /= chipLen
// 			if temp > 0 {
// 				(*decoded)[i] |= (1 << uint(j))
// 			} else {
// 				(*decoded)[i] |= (0 << uint(j))
// 			}
// 		}
// 	}
// }

func main() {

	// （不重要）多机演示，需要给定客户端 IP 地址
	ipA := ""
	fmt.Println("Input Client A's IP:")
	fmt.Scanln(&ipA)

	ipB := ""
	fmt.Println("Input Client B's IP:")
	fmt.Scanln(&ipB)

	UDPAddrA, err := net.ResolveUDPAddr("udp", ipA+":9090")
	if err != nil {
		fmt.Println(err)
		return
	}

	UDPAddrB, _ := net.ResolveUDPAddr("udp", ipB+":9090")
	if err != nil {
		fmt.Println(err)
		return
	}

	udpConnA, err := net.DialUDP("udp", nil, UDPAddrA)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
	}
	defer udpConnA.Close()

	udpConnB, err := net.DialUDP("udp", nil, UDPAddrB)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
	}
	defer udpConnB.Close()

	defer fmt.Println("Done!")

	in := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("-----------------------------------")
		fmt.Println("Please input message sent to A:")
		msgToA, _ := in.ReadBytes('\n')

		fmt.Println("Please input message sent to B:")
		msgToB, _ := in.ReadBytes('\n')

		chipLen := 2
		chipA := []int{1, 1}
		chipB := []int{1, -1}

		fmt.Println("Now we shall use Walsh Code to encode the message.")
		encoded := make([]int, int(float64(chipLen)*8*math.Max(float64(len(msgToA)), float64(len(msgToB)))))
		encode(&encoded, chipA, &msgToA)
		encode(&encoded, chipB, &msgToB)

		fmt.Println("The encoded message is:")
		fmt.Println(encoded)

		fmt.Println("Now sending the message to clients...")
		msgSend, _ := json.Marshal(encoded)

		if _, err := udpConnA.Write(msgSend); err != nil {
			fmt.Println(err)
		}

		if _, err := udpConnB.Write(msgSend); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Sent: ", string(msgSend))
	}
}

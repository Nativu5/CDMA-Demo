package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func decode(decoded *([]byte), chip []int, encoded *([]int)) {
	chipLen := len(chip)
	if len(*encoded)/chipLen/8 != len(*decoded) {
		fmt.Println("Insufficient space to decode message!")
		return
	}
	for i := 0; i < len(*encoded)/chipLen/8; i++ {
		for j := 0; j < 8; j++ {
			temp := 0
			for k := 0; k < chipLen; k++ {
				temp += (*encoded)[i*8*chipLen+j*chipLen+k] * chip[k]
			}
			temp /= chipLen
			if temp > 0 {
				(*decoded)[i] |= (1 << uint(j))
			} else {
				(*decoded)[i] |= (0 << uint(j))
			}
		}
	}
}

func main() {
	var chipLen int
	fmt.Println("Input chip sequence length:")
	fmt.Scanln(&chipLen)

	chip := make([]int, chipLen)
	fmt.Println("Input the chip sequence:")
	for i := 0; i < chipLen; i++ {
		fmt.Scanf("%d", &chip[i])
	}

	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9090,
	})

	if err != nil {
		log.Fatal("Listen failed,", err)
		return
	}
	defer udpConn.Close()

	fmt.Println("Waiting for connection...")
	// 循环读取消息
	for {
		var data [4096]byte
		n, addr, err := udpConn.ReadFromUDP(data[:])
		if err != nil {
			log.Printf("Read from udp server:%s failed, err:%s", addr, err)
			break
		}
		go func() {
			// 返回数据
			fmt.Printf("Addr:%s, data:%v, count:%d\n", addr, string(data[:n]), n)
			fmt.Println("Decoded data:")

			encoded := make([]int, n)
			json.Unmarshal(data[:n], &encoded)

			decoded := make([]byte, len(encoded)/chipLen/8)
			decode(&decoded, chip, &encoded)

			fmt.Println(string(decoded))
		}()
	}
}

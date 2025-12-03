package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	udpConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	r := bufio.NewReader(os.Stdin)

	data := ""
	for {
		fmt.Println(">")
		data, err = r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		udpConn.Write([]byte(data))
	}
}

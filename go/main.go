package main

import (
	"fmt"
	"net"
)

func main() {
	// initialize my listener and close it later
	listener, err := net.Listen("tcp", "localhost:8008")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// loop waiting to accept connections
	for {
		// get the netconnection (net.conn)
		netConnection, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		// split off for different concurrent routines if connections come in
		go func(connection net.Conn) {
			// data to hold input
			data := make([]byte, 1024)

			// read the input into data, and count the bytes read
			bytesRead, err := connection.Read(data)
			if err != nil {
				fmt.Println(err)
			}

			// get the data into a usable slice
			dataSlice := []byte{}
			for idx := range bytesRead {
				dataSlice = append(dataSlice, data[idx])
			}

			wroteData, err := netConnection.Write(dataSlice)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%d, %b", wroteData, dataSlice)
		}(netConnection)
	}
}

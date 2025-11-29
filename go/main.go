package main

import (
	"fmt"
	"io"
	"net"
	"os"
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

			for {
				// read the input into data, and count the bytes read
				bytesRead, err := connection.Read(data)
				if err != nil {
					fmt.Println(err)
					if err == io.EOF {
						fmt.Printf("goodbye :)")
						os.Exit(0)
					}
				}

				// check if there has been a response from client
				if bytesRead == 0 {
					continue
				}
				// get the data into a usable slice
				dataSlice := []byte{}
				for idx := range bytesRead {
					dataSlice = append(dataSlice, data[idx])
				}

				// clear out the data slice
				data = data[:0]

				wroteData, err := connection.Write(dataSlice)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%d, %b", wroteData, dataSlice)
			}
		}(netConnection)
	}
}

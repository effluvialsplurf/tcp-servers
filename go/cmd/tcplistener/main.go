package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	// this will be our return channel
	out := make(chan string, 1)

	go func() {
		// close the file and the channel when they are no longer needed
		defer f.Close()
		defer close(out)

		// initialize the string variable
		str := ""
		// for loop to iterate over the contents of the file
		for {
			// this holds our data
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			// constrain data to the length of the available contents
			data = data[:n]
			// if we are a newline or we are -1
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				// make data empty again using the contents and memory space of the variable
				data = data[i+1:]
				// throw the contents of str in the channel
				out <- str
				str = ""
			}

			// output the string
			str += string(data)
		}

		if len(str) != 0 {
			// if there is anything leftover throw it in the channel
			out <- str
		}
	}()

	return out
}

func main() {
	// grab the listener
	listener, err := net.Listen("tcp", ":42069")
	defer func() {
		listener.Close()
		fmt.Println("conn closed")
	}()
	// if opening fails log the error
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("connected: %s", conn)

		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
	}
}

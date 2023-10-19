package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	var reqBody []byte;
	var escape = "\r\n"

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	req, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	_, err = req.Read(reqBody)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

	s := strings.Split(string(reqBody), escape)

	var resp string;
	if(strings.Split(s[0], " ")[1] == "/") {
		resp = "HTTP/1.1 200 OK\r\n\r\n"
	} else {
		resp = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = req.Write([]byte(resp))
	if err != nil {
		fmt.Println("Error writing request: ", err.Error())
		os.Exit(1)
	}
}

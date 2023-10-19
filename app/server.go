package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	var escape = "\r\n"

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	fmt.Println("listening on Port 4221")
	
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connection Accepted")

	buf, err := io.ReadAll(conn)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

  var resp = "HTTP/1.1 200 OK\r\n\r\n";

	s := strings.Split(string(buf), escape)

	for _, line := range s {
		fmt.Println(line)
	}

	if(strings.Split(s[0], " ")[1] != "/") {
		resp = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	fmt.Println("Writing response", resp)

	_, err = conn.Write([]byte(resp))
	if err != nil {
		fmt.Println("Error writing request: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Response sent")
}

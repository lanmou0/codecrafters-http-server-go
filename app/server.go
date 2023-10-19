package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const ROOT_PATH = "/"
const ECHO_PATH = "/echo/"

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
	
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Connection Accepted")
	
		var buf = make([]byte, 50, 2048)
	
		reader := bufio.NewReader(conn)
	
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			fmt.Println("Error reading request: ", err.Error())
			os.Exit(1)
		}
	
		var resp = buildResponse(404, "Not Found", make(map[header]string), "")

	
		s := strings.Split(string(buf), escape)
		path := strings.Split(s[0], " ")[1]
	
		if(path == "/") {
			resp = buildResponse(200, "OK", make(map[header]string), "")
		}else if(strings.HasPrefix(path, ECHO_PATH)) {
			message := strings.TrimPrefix(path, ECHO_PATH)
			headers := map[header]string{contentType: "text/plain", contentLength: strconv.Itoa(len(message))}
			resp = buildResponse(200, "OK", headers, message)
		}
		fmt.Println("Writing response")
		fmt.Println(resp)
	
		_, err = conn.Write([]byte(resp))
		if err != nil {
			fmt.Println("Error writing request: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Response sent")

		fmt.Println("Connection Closed\r\n")
		conn.Close()
	}
}

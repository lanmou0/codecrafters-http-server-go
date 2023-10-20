package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const ROOT_PATH = "/"
const ECHO_PATH = "/echo/"
const NIL_PATH = "nil"
var ESCAPE = "\r\n"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	createController("nil", func (req HttpRequest) string {
		message := fmt.Sprintf("Path Not Found: %s", req.path)
		headers := make(map[HttpHeader]string)
		headers[ContentType] = "text/plain"
		headers[ContentLength] = strconv.Itoa(len(message))
		return buildResponse(404, "Not Found", headers, message)
	})

	createController("/", func (req HttpRequest) string {
		return buildResponse(200, "OK", make(map[HttpHeader]string), "")
	})

	createController("/echo", func(req HttpRequest) string {
		message := strings.TrimPrefix(strings.Join(req.path[1:], ""), "/")
		headers := make(map[HttpHeader]string)
		headers[ContentType] = "text/plain"
		headers[ContentLength] = strconv.Itoa(len(message))
		return buildResponse(200, "OK", headers, message)
	})

	createController("/user-agent", func (req HttpRequest) string {
		message := req.headers[UserAgent]
		fmt.Println("DEBUG", req.headers)
		headers := make(map[HttpHeader]string)
		headers[ContentType] = "text/plain"
		headers[ContentLength] = strconv.Itoa(len(message))
		return buildResponse(200, "OK", headers, message)
	})

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
		
		buf := make([]byte, 1024)
	
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading request: ", err.Error())
			os.Exit(1)
		}
	
		hRequest := parseRequest(buf)

		resp := getController(hRequest.path[0])(hRequest)

		fmt.Println("Writing response")
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
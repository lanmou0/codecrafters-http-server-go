package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type HttpHeader string;

const (
	ContentType HttpHeader = "Content-Type"
	ContentLength HttpHeader = "Content-Length" 
	Host HttpHeader = "Host"
	UserAgent HttpHeader = "User-Agent"
	Accept HttpHeader = "Accept"
)

type HttpMethod string;

const (
	GET HttpMethod = "GET"
)

type HttpRequest struct {
	version string;
	path []string;
	method HttpMethod;
	headers map[HttpHeader]string;
	body string;
}

func buildResponse(httpCode int, httpMessage string, headers map[HttpHeader]string, body string) string {	
	resp := fmt.Sprintf("HTTP/1.1 %d %s\r\n", httpCode, httpMessage)

	for k, v := range headers {
		resp += fmt.Sprintf("%s: %s\r\n", string(k), v)
	}

	resp += fmt.Sprintf("\r\n%s\r\n", body)

	return resp
}

func parseRequest(request []byte, id string) HttpRequest {
	fmt.Println(id, "Parsing request")

	hRequest := HttpRequest{};
	lines := strings.Split(string(request), ESCAPE)

	meta := strings.Split(lines[0], " ")
	hRequest.version = meta[0]

	hRequest.path = parsePath(meta[1])
	hRequest.method = HttpMethod(meta[2])

	if(lines[1] == "\r\n") {
		return hRequest
	}

	headersIndex := slices.Index(lines[1:], "")
	hRequest.headers = make(map[HttpHeader]string)
	
	for _, h := range lines[1:headersIndex] {
		h1 := strings.Split(h, ":")
		hRequest.headers[HttpHeader(h1[0])] = strings.TrimSpace(h1[1])
	}

	if(len(lines) <= headersIndex) {
		return hRequest
	}

	hRequest.body = lines[headersIndex+1]
	return hRequest
}

func launchServer(addr string, port string) {

	if(validateAddrAndPort(addr, port) < 0) {
		return
	}

	fullAddr := fmt.Sprintf("%s:%s", addr, port)

	l, err := net.Listen("tcp", fullAddr)
	if err != nil {
		fmt.Println("Failed to bind to port ", port)
		os.Exit(1)
	}
	fmt.Println("listening on Port ", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		
		connId := uuid.NewString()+": "
		fmt.Println(connId, "Connection Accepted")

		go handleRequest(conn, connId)
	}
}

func handleRequest(conn net.Conn, id string) {
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(id, "Error reading request: ", err.Error())
		return
	}

	hRequest := parseRequest(buf, id)

	resp := getController(hRequest.path[0], id)(hRequest)

	fmt.Println(id, "Writing response")
	_, err = conn.Write([]byte(resp))
	if err != nil {
		fmt.Println(id, "Error writing request: ", err.Error())
		return
	}
	fmt.Println(id, "Response sent")

	fmt.Println(id, "Connection Closed\r\n")
	conn.Close()
}
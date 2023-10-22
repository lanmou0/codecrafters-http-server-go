package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
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
	POST HttpMethod = "POST"
)

type HttpRequest struct {
	requestId string
	version string;
	path []string;
	method HttpMethod;
	headers map[HttpHeader]string;
	body []byte;
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
	hRequest.requestId = id

	lines := strings.Split(string(request), ESCAPE)

	for _, l := range lines {
		fmt.Println(id, l)
	}

	meta := strings.Split(lines[0], " ")
	hRequest.version = meta[2]

	hRequest.path = parsePath(meta[1])
	hRequest.method = HttpMethod(meta[0])

	if(lines[1] == "") {
		return hRequest
	}

	headersIndex := slices.Index(lines[1:], "")
	hRequest.headers = make(map[HttpHeader]string)

	for _, h := range lines[1:headersIndex+1] {
		h1 := strings.Split(h, ":")
		hRequest.headers[HttpHeader(h1[0])] = strings.TrimSpace(h1[1])
	}

	if(len(lines) <= headersIndex) {
		return hRequest
	}

	bodyLen := 0

	if(hRequest.headers[ContentLength] == "") {
		for i, l := range lines[headersIndex+2:] {
			if(l == "" && lines[i+1] == "") {
				break
			}
			bodyLen += len(l)
		}
	}else {
		bodyLen, _ = strconv.Atoi(hRequest.headers[ContentLength])
	}

	fmt.Println(id, "BodyLen", bodyLen)

	var buf = make([]byte, bodyLen)

	for _, l := range lines[headersIndex+2:] {
		buf = append(buf, []byte(l)...)
	}

	hRequest.body = bytes.Trim(buf, "\x00")
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

	resp := getController(hRequest.method, hRequest.path[0], id)(hRequest)

	fmt.Println(id, "Writing response")
	_, err = conn.Write([]byte(resp))
	if err != nil {
		fmt.Println(id, "Error writing request:", err.Error())
		return
	}
	fmt.Println(id, "Response sent")

	fmt.Println(id, "Connection Closed\r\n")
	conn.Close()
}
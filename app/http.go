package main

import (
	"fmt"
	"regexp"
	"strings"
	"os"
	"net"

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

func parseRequest(request []byte) HttpRequest {
	fmt.Println("Parsing request: \n", string(request))

	hRequest := HttpRequest{};
	lines := strings.Split(string(request), ESCAPE)

	meta := strings.Split(lines[0], " ")
	hRequest.version = meta[0]

	hRequest.path = parsePath(meta[1])
	hRequest.method = HttpMethod(meta[2])

	if(len(lines) <= 2) {
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

func parsePath(rawPath string) []string {
	re := regexp.MustCompile(`/[^/]*[^/]*`)
	return re.FindAllString(rawPath, -1)
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
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

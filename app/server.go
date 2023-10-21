package main

import (
	"fmt"
	"strconv"
	"strings"
)

const ROOT_PATH = "/"
const ECHO_PATH = "/echo/"
const NIL_PATH = "nil"

var ESCAPE = "\r\n"

func main() {
	fmt.Println("Logs from your program will appear here!")

	createController("nil", func(req HttpRequest) string {
		message := fmt.Sprintf("Path Not Found: %s", req.path)
		headers := make(map[HttpHeader]string)
		headers[ContentType] = "text/plain"
		headers[ContentLength] = strconv.Itoa(len(message))
		return buildResponse(404, "Not Found", headers, message)
	})

	createController("/", func(req HttpRequest) string {
		return buildResponse(200, "OK", make(map[HttpHeader]string), "")
	})

	createController("/echo", func(req HttpRequest) string {
		message := strings.TrimPrefix(strings.Join(req.path[1:], ""), "/")
		headers := make(map[HttpHeader]string)
		headers[ContentType] = "text/plain"
		headers[ContentLength] = strconv.Itoa(len(message))
		return buildResponse(200, "OK", headers, message)
	})

	createController("/user-agent", func(req HttpRequest) string {
		message := req.headers[UserAgent]
		headers := make(map[HttpHeader]string)
		headers[ContentType] = "text/plain"
		headers[ContentLength] = strconv.Itoa(len(message))
		return buildResponse(200, "OK", headers, message)
	})

	launchServer("0.0.0.0", "4221")
}
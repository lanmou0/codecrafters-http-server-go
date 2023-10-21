package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

const ROOT_PATH = "/"
const ECHO_PATH = "/echo/"
const NIL_PATH = "nil"

var ESCAPE = "\r\n"

func main() {
	fmt.Println("Logs from your program will appear here!")

	dirF := flag.String("directory", "", "Provide Directory For Files Path")
	flag.Parse()

	createController("nil", func(req HttpRequest) string {
		return buildErrorResponse(404, "Not Found", fmt.Sprintf("Path Not Found: %s", req.path)) 
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

	createController("/files", func(req HttpRequest) string {
		if(*dirF == "") {
			return buildErrorResponse(500, "Internal Server Error", "Dir Was Not Provided")
		}

		fileName :=  strings.TrimPrefix(req.path[1], "/")
		fileContent, err := os.ReadFile(path.Join(*dirF, fileName))
		if(err != nil) {
			return buildErrorResponse(404, "Not Found", "Error Reading File")
		}

		message := string(fileContent)
		headers := make(map[HttpHeader]string)
		headers[ContentType] = "application/octet-stream"
		headers[ContentLength] = strconv.Itoa(len(message))
		return buildResponse(200, "OK", headers, message)
	})

	launchServer("0.0.0.0", "4221")
}
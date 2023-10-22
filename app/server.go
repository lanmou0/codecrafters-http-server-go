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
const NF_PATH = "GET_404"

var ESCAPE = "\r\n"

func main() {
	fmt.Println("Logs from your program will appear here!")

	dirF := flag.String("directory", "", "Provide Directory For Files Path")
	flag.Parse()

	createController(GET ,"404", func(req HttpRequest) string {
		fmt.Printf(req.requestId, "Path Not Found:", req.method, req.path)
		return buildErrorResponse(404, "Not Found", fmt.Sprintf("Path Not Found: %s", req.path[0])) 
	})

	createController(GET, "/", func(req HttpRequest) string {
		return buildResponse(200, "OK", make(map[HttpHeader]string), "")
	})

	createController(GET, "/echo", func(req HttpRequest) string {
		message := strings.TrimPrefix(strings.Join(req.path[1:], ""), "/")
		headers := make(map[HttpHeader]string)
		headers[ContentType] = "text/plain"
		headers[ContentLength] = strconv.Itoa(len(message))
		return buildResponse(200, "OK", headers, message)
	})

	createController(GET, "/user-agent", func(req HttpRequest) string {
		message := req.headers[UserAgent]
		headers := make(map[HttpHeader]string)
		headers[ContentType] = "text/plain"
		headers[ContentLength] = strconv.Itoa(len(message))
		return buildResponse(200, "OK", headers, message)
	})

	createController(GET, "/files", func(req HttpRequest) string {
		if(*dirF == "") {
			return buildErrorResponse(500, "Internal Server Error", "Dir Was Not Provided")
		}

		fileName :=  strings.TrimPrefix(req.path[1], "/")
		fileContent, err := os.ReadFile(path.Join(*dirF, fileName))
		if(err != nil) {
			fmt.Println(req.requestId, "Error reading file:", err.Error())
			return buildErrorResponse(404, "Not Found", fmt.Sprintf("Can't read File: %s", fileName))
		}

		message := string(fileContent)
		headers := make(map[HttpHeader]string)
		headers[ContentType] = "application/octet-stream"
		headers[ContentLength] = strconv.Itoa(len(message))
		return buildResponse(200, "OK", headers, message)
	})

	createController(POST, "/files", func(req HttpRequest) string {
		if(*dirF == "") {
			return buildErrorResponse(500, "Internal Server Error", "Dir Was Not Provided")
		}

		fileName :=  strings.TrimPrefix(req.path[1], "/")
		err := os.WriteFile(path.Join(*dirF, fileName), req.body, 0644)
		if(err != nil) {
			fmt.Println(req.requestId, "Error writing file:", err.Error())
			return buildErrorResponse(404, "Not Found", "Error Writing File")
		}

		headers := make(map[HttpHeader]string)
		return buildResponse(201, "Created", headers, "")
	})

	launchServer("0.0.0.0", "4221")
}
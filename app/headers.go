package main

import (
	"fmt"
)

type header string;

const (
	contentType header = "Content-Type"
	contentSize header = "Content-Size" 
)

func buildResponse(httpCode int, httpMessage string, headers map[header]string, body string) string {	
	resp := fmt.Sprintf("HTTP/1.1 %d %s\r\n", httpCode, httpMessage)

	for k, v := range headers {
		resp += fmt.Sprintf("%s: %s\r\n", string(k), v)
	}

	resp += fmt.Sprintf("\r\n%s\r\n", body)

	return resp
}
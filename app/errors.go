package main

import (
	"strconv"
)

func buildErrorResponse(httpCode int, httpMessage string, responseMessage string) string {
	headers := make(map[HttpHeader]string)
	headers[ContentType] = "text/plain"
	headers[ContentLength] = strconv.Itoa(len(responseMessage))
	return buildResponse(httpCode, httpMessage, headers, responseMessage)
}
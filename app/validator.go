package main

import "fmt"

func validateAddrAndPort(addr string, port string) int {
	if(len(addr)*len(port) == 0) {
		fmt.Println("Address or Port were provided")
		return -1
	}
	//Add More Validation

	return 1
}

package main

import (
	"fmt"
	"os"
)

var controllers = make(map[string]func(HttpRequest) string);

func createController(method HttpMethod, path string, controller func(HttpRequest) string) {
	fmt.Println("Adding controller for ", method, path)
	key := getKey(method, path)
	if(controllers[key] != nil) {
		fmt.Println("Controller already created for ", method, path)
		os.Exit(0)
	}
	controllers[key] = controller
}

func getController(method HttpMethod, path string, id string) func(HttpRequest) string {
	key := getKey(method, path)
	if(controllers[key] == nil) {
		fmt.Println(id, "Path not found ", method, path)
		return controllers[NIL_PATH]
	}
	fmt.Println(id, "Handling request for ", method, path)
	return controllers[key]
}

func getKey(method HttpMethod, path string) string {
	return fmt.Sprint("%s_%s", path, string(method))
}


package main

import (
	"fmt"
)

var controllers = make(map[string]func(HttpRequest) string);

func createController(path string, controller func(HttpRequest) string) {
	fmt.Println("Adding controller for path: ", path)
	controllers[path] = controller
}

func getController(path string, id string) func(HttpRequest) string {
	if(controllers[path] == nil) {
		fmt.Println(id, "Path not found ", path)
		return controllers[NIL_PATH]
	}
	fmt.Println(id, "Handling request for path ", path)
	return controllers[path]
}


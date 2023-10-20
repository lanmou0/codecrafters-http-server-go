package main

import (
	"fmt"
)

var controllers = make(map[string]func(HttpRequest) string);

func createController(path string, controller func(HttpRequest) string) {
	fmt.Println("Adding controller for path: ", path)
	controllers[path] = controller
}

func getController(path string) func(HttpRequest) string {
	if(controllers[path] == nil) {
		return controllers[NIL_PATH]
	}
	return controllers[path]
}


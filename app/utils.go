package main

import "regexp"

func parsePath(rawPath string) []string {
	re := regexp.MustCompile(`/[^/]*[^/]*`)
	return re.FindAllString(rawPath, -1)
}
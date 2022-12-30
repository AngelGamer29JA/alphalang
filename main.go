package main

import (
	Alpha "alpha/alpha"
	"alpha/alpha/io"
	"alpha/alpha/std/color"
	"fmt"
	"os"
	"strings"
)

var AcceptedSuffix []string = []string{
	".sundi",
	".sdi",
	".alpha",
	".alp",
}

func IsCorrectFileExtension(value string) bool {
	for _, suffix := range AcceptedSuffix {
		if strings.HasSuffix(value, suffix) {
			return true
		}
	}
	return false
}

func main() {
	// fmt.Println(Alpha.VARIABLE_DEFINE_PATTERN)
	args_len := len(io.Args)
	if args_len == 0 {
		fmt.Println("No arguments presented")
		os.Exit(0)
	}

	path := io.Args[0]

	if !io.FileExists(path) {
		color.Printf("&4error&r: file not found '%s'\n", path)
		os.Exit(-1)
	}

	content := io.Readfile(path)

	if !IsCorrectFileExtension(path) {
		accepted := strings.Join(AcceptedSuffix, ", ")
		color.Printf("&4error&r: invalid file extension, accepted extensions %s\n", accepted)
		os.Exit(-1)
	}

	instance := Alpha.NewInstance(content, path)
	instance.Execute()
}

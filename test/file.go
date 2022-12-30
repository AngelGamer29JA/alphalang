package main

import (
	"alpha/alpha/io"
	"alpha/alpha/std/color"
)

func main() {
	io.WriteLineAt("data.txt", color.Sprint("hola"), 1)
}

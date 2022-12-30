package color

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"golang.org/x/term"
)

var ansiColors = map[string]string{
	"&0": "\u001b[30m", // Black
	"&1": "\u001b[34m", // Dark Blue
	"&2": "\u001b[32m", // Dark Green
	"&3": "\u001b[36m", // Dark Aqua
	"&4": "\u001b[31m", // Dark Red
	"&5": "\u001b[35m", // Dark Purple
	"&6": "\u001b[33m", // Gold
	"&7": "\u001b[37m", // Gray
	"&8": "\u001b[90m", // Dark Gray
	"&9": "\u001b[94m", // Blue
	"&a": "\u001b[92m", // Green
	"&b": "\u001b[96m", // Aqua
	"&c": "\u001b[91m", // Red
	"&d": "\u001b[95m", // Light Purple
	"&e": "\u001b[93m", // Yellow
	"&f": "\u001b[97m", // White
	"&r": "\u001b[0m",  // Reset
	"&n": "\u001b[4m",  // Underlined
	"&v": "\u001b[7m",  // Reverse
	"&l": "\u001b[1m",  // Bold
	"&i": "\u001b[3m",  // Italics
	"&m": "\u001b[9m",  // Strikethrough
}

func Printf(format string, a ...any) {
	// Create a new Colorable instance
	clb := colorable.NewColorableStdout()

	value := Sprintf(format, a...)
	runes := []rune(value)
	data := []byte(string(runes))
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Fprint(clb, value)
	} else {
		clb.Write(data)
	}
}

func Print(value ...any) {
	// Create a new Colorable instance
	clb := colorable.NewColorableStdout()

	data := []byte(Sprint(value...))
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Fprint(clb, string(data))
	} else {
		clb.Write(data)
	}
}

func Println(value ...any) {
	// Create a new Colorable instance
	clb := colorable.NewColorableStdout()

	data := []byte(Sprintln(value...))
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Fprint(clb, string(data))
	} else {
		clb.Write(data)
	}
}

func Sprintf(input string, a ...any) string {
	// Add reset at end of the line
	input = input + "&r"
	input = fmt.Sprintf(input, a...)
	// Replace Minecraft color codes with corresponding ANSI codes
	for mc, ansi := range ansiColors {
		input = strings.ReplaceAll(input, mc, ansi)
	}
	return input
}

func Sprintln(input ...any) string {

	data := ""
	for _, elm := range input {
		data += fmt.Sprint(elm)
	}

	data = data + "&r"
	// Replace Minecraft color codes with corresponding ANSI codes
	for mc, ansi := range ansiColors {
		data = strings.ReplaceAll(data, mc, ansi)
	}

	return fmt.Sprintln(data)
}

func Sprint(input ...any) string {

	data := ""
	for _, elm := range input {
		data += fmt.Sprint(elm)
	}

	data = data + "&r"
	// Replace Minecraft color codes with corresponding ANSI codes
	for mc, ansi := range ansiColors {
		data = strings.ReplaceAll(data, mc, ansi)
	}

	return fmt.Sprint(data)
}

func Errorf(input string, a ...any) error {
	return errors.New(Sprintf(input, a...))
}

// func PrintColors(input string) {
// 	// Create a new Colorable instance
// 	colorable := colorable.NewColorableStdout()

// 	// Add reset at end of the line
// 	input = input + "&r"

// 	// Replace Minecraft color codes with corresponding ANSI codes
// 	for mc, ansi := range ansiColors {
// 		input = strings.ReplaceAll(input, mc, ansi)
// 	}

// 	// Return formatted string using the Colorable instance
// 	colorable.Write([]byte(input))
// }

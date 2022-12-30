package io

import (
	"alpha/alpha/std/color"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var Args = os.Args[1:]

/*
Writes a line or modifies a line at a specific position in the file.
*/
func WriteLineAt(path, content string, lineIdx int) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	color.Expect(err)
	defer file.Close()

	file.WriteAt([]byte(content), int64(lineIdx))
}

func WriteLine(path, content string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			color.Expect(color.Errorf("file not exists, %s", path))
		}
		color.Expect(err)
	}
	defer file.Close()

	_, err = file.Seek(0, 2)
	color.Expect(err)

	fmt.Fprintf(file, "%s\n", content)

	err = file.Sync()
	color.Expect(err)
}

func Write(path, content string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			color.Printf("&4error&r: '%s' file not exist\n", path)
			return
		}
		color.Printf("&4error&r: %s\n", err)
		return
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		color.Printf("&4error&r: error to write in file %s\n", path)
	}

	file.Close()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Readfile(path string) string {
	// Leer el contenido del archivo
	data, err := os.ReadFile(path)
	if err != nil {
		// Comprobar si el error se debió a que el archivo no existe
		if os.IsNotExist(err) {
			color.Printf("&4error: &rfile not exist %s\n", path)
			return ""
		}
		// Si se produjo otro tipo de error, retornarlo
		fmt.Println(err)
		return ""
	}

	// Si se pudo leer el contenido del archivo, retornarlo como una cadena
	return string(data)
}

func ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	var lines []string
	if err != nil {
		if err == os.ErrExist {
			return []string{}, color.Errorf("&4error&r: error to read file lines, file not found %s", path)
		}
		return []string{}, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	file.Close()
	return lines, nil
}

func CreateFile(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		if err == os.ErrExist {
			file.Close()
			return nil
		}
		return color.Errorf("&4error: &rfile not created.\n\t&4err: &r%s\n", err.Error())
	}
	defer file.Close()
	file.Write([]byte(content))
	return nil
}

func GetAbsolutePath(abs string) string {
	p, err := filepath.Abs(abs)
	if err != nil {
		color.Println("&4error: &rinternal error has occured ", err)
		return ""
	}
	return p
}

// Get current root directory os.Getwd()
func GetCurrentWd() string {
	root, err := os.Getwd()
	if err != nil {
		color.Println("&4error: &rerror to get current working directory.")
	}
	return root
}

func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	if err != nil {
		if err == io.EOF {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err)
	}

	return strings.TrimSpace(input)
}

func Mkdir(path string) error {
	err := os.Mkdir(path, os.FileMode(0700))
	if err != nil {
		if err == os.ErrExist {
			return color.Errorf("&4error&r: folder already exists.")
		}

		return err
	}
	return nil
}

func Delete(path string) error {
	err := os.Remove(path)
	if err != nil {
		if err == os.ErrExist {
			return color.Errorf("&4error&r: file or folder not exists")
		}
		return err
	}
	return nil
}

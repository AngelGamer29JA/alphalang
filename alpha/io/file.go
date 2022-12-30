package io

import (
	"os"
)

type File struct {
	File    os.File
	Content string
	Lines   []string
}

func (f *File) WriteLine(content string) {
	WriteLine(f.File.Name(), content)
}

func (f *File) Write(content string) {
	Write(f.File.Name(), content)
}

func (f *File) Delete() {
	Delete(f.File.Name())
}

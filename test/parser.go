package main

import (
	"alpha/alpha"
	io "alpha/alpha/io"
	"fmt"
	"strings"
)

var w int = 1

func ReadStatement(stm *alpha.BlockStatement) {
	w += 1
	for _, stt := range stm.Scope {
		fmt.Println(strings.Repeat("\t", w), stt.Value)
		if stt.IsStatement {
			ReadStatement(&stt.Block)
		}
	}
}

func main() {

	input := io.Readfile("samples/loop.alpha")
	parser := alpha.NewParser("main.alpha", input)

	parser.Parse()

	if len(parser.Errors) > 0 {
		parser.ShowTraceback()
		return
	}

	parsed := parser.CodeParse
	for _, p := range parsed {
		fmt.Println(p.Value)
		if p.IsStatement {
			for _, st := range p.Block.Scope {
				fmt.Println("\t", st.Value)
				if st.IsStatement {
					ReadStatement(&st.Block)
				}
			}
		}
	}
}

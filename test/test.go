package main

import (
	"fmt"

	"github.com/antonmedv/expr"
)

func Eval() {

}

func main() {
	env := new(interface{})
	code := `5+3,5`

	program, err := expr.Compile(code, expr.Optimize(true))

	if err != nil {
		fmt.Println(err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Println(err)

		return
	}
	fmt.Println(output)

}

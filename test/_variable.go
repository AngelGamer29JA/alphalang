package main

import (
	Alpha "alpha/alpha"
	"fmt"
)

func main() {
	input := `name = n`
	variable, err := Alpha.GetVariableDefinition(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(variable.Name, variable.Content, variable.Type)
}

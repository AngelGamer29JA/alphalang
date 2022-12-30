package main

import (
	Alpha "alpha/alpha"
	"fmt"

	"github.com/dlclark/regexp2"
)

func ReplaceIndex(str string, char rune, index int) string {
	bites := make([]rune, len(str))
	for i := 0; i < len(str); i++ {
		if i == index {
			bites[i] = char
		} else {
			bites[i] = rune(str[i])
		}
	}
	return string(bites)
}

var expressions = []string{
	"sendln %any%",
}

func CheckExpression(input string) bool {
	for _, expression := range expressions {
		template := Alpha.ParseTemplate(expression)
		fmt.Println(template)
		regex := regexp2.MustCompile(template, regexp2.Compiled)
		m, _ := regex.MatchString(input)
		if m {
			return true
		}
	}
	return false
}

func main() {
	temp := CheckExpression("sendln 5")
	fmt.Println(temp)
}

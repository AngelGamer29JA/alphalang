package alpha

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
)

type Expression struct {
	Value        string
	Executor     func(code *Code, instance *Instance)
	Check        bool
	Assignations bool
}

/*
value = "my instruction %string%"

executor = myInstruction(code *Code, instance *Instance)

check = Check variables passed example, my instruction {variable}, if true throw error when variable not found

assigns = Check if instruction change any variable, if false don't check, if true throw error "variable not found"


*/
func NewExpression(value string, executor func(code *Code, instance *Instance), check bool, assigns bool) *Expression {
	return &Expression{Value: value, Executor: executor, Check: check, Assignations: assigns}
}

func CheckExpression(input string) bool {
	if input == "end" || input == "stop" {
		return true
	}
	for _, expression := range BuiltinExpressions {
		template := ParseTemplate(expression.Value)
		regex := regexp2.MustCompile(template, regexp2.Compiled)
		m, _ := regex.MatchString(input)
		if m {
			return true
		}
	}
	return false
}

func GetExpression(value string) (Expression, error) {
	for _, expression := range BuiltinExpressions {
		template := ParseTemplate(expression.Value)
		regex := regexp2.MustCompile(template, regexp2.Compiled)
		m, _ := regex.MatchString(value)
		if m {
			return expression, nil
		}
	}
	return Expression{}, fmt.Errorf("expression not found: %s", value)
}

func RemoveComments(str string) string {
	regex := regexp.MustCompile(`#.+`)
	_str := regex.ReplaceAllString(str, "")
	return strings.TrimSpace(_str)
}

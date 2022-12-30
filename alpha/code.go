package alpha

import (
	"fmt"
	"regexp"
	"strings"
)

type Code struct {
	Value       string
	Expr        Expression
	Assignation VariableDefinition
	Statement Block
	Conditional ConditionalBlock
	Kind      ParserKind
	Line      int
}

type CodeLine struct {
	Value       string
	Line        int
	IsStatement bool
	Block       BlockStatement
	Kind        ParserKind
}

type CodeParser struct {
	Value       string
	Line        int
	IsStatement bool
	Block       BlockStatement
	Kind        ParserKind
}

func NewCodeParser(value string, line int, kind ParserKind, statement bool, block BlockStatement) *CodeParser {
	return &CodeParser{
		Value:       value,
		Line:        line,
		Kind:        kind,
		IsStatement: statement,
		Block:       block,
	}
}

func CreateCode(value string, expr Expression, assign VariableDefinition, block Block, line int, kind ParserKind) *Code {
	return &Code{
		Value:       value,
		Expr:        expr,
		Assignation: assign,
		Statement:   block,
		Line:        line,
		Kind:        kind,
	}
}

func NewCode(value string, expr Expression, assign VariableDefinition, isdefinition bool, line int) *Code {
	return &Code{
		Value:       value,
		Expr:        expr,
		Assignation: assign,
		// IsVariableDefinition: isdefinition,
		Line: line,
	}
}

func (c *Code) Error(name, message string, a ...any) {
	if len(a) > 0 {
		Error(name, c.Line, fmt.Sprintf(message, a...))
	} else {
		Error(name, c.Line, message)
	}
}

func (c *Code) GetVariableInstances() []VariableInstance {
	return GetVariableInstances(c.Value)
}

func (c *Code) GetContents() []string {
	arr := []string{}
	regex, err := regexp.Compile(STR_NUM_BOOL_NULL)
	if err != nil {
		fmt.Println("Err get contents: ", err)
		return []string{}
	}

	value := c.Value
	content := regex.FindAllString(value, -1)

	for _, cap := range content {
		arr = append(arr, strings.ReplaceAll(cap, `"`, ""))
	}
	return arr
}

func (c *Code) GetContentAt(index int) (string, error) {
	contents := c.GetContents()

	if len(contents) <= 0 {
		return "", fmt.Errorf("out of range index content %d", index)
	}
	return contents[index], nil

}

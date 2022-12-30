package alpha

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
)

/*
La estructura Block es la encargada de guardar los bloques de codigo (statement),
y la funcion GetVariableInstances() se encarga de encontrar variables llamadas en una
linea de codigo, un ejemplo sería 'sendln "{myvariable}"' o 'make dir "" stored in {}'.

GetContents() es muy similar a la anterior funcion ya mencionada pero esta busca cualquier
tipo de dato existente en el lenguaje como de tipo number, string, bool, null o un array, un ejemplo
de como los analiza y encuentra sería:

	sendln "este es un contenido" o 5 o true y false o 5, 10, true, ""
	        --------------------
			^
			|
		  string

	sendln 5
	sendln true
	sendln false
	sendln "string"
	sendln 10, 5, true, "string"

Una simple funcion que busca un contenido en el indice especifico, de no existir el indice o que
el array de contenido este vacío retornara un error indicando que no existe o no hay elementos.

Scope - como lo dice es una estructura utilizada para almacenar el codigo que en la funcion Parse() ya
ha convertido a los ejecutores de codigo, bloques de codigo y declaraciones de variable, en esta estructura
se almacenan como antes dicho, variables que utilizan la estructura Variable en un map[string]Variable,
la propiedad Code []Code es un array que contiene las instrucciones para ser ejecutadas en una instancia.

BlockStatement es una estructura utilizada para
*/
type Block struct {
	Start    int
	End      int
	Scope    Scope
	Executor func(*Block, *Instance)
	Value    string
}

const (
	BlockStartPrefix = "do"
	BlockEndPrefix   = "end"
)

func (b *Block) GetVariableInstances() []VariableInstance {
	return GetVariableInstances(b.Value)
}

func (b *Block) GetContents() []string {
	arr := []string{}
	regex, err := regexp.Compile(STR_NUM_BOOL_NULL)
	if err != nil {
		fmt.Println("Err get contents: ", err)
		return []string{}
	}

	value := b.Value
	content := regex.FindAllString(value, -1)

	for _, cap := range content {
		arr = append(arr, strings.ReplaceAll(cap, `"`, ""))
	}
	return arr
}

func (b *Block) GetContentAt(index int) (string, error) {
	contents := b.GetContents()

	if len(contents) <= 0 {
		return "", fmt.Errorf("out of range index content %d", index)
	}
	return contents[index], nil

}

type Scope struct {
	Code       []Code
	Variables  map[string]*Variable
	Statements []Block
}

func (sf *Block) AddStatement(value string, statement Block, line int) {
	sf.Scope.Code = append(sf.Scope.Code, *CreateCode(value, Expression{}, VariableDefinition{}, statement, line, STATEMENT_EXPRESSION))
}

type BlockStatement struct {
	Start int
	End   int
	Scope []CodeParser
	Value string
}

func NewBlockStatement(value string) *BlockStatement {
	return &BlockStatement{Value: value, Start: 0, End: 0, Scope: make([]CodeParser, 0)}
}

func (b *BlockStatement) AddCode(value string, line int, kind ParserKind) {
	b.Scope = append(b.Scope, *NewCodeParser(value, line, kind, false, BlockStatement{}))
}

func (b *BlockStatement) AddStatement(value string, line int, statement BlockStatement) {
	b.Scope = append(b.Scope, *NewCodeParser(value, line, STATEMENT_EXPRESSION, true, statement))
}

func NewBlock(value string, executor func(*Block, *Instance)) *Block {
	return &Block{Executor: executor, Value: value}
}

func CheckStatement(input string) bool {
	input = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(input), "do"))
	for _, statement := range BuiltinStatements {
		template := ParseTemplate(statement.Value)
		regex := regexp2.MustCompile(template, regexp2.Compiled)
		m, _ := regex.MatchString(input)
		if m {
			return true
		}
	}
	return false
}

func GetStatement(value string) (Block, error) {
	value = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(value), "do"))
	for _, statement := range BuiltinStatements {
		template := ParseTemplate(statement.Value)
		regex := regexp2.MustCompile(template, regexp2.Compiled)
		m, _ := regex.MatchString(value)
		if m {
			return statement, nil
		}
	}
	return Block{}, fmt.Errorf("statement not found: %s", value)
}

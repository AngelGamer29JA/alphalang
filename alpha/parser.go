package alpha

import (
	"alpha/alpha/std/color"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	VARIABLE_DEFINITION  ParserKind = 0
	CODE_EXPRESSION      ParserKind = 1
	STATEMENT_EXPRESSION ParserKind = 2
	CONDITION_STATEMENT  ParserKind = 3
)

type ParseValue struct {
	CodeParser
	Definition     VariableDefinition
	CodeExpression Code
}

type ParserKind int32

type ParserError struct {
	Line  int
	Value string
	Name  string
}

type Parser struct {
	Line       int
	Data       []string
	Name       string
	Errors     []ParserError
	CodeParse  []CodeParser
	blockDepth int
	//Scanning bool
}

func NewParser(name string, data string) *Parser {
	name = filepath.Base(name)
	__data := strings.Split(data, "\n")
	return &Parser{Name: name, Data: __data}
}

func ScanParser(name string, data []string) *Parser {
	return &Parser{Name: name, Data: data}
}

func (p *Parser) Compile(instance *Instance) {
	p.Parse()
	for _, code := range p.CodeParse {
		if code.Kind == CODE_EXPRESSION {
			expr, _ := GetExpression(code.Value)
			instance.AddCode(code.Value,
				*NewExpression(
					expr.Value,
					expr.Executor,
					expr.Check,
					expr.Assignations,
				),
				VariableDefinition{},
				code.Line,
				CODE_EXPRESSION,
			)
		} else if code.Kind == VARIABLE_DEFINITION {
			varDef, _ := GetVariableDefinition(code.Value)
			instance.AddCode(code.Value,
				Expression{},
				*NewVariableDefinition(
					varDef.Name,
					varDef.Content,
					varDef.Type,
				),
				code.Line,
				VARIABLE_DEFINITION,
			)
		} else if code.Kind == STATEMENT_EXPRESSION {
			blockStatement := p.CompileStatement(code.Block)
			instance.AddStatement(code.Value, blockStatement, code.Line)
		}
	}
}

func (p *Parser) CompileStatement(b BlockStatement) Block {
	bb, _ := GetStatement(b.Value)
	blockStatement := NewBlock(b.Value, bb.Executor)
	statement := b

	blockStatement.Start = statement.Start
	blockStatement.End = statement.End
	for _, nested := range b.Scope {
		if nested.Kind == VARIABLE_DEFINITION {
			variable, _ := GetVariableDefinition(nested.Value)
			varDef := *NewVariableDefinition(
				variable.Name,
				variable.Content,
				variable.Type,
			)
			blockStatement.Scope.Code = append(blockStatement.Scope.Code, *CreateCode(nested.Value, Expression{}, varDef, Block{}, nested.Line, nested.Kind))
		}

		if nested.Kind == STATEMENT_EXPRESSION {
			blockk := p.CompileStatement(nested.Block)
			blockStatement.AddStatement(nested.Value, blockk, nested.Line)
		}

		if nested.Kind == CODE_EXPRESSION {
			expr, _ := GetExpression(nested.Value)
			blockStatement.Scope.Code = append(blockStatement.Scope.Code, *CreateCode(nested.Value, expr, VariableDefinition{}, Block{}, nested.Line, nested.Kind))
		}
	}

	return *blockStatement
}

func (p *Parser) Analyze() {
	maxLines := len(p.Data)
	p.blockDepth = 0
	for i := 0; i < maxLines; i++ {
		line := RemoveComments(p.Data[i])
		if !isEmpty(line) {
			err := p.CheckLine(line, i)
			if err != nil {
				p.AddError(err.Error(), i+1)
			}
		}
	}
}

func (p *Parser) CheckLine(line string, idx int) error {
	regex := regexp.MustCompile(VARIABLE_INSTANCE_PATTERN)
	match := regex.MatchString(line)

	if IsVariableDefinition(line) && GetVariableType(line) == "invalid type" && !CheckStatement(line) {
		if !match {
			// Si la línea no es una definición de variable válida, agregamos un error a la lista de errores
			//p.AddError(fmt.Sprintf("invalid variable definition %s", line), i)
			return fmt.Errorf("invalid variable definition %s", line)
		}

	}

	if !IsVariableDefinition(line) && !CheckExpression(line) {
		if !strings.HasSuffix(line, BlockStartPrefix) {
			//p.AddError(fmt.Sprintf("invalid expression syntax %s", line), i) // Expression
			return fmt.Errorf("invalid expression syntax %s", line)
		} else if strings.HasSuffix(line, BlockStartPrefix) {
			if !CheckStatement(line) {
				return fmt.Errorf("invalid statement header %s", line)
			} else if CheckStatement(line) {
				p.blockDepth++
				p.AnalyzeStatement(idx)
			}
		}
	}

	// Si se encuentra un fin de bloque, se disminuye en uno la profundidad del bloque
	if strings.TrimSpace(line) == BlockEndPrefix {
		// Si se encuentra un fin de bloque sin un inicio de bloque previo, se produce un error
		if p.blockDepth == 0 {
			return fmt.Errorf("unexpected end of block")
		}
		p.blockDepth--
	}
	return nil
}

/*
Esta funcion analiza un bloque de codigo (statement) pasando como parametro la linea
donde se encuentra el bloque de codigo encontrado para analizarlo desde esa linea hasta donde ya no encuentre
un fin de bloque 'end'
*/
func (p *Parser) AnalyzeStatement(start int) { // function: Analyze Statement
	//subList := p.Data[i:]
	maxLines := len(p.Data) - 1
	subLine := start + 1
	blockDepth := 1

	for blockDepth > 0 {
		if subLine > maxLines {
			// Si se ha llegado al final de la lista y todavía hay bloques abiertos, se producirá un error
			startLine := strings.TrimSpace(p.Data[start])
			endLine := "eof"
			closureMessage := fmt.Sprintf("invalid statement closure\n\n\t&6%d >> &r%s\n\t...\n\t&6%d >> &rexpected 'end' but '%s'", start+1, startLine, subLine+1, endLine)
			p.AddError(closureMessage, subLine+1)
			break
		}

		line := strings.TrimSpace(p.Data[subLine])
		if strings.TrimSpace(line) == BlockEndPrefix {
			// Si se encuentra un fin de bloque, se reduce en uno la profundidad del bloque
			blockDepth -= 1
		} else if strings.HasSuffix(line, BlockStartPrefix) {
			// Si se encuentra un nuevo bloque, se aumenta en uno la profundidad del bloque
			blockDepth += 1
		}

		subLine += 1
	}
}

/*func (p *Parser) AnalyzeStatement(start int) { // function: AnalyzeStatement
	//subList := p.Data[i:]
	maxLines := len(p.Data) - 1

	// Recorre todas las líneas de código desde la línea de inicio del bloque hasta la última línea de código
	for i := start + 1; i <= maxLines; i++ {
		// Verifica si existe una palabra clave "end" en la línea actual
		if strings.TrimSpace(p.Data[i]) == "end" {
			// Si se encuentra una palabra clave "end", finaliza la función
			return
		}
	}

	// Si no se encuentra una palabra clave "end", muestra un mensaje de error
	startLine := p.Data[start]
	endLine := strings.TrimSpace(p.Data[maxLines])
	if endLine == "" {
		endLine = "eof"
	}
	closureMessage := fmt.Sprintf("invalid statement closure\n\n\t&6%d >> &r%s\n\t...\n\t&6%d >> &rexpected 'end' but '%s'", start, startLine, maxLines, endLine)
	p.AddError(closureMessage, maxLines+1)

	// subLine := start + 1

	// for strings.TrimSpace(p.Data[subLine]) != "end" {
	// 	if subLine == maxLines {
	// 		if strings.TrimSpace(p.Data[subLine]) != "end" {
	// 			startLine := p.Data[start]
	// 			endLine := strings.TrimSpace(p.Data[subLine])
	// 			if endLine == "" {
	// 				endLine = "eof"
	// 			}
	// 			closureMessage := fmt.Sprintf("invalid statement closure\n\n\t&6%d >> &r%s\n\t...\n\t&6%d >> &rexpected 'end' but '%s'", start, startLine, subLine, endLine)
	// 			p.AddError(closureMessage, subLine+1)
	// 			break
	// 		}
	// 	}

	// 	subLine += 1
	// }
}*/

/*

	Esta funcion se encarga de convertir los datos de la funcion
	Analyze() a algo similar a un arbol de sintaxis para manejar las lineas con codigo
	ya no es necesario proecuparse si hay erroes en esta funcion ya que la funcion que se menciono
	es la encargada de encontrar errores, sin embargo solo encuentra errores de sintaxis y parser.

*/

func (p *Parser) Parse() {
	/*
		`p.Data` es la variable que almacea las lineas de un archivo, cada linea es un string y lo que hago
		es recorrer toda las lineas con un for, y en la primera linea del for esta la variable line que almacena
		la linea actua que a su vez se almacena el numero de linea en la propiedad Line del stuct Parser
	*/

	p.Analyze()

	maxLines := len(p.Data)
	if len(p.Errors) > 0 {
		return
	}

	for i := 0; i < maxLines; i++ {
		if isEmpty(p.CurrentLine()) {
			p.Line += 1
		} else {
			if p.Line != maxLines {
				// codeLine, check := p.CheckCode()
				if IsVariableDefinition(p.CurrentLine()) && !CheckExpression(p.CurrentLine()) {
					p.AddCode(p.CurrentLine(), p.Line, VARIABLE_DEFINITION)
					p.Line += 1
					// return CodeLine{Value: p.CurrentLine(), Line: p.Line, Kind: VARIABLE_DEFINITION}, true
				}

				if CheckStatement(p.CurrentLine()) &&
					strings.HasSuffix(p.CurrentLine(), BlockStartPrefix) &&
					!CheckExpression(p.CurrentLine()) &&
					!IsVariableDefinition(p.CurrentLine()) {
					statement := NewBlockStatement(p.CurrentLine())

					statement.Start = p.Line

					p.Line += 1

					p.Line = p.ReadStatement(statement, p.Line)
					statement.End = p.Line

					p.AddStatement(p.Data[statement.Start], statement.Start, *statement)
					if statement.End != len(p.Data)-1 {
						p.Line += 1
					}
				}

				if CheckExpression(p.CurrentLine()) {
					if p.CurrentLine() != BlockEndPrefix {
						p.AddCode(p.CurrentLine(), p.Line, CODE_EXPRESSION)
						p.Line += 1
					}
				}
			}
		}
	}
}

func (p *Parser) CheckCode() (CodeLine, bool) {
	if IsVariableDefinition(p.CurrentLine()) && !CheckExpression(p.CurrentLine()) {
		//p.AddCode(p.CurrentLine(), p.Line, VARIABLE_DEFINITION)
		return CodeLine{Value: p.CurrentLine(), Line: p.Line, Kind: VARIABLE_DEFINITION}, true
	}

	if CheckStatement(p.CurrentLine()) &&
		strings.HasSuffix(p.CurrentLine(), BlockStartPrefix) &&
		!CheckExpression(p.CurrentLine()) &&
		!IsVariableDefinition(p.CurrentLine()) {
		fmt.Println(p.CurrentLine())
		statement := NewBlockStatement(p.CurrentLine())

		statement.Start = p.Line

		p.Line += 1

		p.Line = p.ReadStatement(statement, p.Line)
		statement.End = p.Line

		//p.AddStatement(p.Data[statement.Start], statement.Start, *statement)
		return CodeLine{Value: p.Data[statement.Start], Line: statement.Start, Kind: STATEMENT_EXPRESSION, Block: *statement, IsStatement: true}, true
	}

	if CheckExpression(p.CurrentLine()) {
		if p.CurrentLine() != BlockEndPrefix || p.CurrentLine() != "stop" {
			//p.AddCode(p.CurrentLine(), p.Line, CODE_EXPRESSION)
			return CodeLine{Value: p.CurrentLine(), Line: p.Line, Kind: CODE_EXPRESSION, IsStatement: false}, true
		}
	}
	return CodeLine{}, false
}

func (p *Parser) ReadStatement(statement *BlockStatement, idx int) int {
	lineIdx := idx
	for strings.TrimSpace(p.Data[lineIdx]) != BlockEndPrefix {
		line := strings.TrimSpace(p.Data[lineIdx])
		if !isEmpty(line) {
			if IsVariableDefinition(line) && !CheckExpression(line) {
				statement.AddCode(line, lineIdx, VARIABLE_DEFINITION)
			}

			if CheckStatement(line) &&
				strings.HasSuffix(line, BlockStartPrefix) {
				nestedStatement := NewBlockStatement(line)

				nestedStatement.Start = lineIdx

				lineIdx += 1

				lineIdx = p.ReadStatement(nestedStatement, lineIdx)
				nestedStatement.End = lineIdx
				statement.AddStatement(line, lineIdx, *nestedStatement)
			}

			if CheckExpression(line) {
				statement.AddCode(line, lineIdx, CODE_EXPRESSION)
			}
		}

		lineIdx += 1
	}
	return lineIdx
}

func (p *Parser) AddCode(value string, line int, kind ParserKind) {
	p.CodeParse = append(p.CodeParse, *NewCodeParser(value, line, kind, false, BlockStatement{}))
}

func (p *Parser) AddStatement(value string, line int, statement BlockStatement) {
	p.CodeParse = append(p.CodeParse, *NewCodeParser(value, line, STATEMENT_EXPRESSION, true, statement))
}

func (p *Parser) ShowTraceback() {
	for _, err := range p.Errors {
		p.Error(color.Sprintf(err.Value), err.Line)
	}
}

func (p *Parser) Error(message string, line int) {
	messageFormatted := color.Sprintf("\n\t%s\n\t&4at &r%s", message, p.GetPath())
	Error(p.Name, line, messageFormatted)
}

func (p *Parser) AddError(err string, line int) {
	p.Errors = append(p.Errors, ParserError{Value: err, Line: line, Name: p.Name})
}

/*
Envar un error formateado a la consola
*/
func Error(name string, line int, message string) {
	color.Printf("\n&4%s &r%d: %s\n", name, line, message)
}

// Get file name of path
func (p *Parser) GetName() string {
	return filepath.Base(p.Name)
}

func (p *Parser) GetPath() string {
	abs_path, err := filepath.Abs(p.Name)
	if err != nil {
		color.Println("&4error&r: error to get path of file")
		return p.Name
	}
	return abs_path
}

func (p *Parser) CurrentLine() string {
	maxLines := len(p.Data) - 1
	if p.Line > maxLines {
		return RemoveComments(p.Data[maxLines])
	}
	return RemoveComments(p.Data[p.Line])
}

func isEmpty(line string) bool {
	return line == ""
}

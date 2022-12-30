package alpha

import (
	"alpha/alpha/std/color"
	"fmt"
	"os"
	"path"
	"strings"
)

type Instance struct {
	Main   *Scope
	Parser Parser
}

func NewInstance(value string, name string) *Instance {
	parser := NewParser(path.Base(name), value)

	return &Instance{
		Parser: *parser,
		Main: &Scope{
			Variables:  BuiltinVariables,
			Code:       make([]Code, 0),
			Statements: make([]Block, 0),
		},
	}
}

func (sf *Instance) Execute() {
	sf.Parser.Compile(sf)

	if len(sf.Parser.Errors) > 0 {
		sf.Parser.ShowTraceback()
		os.Exit(-1)
	}

	for _, expr := range sf.Main.Code {
		sf.checkCode(&expr, false)
	}
}

func (sf *Instance) GetVariable(name string) (*Variable, error) {
	variable, ok := sf.Main.Variables[name]
	if !ok {
		return &Variable{}, fmt.Errorf("variable '%s' not found", name)
	}
	return variable, nil
}

func (sf *Instance) SetMeta(name string, meta VariableMeta) error {
	vars, err := sf.GetVariable(name)
	if err != nil {
		return err
	}

	vars.MetaContent = fmt.Sprintf("[%s:%s;Attr:%s]", meta.Type, meta.Content, strings.Join(meta.Attributes, ","))
	return nil
}

func (sf *Instance) SetVariableType(name, typevar string) error {
	vars, err := sf.GetVariable(name)
	if err != nil {
		return err
	}

	vars.Type = typevar

	return nil
}

func (sf *Instance) SetVariable(name string, content string) error {
	vars, err := sf.GetVariable(name)
	if err != nil {
		return fmt.Errorf("variable '%s' not found", name)
	} else if vars.IsConstant {
		return fmt.Errorf("variable '%s' is constant, it cannot be reassigned", name)
	}

	variable_type := GetVariableType(content)
	if variable_type == "invalid type" {
		variable_type = "string"
	}
	// if variable_type != "number" && variable_type != "null" && variable_type != "undefined" && variable_type != "nil" && variable_type != "array" {
	// 	variable_type = "string"
	// }
	vars.Content = content
	vars.MetaContent = fmt.Sprintf("[%s:%s;Attr:]", variable_type, content)
	vars.Type = variable_type
	return nil
}

/*
Create of reassign variable only if variable not constant
*/
func (sf *Instance) Set(name, content string, constant bool) error {
	_, err := sf.GetVariable(name)

	if err != nil {
		variable_type := GetVariableType(content)
		if variable_type != "number" && variable_type != "null" {
			variable_type = "string"
		}
		sf.Main.Variables[name] = &Variable{Name: name, Content: content, IsConstant: constant, Type: variable_type}
	} else {
		err := sf.SetVariable(name, content)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sf *Instance) CreateVariable(name string, content string, constant bool) error {
	_, err := sf.GetVariable(name)
	if err != nil {
		variable_type := GetVariableType(content)
		if variable_type == "invalid type" {
			variable_type = "string"
		}

		sf.Main.Variables[name] = &Variable{Name: name, Content: content, IsConstant: constant, Type: variable_type}
		sf.Main.Variables[name].MetaContent = fmt.Sprintf("[%s:%s;Attr:]", variable_type, content)
		return nil
	}
	return fmt.Errorf("variable %s already defined", name)
	// return fmt.Errorf("variable %s already defined", name)
}

func (sf *Instance) AddCode(value string, expr Expression, assign VariableDefinition, line int, kind ParserKind) {
	sf.Main.Code = append(sf.Main.Code, *CreateCode(value, expr, assign, Block{}, line, kind))
}

func (sf *Instance) AddStatement(value string, statement Block, line int) {
	sf.Main.Code = append(sf.Main.Code, *CreateCode(value, Expression{}, VariableDefinition{}, statement, line, STATEMENT_EXPRESSION))
}

func (sf *Instance) Print(value string) {
	color.Print(sf.format(value))
}

func (sf *Instance) Println(value string) {
	color.Println(sf.format(value))
}

func (sf *Instance) format(value string) string {
	varInstances := GetVariableInstances(value)
	sended := value
	if len(varInstances) > 0 {
		for _, variable := range varInstances {
			__var, _ := sf.GetVariable(variable.Name)
			sended = strings.ReplaceAll(value, fmt.Sprintf("{%s}", variable.Name), __var.Content)
		}
	}

	formatList := map[string]string{
		"\\n": "\n",
		"\\t": "\t",
		"\\r": "\r",
	}

	for key, value := range formatList {
		sended = strings.ReplaceAll(sended, key, value)
	}

	return sended
}

func (sf *Instance) checkCode(expr *Code, isStatement bool) {
	if expr.Kind == VARIABLE_DEFINITION {
		_, exists := sf.GetVariable(expr.Assignation.Name)
		if exists != nil {
			varInstances := GetVariableInstances(expr.Value)
			if len(varInstances) != 0 {
				varInst := varInstances[0]
				variable, exists := sf.GetVariable(varInst.Name)
				if exists != nil {
					message := color.Sprintf("%s\n\t&6>> &r%s", exists.Error(), expr.Value)
					sf.Parser.Error(message, expr.Line+1)
					os.Exit(-1)
				}
				sf.CreateVariable(expr.Assignation.Name, variable.Content, false)
			} else {
				sf.CreateVariable(expr.Assignation.Name, expr.Assignation.Content, false)
			}
		} else {
			varInstances := GetVariableInstances(expr.Value)
			content := expr.Assignation.Content
			if len(varInstances) != 0 {
				varInst := varInstances[0]
				variable, exists := sf.GetVariable(varInst.Name)
				if exists != nil {
					message := color.Sprintf("%s\n\t&6>> &r%s", exists.Error(), expr.Value)
					sf.Parser.Error(message, expr.Line+1)
					os.Exit(-1)
				}
				content = variable.Content
			}
			if err := sf.SetVariable(expr.Assignation.Name, content); err != nil {
				variable_position := strings.Index(expr.Value, expr.Assignation.Name)
				variable_length := len(expr.Assignation.Name)
				repeat := strings.Repeat(" ", variable_position+3)
				message := color.Sprintf("%s\n\t&6>> &r%s\n\t%s%s", err.Error(), expr.Value, repeat, strings.Repeat("^", variable_length))
				sf.Parser.Error(message, expr.Line+1)
				os.Exit(-1)
			}

		}
	}

	if expr.Kind == STATEMENT_EXPRESSION {
		expr.Statement.Executor(&expr.Statement, sf)
	}

	if expr.Kind == CODE_EXPRESSION {
		varInstances := expr.GetVariableInstances()
		for _, varInstance := range varInstances {
			variable, exists := sf.GetVariable(varInstance.Name)
			if exists != nil {
				if expr.Expr.Check {
					message := color.Sprintf("%s\n\t&6>> &r%s", exists.Error(), expr.Value)
					sf.Parser.Error(message, expr.Line+1)
					os.Exit(-1)
				}
			}
			if variable.IsConstant && expr.Expr.Assignations {
				repeat := strings.Repeat(" ", varInstance.Pos-1)
				repeat_arrow := strings.Repeat("^", len(varInstance.String()))
				repeat_value := color.Sprintf("\n\t%s\n\t%s%s", expr.Value, repeat, repeat_arrow)
				var_message := fmt.Sprintf("variable %s is constant", varInstance.Name)
				message := color.Sprintf("%s%s", var_message, repeat_value)
				sf.Parser.Error(message, expr.Line+1)
				os.Exit(-1)
			}

		}
		expr.Expr.Executor(expr, sf)
	}
}

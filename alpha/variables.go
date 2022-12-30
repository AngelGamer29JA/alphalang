package alpha

import (
	"alpha/alpha/io"
	"fmt"
	"regexp"
	"strings"
)

type VariableDefinition struct {
	Name    string
	Content string
	Type    string
}

type VariableMeta struct {
	Type       string
	Content    string
	Attributes []string
}

type Variable struct {
	Name        string
	Content     string
	Type        string
	MetaContent string
	Data        interface{}
	IsConstant  bool
}

var BuiltinVariables map[string]*Variable = map[string]*Variable{
	"io::args": CreateVariable("io::args", strings.Join(io.Args, ", "), "array", true),
	"io::cwd":  CreateVariable("io::cwd", io.GetCurrentWd(), "string", true),
}

func NewVariableDefinition(name, content, tipo string) *VariableDefinition {
	return &VariableDefinition{
		Name:    name,
		Content: content,
		Type:    tipo,
	}
}

func CreateVariable(name, content, tipo string, constant bool) *Variable {
	return NewVariable(
		name,
		content,
		tipo,
		fmt.Sprint("[", tipo, ":", content, ";Attr:", "]"),
		constant,
	)
}

func NewVariable(name, content, variabletype, meta string, constant bool) *Variable {
	return &Variable{
		Name:        name,
		Content:     content,
		Type:        variabletype,
		MetaContent: meta,
		IsConstant:  constant,
	}
}

func GetVariableType(str string) string {
	for key_pattern, regex_pattern := range patterns_types {
		ok, _ := regexp.MatchString(regex_pattern, str)
		if ok {
			return key_pattern
		}
	}
	return "invalid type"
}

func IsVariableDefinition(str string) bool {
	regex := regexp.MustCompile(VARIABLE_NAME_PATTERN)
	match := regex.MatchString(str)
	return match
}

func CheckVariableDefinition(value string) bool {
	regex := regexp.MustCompile(VARIABLE_DEFINE_PATTERN)
	match := regex.MatchString(value)
	return match
}

func GetVariableDefinition(value string) (*Variable, error) {
	if CheckVariableDefinition(value) {
		regexContent := regexp.MustCompile(ANY_PATTERN)
		Content := strings.ReplaceAll(regexContent.FindString(value), `"`, "")
		assignIndex := strings.IndexRune(value, '=') - 1

		VariableType := GetVariableType(Content)

		VariableName := value[0:assignIndex]
		return &Variable{
			Name:       VariableName,
			Content:    Content,
			IsConstant: false,
			Type:       VariableType,
		}, nil
	}
	return &Variable{}, fmt.Errorf("invalid variable definition")
}

func ReadMetaContent(value string) VariableMeta {
	meta := VariableMeta{}
	// [Type: content;Attrattr,attr]
	// [string:soy un string;Attr:]
	value = strings.TrimSuffix(strings.TrimPrefix(value, "["), "]")

	metaTypeIndex := strings.IndexRune(value, ':') + 1
	metaAttrIndex := strings.Index(value, ";Attr:")

	metaType := value[:metaTypeIndex-1]
	metaContent := value[metaTypeIndex:metaAttrIndex]
	metaAttributes := strings.Split(value[metaAttrIndex+6:], ",")

	meta.Type = metaType
	meta.Content = metaContent
	meta.Attributes = metaAttributes

	return meta
}

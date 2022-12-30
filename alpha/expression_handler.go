package alpha

import (
	"regexp"
	"strings"
)

/*
string = STRING_PATTERN
number = NUMBER_PATTERN // float64
integer = INT_PATTERN // int64
null = NULL_PATTERN
bool = BOOL_PATTERN

%(.+)%

sendln %string|number|null|bool%
*/

const (
	STRING_PATTERN            string = `\"\"|\"([^\"]*?)\"`
	NUMBER_PATTERN            string = `([0-9.]+)`
	INTEGER_PATTERN           string = `([0-9]+)`
	NULL_PATTERN              string = `(undefined|null|nil)`
	BOOL_PATTERN              string = `(true|false)`
	ARRAY_PATTERN             string = "^((" + NUMBER_PATTERN + "+|" + STRING_PATTERN + "|" + BOOL_PATTERN + "|" + NUMBER_PATTERN + "+|" + BOOL_PATTERN + "), )*(" + NUMBER_PATTERN + "+|" + STRING_PATTERN + "|" + BOOL_PATTERN + "|" + NUMBER_PATTERN + "+|" + NULL_PATTERN + ")|" + VARIABLE_INSTANCE_PATTERN + "$"
	STR_NUM                   string = "(" + STRING_PATTERN + "|" + NUMBER_PATTERN + ")"
	STR_NUM_BOOL              string = "(" + STR_NUM + "|" + BOOL_PATTERN + ")"
	STR_NUM_BOOL_NULL         string = "(" + STR_NUM_BOOL + "|" + NULL_PATTERN + ")"
	VARIABLE_NAME_PATTERN     string = `([$a-zA-Z_]+)\s*=\s*`
	VARIABLE_INSTANCE_PATTERN string = "({([a-zA-Z_$:]+)})"
	ANY_PATTERN               string = "(" + STR_NUM_BOOL_NULL + "|" + ARRAY_PATTERN + "|" + VARIABLE_INSTANCE_PATTERN + ")"
	VARIABLE_DEFINE_PATTERN   string = VARIABLE_NAME_PATTERN + ANY_PATTERN + "$"
)

var patterns_types = map[string]string{
	"string":  STRING_PATTERN,
	"number":  NUMBER_PATTERN,
	"integer": INTEGER_PATTERN,
	"bool":    BOOL_PATTERN,
	"null":    NULL_PATTERN,
}

// parse sendln %strings% - sendln (".*")
func ParseTemplate(input string) string {
	regex := regexp.MustCompile(`%([a-zA-Z|]+)%`)

	split_types := make([]string, 0)

	for _, m := range regex.FindAllString(input, -1) {
		m = strings.TrimSpace(strings.Replace(m, "%", "", -1))

		splitted := strings.Split(m, "|")
		// fmt.Println(splitted)
		split_types = append(split_types, splitted...)
	}

	input = regex.ReplaceAllString(input, "($1)")

	for _, value := range split_types {
		// fmt.Println(value)
		if value == "any" {
			input = strings.ReplaceAll(input, "any", ANY_PATTERN)
		} else if value == "variable" {
			input = strings.ReplaceAll(input, "variable", VARIABLE_INSTANCE_PATTERN)
		} else if ptype, ok := patterns_types[value]; ok {
			input = strings.ReplaceAll(input, value, ptype)
		}
	}

	return "(" + input + ")$"
}

package alpha

import "regexp"

type VariableInstance struct {
	Name string
	Pos  int
}

func (sf *VariableInstance) String() string {
	return "{" + sf.Name + "}"
}

func GetVariableInstances(str string) []VariableInstance {
	regex := regexp.MustCompile(`".+"`)
	str = regex.ReplaceAllString(str, "")

	regex = regexp.MustCompile("{([a-zA-Z_$:]+)}")
	matches := regex.FindAllStringSubmatchIndex(str, -1)
	results := make([]VariableInstance, len(matches))
	for i, match := range matches {
		value := str[match[2]:match[3]]
		results[i] = VariableInstance{
			Pos:  match[2],
			Name: value,
		}
	}
	return results
}

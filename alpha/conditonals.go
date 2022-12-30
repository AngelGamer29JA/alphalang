package alpha

type ConditonalExpression struct {
	Value     string
	Validator func(code *Code, instance *Instance) bool
}

type ConditionalStatement struct {
	Value string
}

type ConditionalBlock struct {
	Value     string
	Validator func(code *Code, instance *Instance) bool
	Start     int
	End       int
	Else      []Code
}

var BuiltinConditions []ConditonalExpression = []ConditonalExpression{
	*NewConditionalExpression("%any% is %any%", IsEqualTo),
}

const ConditionPrefix string = "if"

func NewConditionalExpression(value string, validator func(code *Code, instance *Instance) bool) *ConditonalExpression {
	return &ConditonalExpression{
		Value:     value,
		Validator: validator,
	}
}

func IsEqualTo(code *Code, instance *Instance) bool {
	return false
}

package alpha

import (
	"fmt"
	"strconv"
	"strings"
)

var BuiltinStatements = []Block{
	*NewBlock("loop %integer% times", loopblock),
	*NewBlock("loop", loopblock),
}

// loop %number% times
func loopblock(block *Block, instance *Instance) {
	content, _ := block.GetContentAt(0)

	if strings.HasPrefix(block.Value, "loop") {
		loopIndex := 0
	Loop:
		for {
			for _, code := range block.Scope.Code {
				if RemoveComments(code.Value) == "stop" {
					break Loop
				}
				index := fmt.Sprintf("%d", loopIndex)
				code.Value = strings.Replace(code.Value, "%loop-index%", index, -1)
				instance.checkCode(&code, true)
			}
			loopIndex += 1
		}
	} else {
		loopMax, _ := strconv.Atoi(content)

		for loopIndex := 0; loopIndex < loopMax; loopIndex++ {
			for _, code := range block.Scope.Code {
				index := fmt.Sprintf("%d", loopIndex)
				code.Value = strings.Replace(code.Value, "%loop-value%", index, -1)
				instance.checkCode(&code, true)
			}
		}
	}
}

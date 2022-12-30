package alpha

import (
	"alpha/alpha/io"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"golang.org/x/term"
)

var BuiltinExpressions = []Expression{
	*NewExpression(`sendln %any%`, send, true, false),                                   // print in consol with newline
	*NewExpression(`send %any%`, send, true, false),                                     // print in console without newline
	*NewExpression(`prompt %string% stored in %variable%`, prompt, true, true),          // prompt in console
	*NewExpression(`set %any% to %variable%`, setto, false, true),                       // reassign variable to another value
	*NewExpression(`get type %variable% to %variable%`, getTypeTo, true, true),          // Obtener el tipo de una variable
	*NewExpression(`wait %integer% secon(ds|d)`, wait, true, false),                     // Esperar un tiempo determinado en segundos
	*NewExpression(`make dir to %string%`, makeDir, true, false),                        // Crea una carpeta
	*NewExpression(`create file %string%`, createFile, true, false),                     // Crear un archivo
	*NewExpression(`create file %string% stored in %variable%`, createFile, true, true), //
	*NewExpression(`read file %string% stored in %variable%`, readFile, true, true),
	*NewExpression(`(write|writeln) %string% in %string%`, writeFile, true, false),
	*NewExpression(`(write|writeln) %string% stored in %variable%`, writeFile, true, false),
	*NewExpression(`add %number% to %variable%`, addTo, true, true),
	*NewExpression(`clear file %string%`, clearFile, true, false),
	*NewExpression(`readline stored in %variable%`, readLine, true, true),
	*NewExpression(`delete %string%`, deleteFileOrFolder, true, false),
	*NewExpression(`replace %string% to %string% stored in %variable%`, replaceTo, true, false),
	*NewExpression(`writeln %string% in %string% at %integer%`, writeFile, true, false),
	*NewExpression(`clear terminal`, clearTerminal, true, false),
}

func send(code *Code, instance *Instance) {
	var content string
	content, err := code.GetContentAt(0)
	varInstances := code.GetVariableInstances()

	if err != nil { // Any content
		_var, _ := instance.GetVariable(varInstances[0].Name)
		content = _var.Content
	} else {
		for _, ins := range varInstances {
			_var, _ := instance.GetVariable(ins.Name)
			content = strings.Replace(content, ins.String(), _var.Content, -1)
		}
	}

	if strings.HasPrefix(code.Value, "sendln") {
		instance.Println(content)
	} else if strings.HasPrefix(code.Value, "send") {
		instance.Print(content)
	}
}

func prompt(code *Code, instance *Instance) {
	arg, _ := code.GetContentAt(0)
	_varInstanced := code.GetVariableInstances()[0]

	fmt.Printf("%s ", arg)
	return_value := io.ReadLine()
	instance.SetVariable(_varInstanced.Name, return_value)
}

func setto(code *Code, instance *Instance) {
	_varInstanced := code.GetVariableInstances()[0]
	value, _ := code.GetContentAt(0)

	instance.Set(_varInstanced.Name, value, false)
}

func getTypeTo(code *Code, instance *Instance) {
	_varInstanced := code.GetVariableInstances()[0]
	_varTo := code.GetVariableInstances()[1]
	__var, _ := instance.GetVariable(_varInstanced.Name)
	instance.SetVariable(_varTo.Name, __var.Type)
}

func wait(code *Code, instance *Instance) {
	seconds, _ := code.GetContentAt(0)
	timeAwaitted, _ := strconv.Atoi(seconds)
	time.Sleep(time.Duration(timeAwaitted) * time.Second)
}

func makeDir(code *Code, instance *Instance) {
	dirpath, _ := code.GetContentAt(0)
	fmt.Println("Making dir in : ", dirpath)
}

func createFile(code *Code, instance *Instance) {
	fpath, _ := code.GetContentAt(0)
	value := code.GetVariableInstances()
	if len(value) > 0 {
		stored := value[0]
		if err := io.CreateFile(fpath, ""); err != nil {
			instance.Parser.Error(err.Error(), code.Line)
			return
		}

		fileContent := io.Readfile(fpath)
		metaNew := VariableMeta{
			Type:       "file",
			Content:    fileContent,
			Attributes: []string{"filename=" + fpath},
		}

		instance.SetVariable(stored.Name, fileContent)
		instance.SetMeta(stored.Name, metaNew)
		instance.SetVariableType(stored.Name, "file")
		//fmt.Printf("File %s created and stored in %s", fpath, value[0].Name)
	} else {
		io.CreateFile(fpath, "")
		//fmt.Printf("File %s created", fpath)
	}
}

func readFile(code *Code, instance *Instance) {
	fpath, _ := code.GetContentAt(0)
	stored := code.GetVariableInstances()[0]
	fileContent := io.Readfile(fpath)

	metaNew := VariableMeta{
		Type:       "file",
		Content:    fileContent,
		Attributes: []string{"filename=" + fpath},
	}

	instance.SetVariable(stored.Name, fileContent)
	instance.SetMeta(stored.Name, metaNew)
	instance.SetVariableType(stored.Name, "file")
}

func writeFile(code *Code, instance *Instance) {
	fcontent, _ := code.GetContentAt(0) // write data to file
	fcontent = instance.format(fcontent)
	variables := code.GetVariableInstances()
	if len(variables) > 0 {
		variable, _ := instance.GetVariable(variables[0].Name)
		meta := ReadMetaContent(variable.MetaContent)
		fnidx := strings.IndexRune(meta.Attributes[0], '=')
		filename := meta.Attributes[0][fnidx+1:]
		if strings.HasPrefix(code.Value, "writeln") {
			io.WriteLine(filename, fcontent)
		} else if strings.HasPrefix(code.Value, "write") {
			io.Write(filename, fcontent)
		}
	} else {
		fpath, _ := code.GetContentAt(1)
		if strings.HasPrefix(code.Value, "writeln") {
			io.WriteLine(fpath, fcontent)
		} else if strings.HasPrefix(code.Value, "write") {
			io.Write(fpath, fcontent)
		}
		if strings.HasSuffix(code.Expr.Value, "at %integer%") {
			idxContent, _ := code.GetContentAt(2) // 0 1 2
			indexFile, _ := strconv.Atoi(idxContent)
			io.WriteLineAt(fpath, fcontent, indexFile)
		}
	}
}

func addTo(code *Code, instance *Instance) {
	countValue, _ := code.GetContentAt(0)
	countNumber, _ := strconv.ParseFloat(countValue, 32)

	varInst := code.GetVariableInstances()[0]
	varToAdd, _ := instance.GetVariable(varInst.Name)
	countVar, _ := strconv.ParseFloat(varToAdd.Content, 32)
	countVar += countNumber
	varToAdd.Content = strconv.FormatFloat(countVar, 'f', -1, 32)
}

func clearFile(code *Code, instance *Instance) {
	fpath, _ := code.GetContentAt(0)
	if !io.FileExists(fpath) {
		instance.Parser.Error(fmt.Sprintf("file '%s' not exists, please check the file name", fpath), code.Line)
		os.Exit(-1)
	}
	io.CreateFile(fpath, "")
}

func readLine(code *Code, instance *Instance) {
	variableStored := code.GetVariableInstances()[0]

	lineContent := io.ReadLine()
	instance.SetVariable(variableStored.Name, lineContent)
}

func deleteFileOrFolder(code *Code, instance *Instance) {
	fpath, _ := code.GetContentAt(0)
	io.Delete(fpath)
}

func replaceTo(code *Code, instance *Instance) {
	matchString, _ := code.GetContentAt(0)
	replaceString, _ := code.GetContentAt(1)
	storedIn := code.GetVariableInstances()[0]
	variableStored, _ := instance.GetVariable(storedIn.Name)

	newString := strings.Replace(variableStored.Content, matchString, replaceString, -1)
	instance.SetVariable(variableStored.Name, newString)
}

func clearTerminal(code *Code, instance *Instance) {
	_, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		code.Error(instance.Parser.Name, err.Error())
		return
	}

	// Enviar códigos ANSI para borrar la pantalla y posicionar el cursor en la esquina superior izquierda
	fmt.Fprintf(colorable.NewColorableStdout(), "\033[2J")           // Borrar pantalla
	fmt.Fprintf(colorable.NewColorableStdout(), "\033[%d;%dH", 0, 0) // Posicionar cursor
}

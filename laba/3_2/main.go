package main

import (
	"bufio"
	"fmt"
	. "laba/Tree"
	. "laba/types"
	"os"
	"strings"
	"unicode"
)

func main() {
	variables := make(map[string]*Variable)
	functions := make(map[string]*Function)

	var fileName string
	fmt.Print("Введите название файла: ")
	if _, err := fmt.Scan(&fileName); err != nil {
		fmt.Println("Error fmt.Scan - ", err)
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error open file - %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		var name, body string
		for index, char := range line {
			if char == ':' {
				name = line[:index]
				body = line[index : len(line)-1]
				ParseFunction(name, body, functions)
				break
			} else if char == '=' {
				name = line[:index]
				body = line[index : len(line)-1]
				ParseVariable(name, body, variables, functions)
				break
			}
		}

		if name == "" && body == "" {
			ParsePrint(line, variables)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func ParseVariable(name, body string, variables map[string]*Variable, functions map[string]*Function) {
	var vT, varName string
	flag := false
	for _, char := range name {
		if char == '(' {
			flag = true
			continue
		} else if char == ')' {
			flag = false
			continue
		}

		if flag {
			vT += string(char)
		} else {
			varName += string(char)
		}
	}

	var variableType VarType
	switch vT {
	case "i":
		variableType = Integer
	case "f":
		variableType = Float
	default:
		variableType = Float
	}

	body = strings.TrimSpace(body)

	tokens := Tokenize(body)
	root := Parse(tokens)
	resultBodyVariable := Evaluate(root, variables, functions)

	variables[varName] = NewVariable(varName, variableType, resultBodyVariable)
}

func ParseFunction(name, body string, functions map[string]*Function) {
	args := []string{}
	var nameFunction string

	flagArgs := false
	for index := 0; index < len(name); {
		char := rune(name[index])
		if char == '(' {
			nameFunction = name[0:index]
			flagArgs = true
			index++
		} else if flagArgs && unicode.IsLetter(char) {
			start := index
			for index < len(name) && (unicode.IsLetter(rune(name[index])) || unicode.IsDigit(rune(name[index]))) {
				index++
			}
			args = append(args, name[start:index])
		} else {
			index++
		}
	}

	tokens := Tokenize(body)
	root := Parse(tokens)

	functions[nameFunction] = NewFunction(nameFunction, root, args)
}

func ParsePrint(line string, variables map[string]*Variable) {
	lineLen := len(line)
	nameArg := ""
	flagRparen := false
	for i := 0; i < lineLen; {
		char := rune(line[i])

		if char == ' ' {
			flagRparen = true
			i++
		} else if flagRparen && unicode.IsLetter(char) {
			start := i
			for i < lineLen && (unicode.IsLetter(char) || unicode.IsDigit(char)) {
				i++
			}

			nameArg = line[start:i]
		} else {
			i++
		}
	}

	if !flagRparen {
		fmt.Println("Выводим все переменные - print")
		for _, variable := range variables {
			fmt.Printf("%s = %f\n", variable.Name, variable.Value)
		}
	} else {
		value, ok := variables[nameArg]
		if !ok {
			fmt.Println("Ошибка в print - такой переменной не существует")
			return
		}
		fmt.Printf("%s = %f\n", nameArg, value.Value)
	}
}

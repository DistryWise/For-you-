package main

import (
	"bufio"
	"fmt"
	. "laba/3_1/scope"
	"os"
)

func main() {
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

	// Слайс где будут храниться области видимости
	var sliceScopes []Scope
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		switch line {
		case "{":
			sliceScopes = append(sliceScopes, NewScope())
		case "}":
			sliceScopes = sliceScopes[:len(sliceScopes)-1]
		case "ShowVar;":
			showVariables := make(map[string]string)

			for _, element := range sliceScopes {
				for key, value := range element {
					showVariables[key] = value
				}
			}
			fmt.Println("ShowVar: ", showVariables)
		default:
			var variableName, variableValue string
			eqBool := false

			for _, char := range line {
				if char == '=' {
					eqBool = true
					continue
				}
				if eqBool {
					variableValue += string(char)
				} else {
					variableName += string(char)
				}
			}

			if len(sliceScopes) > 0 {
				scope := sliceScopes[len(sliceScopes)-1]
				scope[variableName] = variableValue
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v", err)
	}
}

package main

import (
	"bufio"
	"fmt"
	. "laba/2/array"
	"os"
	"strconv"
	"strings"
)

func main() {
	var containerArrays map[string]Array = make(map[string]Array)

	var fileName string
	fmt.Print("Введите название файла: ")
	if _, err := fmt.Scan(&fileName); err != nil {
		fmt.Printf("Error scan Stdin: %v", err)
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error open file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		commandSplit := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == ',' || r == ';' || r == '(' || r == ')'
		})

		searchCommand(commandSplit, containerArrays)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v", err)
	}
}

func searchCommand(command []string, arrays map[string]Array) {
	if len(command) == 0 {
		fmt.Println("Empty command!!!")
		return
	}

	nameCommand := command[0]

	switch nameCommand {
	case "load":
		if err := LoadArray(command[1], command[2], arrays); err != nil {
			fmt.Printf("Error: %v", err)
		}
	case "save":
		if err := SaveArray(command[1], command[2], arrays); err != nil {
			fmt.Printf("Error: %v", err)
		}
	case "rand":
		count, err := strconv.Atoi(command[2])
		if err != nil {
			fmt.Printf("Error converting count: %v", err)
			return
		}
		lb, err := strconv.Atoi(command[3])
		if err != nil {
			fmt.Printf("Error converting lb: %v", err)
			return
		}
		rb, err := strconv.Atoi(command[4])
		if err != nil {
			fmt.Printf("Error converting rb: %v", err)
			return
		}
		if err := RandArray(command[1], count, lb, rb, arrays); err != nil {
			fmt.Printf("Error: %v", err)
		}
	case "concat":
		if err := ConcatArray(command[1], command[2], arrays); err != nil {
			fmt.Printf("Error: %v", err)
		}
	case "free":
		if err := FreeArray(command[1], arrays); err != nil {
			fmt.Printf("Error: %v", err)
		}
	case "remove":
		index, err := strconv.Atoi(command[2])
		if err != nil {
			fmt.Printf("Error converting index: %v", err)
			return
		}
		count, err := strconv.Atoi(command[3])
		if err != nil {
			fmt.Printf("Error converting count: %v", err)
			return
		}
		if err := RemoveArray(command[1], index, count, arrays); err != nil {
			fmt.Printf("Error: %v", err)
		}
	case "copy":
		lb, err := strconv.Atoi(command[2])
		if err != nil {
			fmt.Printf("Error converting lb: %v", err)
			return
		}
		rb, err := strconv.Atoi(command[3])
		if err != nil {
			fmt.Printf("Error converting rb: %v", err)
			return
		}
		if err := CopyArray(command[1], command[4], lb, rb, arrays); err != nil {
			fmt.Printf("Error: %v", err)
		}
	case "sort":
		if err := SortArray(command[1], arrays); err != nil {
			fmt.Printf("Error: %v", err)
		}
	case "shuffle":
		if err := ShuffleArray(command[1], arrays); err != nil {
			fmt.Printf("Error: %v", err)
		}
	case "stats":
		if err := StatsArray(command[1], arrays); err != nil {
			fmt.Printf("Error: %v", err)
		}
	case "print":
		if len(command) == 3 {
			if err := PrintArray(command[1], command[2], arrays); err != nil {
				fmt.Printf("Error: %v", err)
			}
		} else {
			lb, err := strconv.Atoi(command[2])
			if err != nil {
				fmt.Printf("Error converting lb: %v", err)
				return
			}
			rb, err := strconv.Atoi(command[3])
			if err != nil {
				fmt.Printf("Error converting rb: %v", err)
				return
			}
			if err := PrintRangeArray(command[1], lb, rb, arrays); err != nil {
				fmt.Printf("Error: %v", err)
			}
		}
	default:
		fmt.Println("Invalid command!!!")
	}
}

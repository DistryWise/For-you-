package parseTree

import (
	"fmt"
	. "project_dict/laba3_2/types"
	"strconv"
	"unicode"
)

type Token struct {
	Type, Value string
}

func Tokenize(expression string) []Token {
	var tokens []Token
	expressionLen := len(expression)
	for i := 0; i < expressionLen; {
		ch := rune(expression[i])

		if unicode.IsDigit(ch) {
			start := i
			for i < expressionLen && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				i++
			}
			tokens = append(tokens, Token{Type: NUMBER, Value: expression[start:i]})
		} else if unicode.IsLetter(ch) {
			start := i
			for i < expressionLen && (unicode.IsLetter(rune(expression[i])) || unicode.IsDigit(rune(expression[i]))) {
				i++
			}
			value := expression[start:i]

			if i < expressionLen && expression[i] == '(' {
				tokens = append(tokens, Token{Type: FUNCTION, Value: value})
			} else {
				tokens = append(tokens, Token{Type: IDENT, Value: expression[start:i]})
			}
		} else if ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == ',' {
			tokens = append(tokens, Token{Type: OPERATOR, Value: string(ch)})
			i++
		} else if ch == '(' {
			tokens = append(tokens, Token{Type: LPAREN, Value: string(ch)})
			i++
		} else if ch == ')' {
			tokens = append(tokens, Token{Type: RPAREN, Value: string(ch)})
			i++
		} else {
			i++
		}
	}
	return tokens
}

func Parse(tokens []Token) *Node {
	return parseExpression(tokens, 0).Node
}

type ParseResult struct {
	Node  *Node
	Index int
}

func parseExpression(tokens []Token, index int) ParseResult {
	leftResult := parseTerm(tokens, index)
	node := leftResult.Node
	index = leftResult.Index

	tokensLen := len(tokens)
	for index < tokensLen && (tokens[index].Type == OPERATOR && (tokens[index].Value == "+" || tokens[index].Value == "-")) {
		op := tokens[index]
		rightResult := parseTerm(tokens, index+1)
		node = &Node{Type: OPERATOR, Value: op.Value, Left: node, Right: rightResult.Node}
		index = rightResult.Index
	}

	return ParseResult{Node: node, Index: index}
}

func parseTerm(tokens []Token, index int) ParseResult {
	leftResult := parseFactor(tokens, index)
	node := leftResult.Node
	index = leftResult.Index

	for index < len(tokens) && (tokens[index].Type == OPERATOR && (tokens[index].Value == "*" || tokens[index].Value == "/")) {
		op := tokens[index]
		rightResult := parseFactor(tokens, index+1)
		node = &Node{Type: OPERATOR, Value: op.Value, Left: node, Right: rightResult.Node}
		index = rightResult.Index
	}

	return ParseResult{Node: node, Index: index}
}

func parseFactor(tokens []Token, index int) ParseResult {
	token := tokens[index]
	if token.Type == NUMBER {
		return ParseResult{Node: &Node{Type: NUMBER, Value: token.Value}, Index: index + 1}
	} else if token.Type == FUNCTION {
		funcName := token.Value
		index += 2 // пропускаем идентификатор и левую скобку
		args := []*Node{}
		for tokens[index].Type != RPAREN {
			argResult := parseExpression(tokens, index)
			args = append(args, argResult.Node)
			index = argResult.Index
			if tokens[index].Type == OPERATOR && tokens[index].Value == "," {
				index++ // пропускаем запятую
			}
		}
		index++                                                    // пропускаем правую скобку
		node := &Node{Type: FUNCTION, Value: funcName, Args: args} // в данном примере поддерживается только один аргумент
		return ParseResult{Node: node, Index: index}
	} else if token.Type == IDENT {
		return ParseResult{Node: &Node{Type: IDENT, Value: token.Value}, Index: index + 1}
	} else if token.Type == LPAREN {
		result := parseExpression(tokens, index+1)
		if tokens[result.Index].Type == RPAREN {
			return ParseResult{Node: result.Node, Index: result.Index + 1}
		}
	}
	return ParseResult{Node: nil, Index: index}
}

func printParseTree(node *Node, indent string) {
	if node == nil {
		return
	}

	fmt.Println(indent + node.Type + ": " + node.Value)
	if node.Left != nil {
		printParseTree(node.Left, indent+"  ")
	}
	if node.Right != nil {
		printParseTree(node.Right, indent+"  ")
	}

	for _, arg := range node.Args {
		printParseTree(arg, indent+"  ")
	}
}

func Evaluate(node *Node, variables map[string]*Variable, functions map[string]*Function) float64 {
	switch node.Type {
	case NUMBER:
		val, _ := strconv.ParseFloat(node.Value, 64)
		return val
	case IDENT:
		value, ok := variables[node.Value]
		if !ok {
			fmt.Printf("Error - variable not found")
			return -1
		}
		return value.Value
	case FUNCTION:
		args := make([]float64, len(node.Args))
		for i, argNode := range node.Args {
			args[i] = Evaluate(argNode, variables, functions)
		}
		if function, exists := functions[node.Value]; exists {
			if len(function.Arg) != len(args) {
				fmt.Println("Error - неправильные аргументы функции")
				panic(fmt.Sprintf("Error - неправильные аргументы функции - %s", function.Name))
			}
			tempVariables := make(map[string]*Variable)
			for index, element := range args {
				tempVariables[function.Arg[index]] = NewVariable(function.Arg[index], Null, element)
			}
			return Evaluate(function.Body, tempVariables, functions)
		}
		fmt.Printf("Error - function not found")
		return -1
	case OPERATOR:
		leftVal := Evaluate(node.Left, variables, functions)
		rightVal := Evaluate(node.Right, variables, functions)
		switch node.Value {
		case "+":
			return leftVal + rightVal
		case "-":
			return leftVal - rightVal
		case "*":
			return leftVal * rightVal
		case "/":
			return leftVal / rightVal
		}
	}
	return 0
}

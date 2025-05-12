package RPN

import (
	"fmt"
	"project_yandex_lms/lms-calc-with-gorutine/variables"
	"strconv"
	"unicode"
)

func InfixToRPN(expression string) ([]string, error) {
	var output []string
	var stack []string
	var currentNum string

	isOperator := func(s string) bool {
		_, ok := variables.Operators[s]
		return ok
	}

	precedence := func(op string) int {
		return variables.Operators[op]
	}

	for i, char := range expression {
		token := string(char)

		if token == " " {
			continue
		}

		if _, err := strconv.Atoi(token); err == nil {

			currentNum += token

			if i == len(expression)-1 || !unicode.IsDigit(rune(expression[i+1])) {
				output = append(output, currentNum)
				currentNum = ""
			}
		} else if isOperator(token) {

			for len(stack) > 0 && isOperator(stack[len(stack)-1]) && precedence(token) <= precedence(stack[len(stack)-1]) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {

			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("несоответствие скобок")
			}
			stack = stack[:len(stack)-1]
		} else {
			return nil, fmt.Errorf("недопустимый символ: %s", token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" || stack[len(stack)-1] == ")" {
			return nil, fmt.Errorf("несоответствие скобок")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

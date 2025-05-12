package calc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	var stack []float64
	var ops []rune

	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	applyOp := func(op rune) error {
		if len(stack) < 2 {
			return fmt.Errorf("invalid character low operators %c", op)
		}
		b, a := stack[len(stack)-1], stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		switch op {
		case '+':
			stack = append(stack, a+b)
		case '-':
			stack = append(stack, a-b)
		case '*':
			stack = append(stack, a*b)
		case '/':
			if b == 0 {
				return fmt.Errorf("invalid character toshe delit of 0")
			}
			stack = append(stack, a/b)
		}
		return nil
	}

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		if unicode.IsDigit(char) {
			numStr := string(char)
			for i+1 < len(expression) && (unicode.IsDigit(rune(expression[i+1])) || expression[i+1] == '.') {
				i++
				numStr += string(expression[i])
			}
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid character of digit: %s", numStr)
			}
			stack = append(stack, num)

		} else if char == '(' {
			ops = append(ops, char)

		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				if err := applyOp(op); err != nil {
					return 0, err
				}
			}
			if len(ops) == 0 {
				return 0, fmt.Errorf("invalid character of delit zero")
			}
			ops = ops[:len(ops)-1]

		} else if precedence[char] > 0 {
			for len(ops) > 0 && precedence[ops[len(ops)-1]] >= precedence[char] {
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				if err := applyOp(op); err != nil {
					return 0, err
				}
			}
			ops = append(ops, char)
		} else {
			return 0, fmt.Errorf("invalid character of chars: %c", char)
		}
	}

	for len(ops) > 0 {
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]
		if err := applyOp(op); err != nil {
			return 0, err
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid character of len stack")
	}
	return stack[0], nil
}

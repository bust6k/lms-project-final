package RPN

import (
	"project_yandex_lms/lms-calc-with-gorutine/variables"
	"testing"
)

func TestInfixToRPN(t *testing.T) {
	
	originalOperators := variables.Operators
	defer func() { variables.Operators = originalOperators }()

	// Настраиваем тестовые операторы
	variables.Operators = map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"^": 3,
	}

	tests := []struct {
		name        string
		expression  string
		expected    []string
		expectError bool
		errMsg      string
	}{
		// Базовые арифметические операции
		{
			name:       "simple addition",
			expression: "3 + 4",
			expected:   []string{"3", "4", "+"},
		},
		{
			name:       "simple multiplication",
			expression: "3 * 4",
			expected:   []string{"3", "4", "*"},
		},
		{
			name:       "mixed operations",
			expression: "3 + 4 * 2",
			expected:   []string{"3", "4", "2", "*", "+"},
		},
		{
			name:       "with parentheses",
			expression: "(3 + 4) * 2",
			expected:   []string{"3", "4", "+", "2", "*"},
		},
		{
			name:       "complex expression",
			expression: "3 + 4 * 2 / (1 - 5) ^ 2",
			expected:   []string{"3", "4", "2", "*", "1", "5", "-", "2", "^", "/", "+"},
		},
		{
			name:       "multiple digits",
			expression: "123 + 456",
			expected:   []string{"123", "456", "+"},
		},
		{
			name:       "unary minus",
			expression: "-3 + 4",
			expected:   []string{"-3", "4", "+"},
		},
		{
			name:       "exponentiation right associativity",
			expression: "2 ^ 3 ^ 2",
			expected:   []string{"2", "3", "2", "^", "^"},
		},

		// Ошибочные случаи
		{
			name:        "mismatched parentheses",
			expression:  "(3 + 4",
			expectError: true,
			errMsg:      "несоответствие скобок",
		},
		{
			name:        "invalid character",
			expression:  "3 $ 4",
			expectError: true,
			errMsg:      "недопустимый символ: $",
		},
		{
			name:        "empty expression",
			expression:  "",
			expectError: true,
			errMsg:      "",
		},
		{
			name:        "only operator",
			expression:  "+",
			expectError: true,
			errMsg:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := InfixToRPN(tt.expression)

			if tt.expectError {
				if err == nil {
					t.Error("ожидалась ошибка, но не получена")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("неправильное сообщение об ошибке: получено '%v', ожидается '%v'", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("неожиданная ошибка: %v", err)
				return
			}

			if !compareStringSlices(result, tt.expected) {
				t.Errorf("результат не соответствует ожидаемому значению: получено %v, ожидается %v", result, tt.expected)
			}
		})
	}
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestInfixToRPNEdgeCases(t *testing.T) {
	originalOperators := variables.Operators
	defer func() { variables.Operators = originalOperators }()
	variables.Operators = map[string]int{
		"+": 1,
		"*": 2,
	}

	tests := []struct {
		name       string
		expression string
		expected   []string
	}{
		{
			name:       "multiple spaces",
			expression: "   3   +   4   *  2   ",
			expected:   []string{"3", "4", "2", "*", "+"},
		},
		{
			name:       "no spaces",
			expression: "3+4*2",
			expected:   []string{"3", "4", "2", "*", "+"},
		},
		{
			name:       "nested parentheses",
			expression: "((3 + 4) * 2)",
			expected:   []string{"3", "4", "+", "2", "*"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := InfixToRPN(tt.expression)
			if err != nil {
				t.Errorf("неожиданная ошибка: %v", err)
				return
			}

			if !compareStringSlices(result, tt.expected) {
				t.Errorf("результат не соответствует ожидаемому значению: получено %v, ожидается %v", result, tt.expected)
			}
		})
	}
}

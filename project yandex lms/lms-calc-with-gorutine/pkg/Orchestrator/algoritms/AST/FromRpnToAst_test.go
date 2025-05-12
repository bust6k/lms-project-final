package AST

import (
	"fmt"
	"project_yandex_lms/lms-calc-with-gorutine/entites"
	"project_yandex_lms/lms-calc-with-gorutine/variables"
	"testing"
)

func init() {
	// Инициализируем операторы для тестов
	variables.Operators = map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}
}

func TestBuildAST(t *testing.T) {
	tests := []struct {
		name        string
		rpn         []string
		want        *entites.ASTNode
		expectError bool
		errMsg      string
	}{
		{
			name: "simple addition",
			rpn:  []string{"2", "3", "+"},
			want: &entites.ASTNode{
				Type:  variables.OperatorNode,
				Value: "+",
				Left: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "2",
				},
				Right: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "3",
				},
			},
		},
		{
			name: "complex expression",
			rpn:  []string{"2", "3", "4", "*", "+"},
			want: &entites.ASTNode{
				Type:  variables.OperatorNode,
				Value: "+",
				Left: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "2",
				},
				Right: &entites.ASTNode{
					Type:  variables.OperatorNode,
					Value: "*",
					Left: &entites.ASTNode{
						Type:  variables.NumberNode,
						Value: "3",
					},
					Right: &entites.ASTNode{
						Type:  variables.NumberNode,
						Value: "4",
					},
				},
			},
		},
		{
			name: "division operation",
			rpn:  []string{"6", "2", "/"},
			want: &entites.ASTNode{
				Type:  variables.OperatorNode,
				Value: "/",
				Left: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "6",
				},
				Right: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "2",
				},
			},
		},
		{
			name: "single number",
			rpn:  []string{"42"},
			want: &entites.ASTNode{
				Type:  variables.NumberNode,
				Value: "42",
			},
		},
		{
			name:        "missing operands",
			rpn:         []string{"+"},
			expectError: true,
			errMsg:      "недостаточно операндов для оператора: +",
		},
		{
			name:        "invalid RPN expression",
			rpn:         []string{"2", "3", "+", "4"},
			expectError: true,
			errMsg:      "недопустимое RPN выражение",
		},
		{
			name:        "empty input",
			rpn:         []string{},
			expectError: true,
			errMsg:      "недопустимое RPN выражение",
		},
		{
			name: "multi-digit numbers",
			rpn:  []string{"123", "456", "+"},
			want: &entites.ASTNode{
				Type:  variables.OperatorNode,
				Value: "+",
				Left: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "123",
				},
				Right: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "456",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildAST(tt.rpn)

			if tt.expectError {
				if err == nil {
					t.Fatal("ожидалась ошибка, но не получена")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("неправильное сообщение ошибки: получено '%v', ожидается '%v'", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			if !compareAST(got, tt.want) {
				t.Errorf("полученное AST не соответствует ожидаемому")
				printAST(got, "Получено:")
				printAST(tt.want, "Ожидалось:")
			}
		})
	}
}

// compareAST рекурсивно сравнивает два AST-дерева
func compareAST(a, b *entites.ASTNode) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Type != b.Type || a.Value != b.Value {
		return false
	}
	return compareAST(a.Left, b.Left) && compareAST(a.Right, b.Right)
}

// printAST рекурсивно печатает AST для отладки
func printAST(node *entites.ASTNode, prefix string) {
	if node == nil {
		return
	}
	fmt.Printf("%s {Type: %v, Value: %s}\n", prefix, node.Type, node.Value)
	if node.Left != nil {
		printAST(node.Left, prefix+"  L:")
	}
	if node.Right != nil {
		printAST(node.Right, prefix+"  R:")
	}
}

func TestBuildASTWithCustomOperators(t *testing.T) {
	// Сохраняем оригинальные операторы
	originalOperators := variables.Operators
	defer func() { variables.Operators = originalOperators }()

	// Добавляем тестовые операторы
	variables.Operators = map[string]int{
		"add": 1,
		"mul": 2,
	}

	tests := []struct {
		name string
		rpn  []string
		want *entites.ASTNode
	}{
		{
			name: "custom operators",
			rpn:  []string{"2", "3", "add"},
			want: &entites.ASTNode{
				Type:  variables.OperatorNode,
				Value: "add",
				Left: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "2",
				},
				Right: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "3",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildAST(tt.rpn)
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			if !compareAST(got, tt.want) {
				t.Errorf("полученное AST не соответствует ожидаемому")
			}
		})
	}
}

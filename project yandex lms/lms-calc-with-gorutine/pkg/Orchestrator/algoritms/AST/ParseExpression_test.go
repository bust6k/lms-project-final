package AST_test

import (
	"project_yandex_lms/lms-calc-with-gorutine/entites"
	"project_yandex_lms/lms-calc-with-gorutine/pkg/Orchestrator/algoritms/AST"
	"project_yandex_lms/lms-calc-with-gorutine/variables"
	"testing"
)

func TestSplitAST(t *testing.T) {
	tests := []struct {
		name     string
		input    *entites.ASTNode
		expected []entites.Task
	}{
		{
			name: "Single operation",
			input: &entites.ASTNode{
				Type:  variables.OperatorNode,
				Value: "+",
				Left: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "1",
				},
				Right: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "2",
				},
			},
			expected: []entites.Task{
				{
					Id:        0,
					Arg1:      1,
					Operation: "+",
					Arg2:      2,
				},
			},
		},
		{
			name: "Nested operations",
			input: &entites.ASTNode{
				Type:  variables.OperatorNode,
				Value: "*",
				Left: &entites.ASTNode{
					Type:  variables.OperatorNode,
					Value: "+",
					Left:  &entites.ASTNode{Type: variables.NumberNode, Value: "1"},
					Right: &entites.ASTNode{Type: variables.NumberNode, Value: "2"},
				},
				Right: &entites.ASTNode{
					Type:  variables.NumberNode,
					Value: "3",
				},
			},
			expected: []entites.Task{
				{
					Id:        0,
					Arg1:      1,
					Operation: "+",
					Arg2:      2,
				},
				{
					Id:        1,
					Arg1:      (1 + 2),
					Operation: "*",
					Arg2:      3,
				},
			},
		},
		// Добавьте дополнительные тесты по мере необходимости
	}

	for _, tt := range tests {
		// Сбросим счетчик перед каждым тестом
		variables.CurrentCountOfUnprocessedUserExpressions = 0

		t.Run(tt.name, func(t *testing.T) {
			got := AST.SplitAST(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("expected %d tasks, got %d", len(tt.expected), len(got))
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("expected task %v, got %v", tt.expected[i], got[i])
				}
			}
		})
	}
}

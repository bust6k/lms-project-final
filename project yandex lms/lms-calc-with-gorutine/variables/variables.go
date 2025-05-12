package variables

import (
	"project_yandex_lms/lms-calc-with-gorutine/entites"
	"project_yandex_lms/lms-calc-with-gorutine/models"
)

var (
	CurrentCountOfUnprocessedUserExpressions int
	Expressions                              []models.ProcessedExpression

	CurrentTask entites.Task
	//the tasks
	TheTasks []entites.Task

	Operators = map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}
)

const (
	NumberNode entites.NodeType = iota
	OperatorNode
)

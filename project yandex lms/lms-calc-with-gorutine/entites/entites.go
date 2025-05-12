package entites

import (
	"time"
)

type NodeType int

type UnprocessedUserExpression struct {
	Id         int    `json:"id"`
	Expression string `json:"expression" `
}

type ASTNode struct {
	Type  NodeType
	Value string // Значение (число или оператор).
	Left  *ASTNode
	Right *ASTNode
}

type Task struct {
	Id             int           `json:"id"`
	Arg1           float64       `json:"arg1"`
	Arg2           float64       `json:"arg2"`
	Operation      string        `json:"operation"`
	Operation_time time.Duration `json:"operation_time"`
}

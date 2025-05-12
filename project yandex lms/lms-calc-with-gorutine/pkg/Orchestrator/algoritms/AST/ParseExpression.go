package AST

import (
	"fmt"
	pb "github.com/bust6k/protoLMS"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"project_yandex_lms/lms-calc-with-gorutine/config"
	"project_yandex_lms/lms-calc-with-gorutine/entites"

	"context"
	"project_yandex_lms/lms-calc-with-gorutine/variables"
	"strconv"
)

func SplitAST(node *entites.ASTNode) []entites.Task {
	if node == nil {
		return nil
	}

	var tasks []entites.Task

	if node.Type == variables.OperatorNode {
		leftTasks := SplitAST(node.Left)
		tasks = append(tasks, leftTasks...)

		rightTasks := SplitAST(node.Right)
		tasks = append(tasks, rightTasks...)

		task := entites.Task{
			Id:        variables.CurrentCountOfUnprocessedUserExpressions,
			Arg1:      parseOperand(node.Left),
			Operation: node.Value,
			Arg2:      parseOperand(node.Right),
		}
		tasks = append(tasks, task)
	}

	return tasks
}
func parseOperand(node *entites.ASTNode) float64 {
	if node.Type == variables.NumberNode {
		value, _ := strconv.ParseFloat(node.Value, 64)
		return value
	} else if node.Type == variables.OperatorNode {
		left := parseOperand(node.Left)
		right := parseOperand(node.Right)
		switch node.Value {
		case "+":
			return left + right
		case "-":
			return left - right
		case "*":
			return left * right
		case "/":
			return left / right
		}
	}
	return 0
}

func collectExpression(node *entites.ASTNode) string {
	if node == nil {
		return ""
	}
	if node.Type == variables.NumberNode {
		return node.Value
	}

	if node.Value == "+" || node.Value == "-" {
		return "(" + collectExpression(node.Left) + " " + node.Value + " " + collectExpression(node.Right) + ")"
	}
	return collectExpression(node.Left) + " " + node.Value + " " + collectExpression(node.Right)
}

func PostTasksToServer(tasks []entites.Task) error {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		return fmt.Errorf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	conn, err := grpc.Dial(config.DefaultGrpcConfig().InternalServ, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("не удалось подключиться к gRPC скрверу с ошибкой:%v", err)
		logger.Warn("ошибка при попытке подключиться к gRPC серверу internalserv ", zap.Error(err))

	}
	defer conn.Close()
	client := pb.NewInternalServiceClient(conn)

	var sliceTaskRequest pb.ListTaskRequest
	for _, task := range tasks {
		TaskReq := pb.TaskRequest{Id: int32(task.Id), Arg1: float32(task.Arg1), Arg2: float32(task.Arg2), Operation: task.Operation, OperationTime: task.Operation_time.String()}
		sliceTaskRequest.List = append(sliceTaskRequest.List, &TaskReq)

	}
	_, err = client.PushTasks(context.Background(), &sliceTaskRequest)
	if err != nil {
		logger.Warn("ошибка при попытке запушить задачи на энд поинт /internal ", zap.Error(err))
		return fmt.Errorf("ошибка при  попытке отправить задачи на /internal с ошибкой: %v", err)
	}

	return nil
}

package taskserv

import (
	"context"
	"fmt"
	pb "github.com/bust6k/protoLMS"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"project_yandex_lms/lms-calc-with-gorutine/entites"

	"project_yandex_lms/lms-calc-with-gorutine/grpc/useful"
	"project_yandex_lms/lms-calc-with-gorutine/variables"
)

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
}

func NewTaskService() *TaskServer {
	return &TaskServer{}
}

func (t *TaskServer) PushTask(ctx context.Context, task *pb.TaskRequest) (*emptypb.Empty, error) {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	taskRequest, err := useful.FromTaskRequestToTask(task)

	if err != nil {
		logger.Warn("при преобразовании типа  TaskRequest в Task произошла ошибка", zap.Error(err))
		return nil, fmt.Errorf("ошибка при преобразовании типа TaskRequest ы task с ошибкой:%v", err)
	}
	logger.Debug("task request", zap.Reflect("task request", taskRequest))
	variables.CurrentTask = taskRequest

	return nil, nil
}

func (t *TaskServer) GetTask(ctx context.Context, em *emptypb.Empty) (*pb.TaskRequest, error) {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)
	}
	currentTask := variables.CurrentTask

	taskRequest := useful.FromTaskToTaskRequest(currentTask)

	variables.CurrentTask = entites.Task{0, 0, 0, "", currentTask.Operation_time}
	return taskRequest, nil
}

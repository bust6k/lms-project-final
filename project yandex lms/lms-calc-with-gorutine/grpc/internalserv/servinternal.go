package internalserv

import (
	"context"
	"fmt"
	pb "github.com/bust6k/protoLMS"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"project_yandex_lms/lms-calc-with-gorutine/config"

	"project_yandex_lms/lms-calc-with-gorutine/grpc/useful"
	"project_yandex_lms/lms-calc-with-gorutine/variables"
)

type InternalServer struct {
	pb.UnimplementedTaskServiceServer
	pb.InternalServiceServer
}

func New() *InternalServer {

	return &InternalServer{}

}

func (i *InternalServer) PushTasks(ctx context.Context, list *pb.ListTaskRequest) (*emptypb.Empty, error) {

	logger, err := zap.NewDevelopment()

	defer logger.Sync()
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	for _, taskElem := range list.List {
		taskElemInWiewTask, err := useful.FromTaskRequestToTask(taskElem)

		if err != nil {
			logger.Warn("ошибка при преобразовании типа TaskRequest в Task ", zap.Error(err))
			return nil, fmt.Errorf("ошибка при преобразовании типа TaskRequest в Task с ошибкой:%v", err)

		}

		variables.TheTasks = append(variables.TheTasks, taskElemInWiewTask)
	}

	for len(variables.TheTasks) > 0 {
		firstElement := variables.TheTasks[0]

		variables.TheTasks = variables.TheTasks[1:]

		conn, err := grpc.Dial(config.DefaultGrpcConfig().TaskServ, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Warn("ошибка при попытке установить соединение с gRPC сервером task", zap.Error(err))
			return nil, fmt.Errorf("не получилось подключиться к gRPC серверу с ошибкой: ", err)

		}
		defer conn.Close()

		logger.Debug("первый элемент", zap.Reflect("first element", firstElement))
		client := pb.NewTaskServiceClient(conn)

		client.PushTask(context.Background(), useful.FromTaskToTaskRequest(firstElement))
	}

	return nil, nil
}

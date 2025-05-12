package main

import (
	"fmt"
	pb "github.com/bust6k/protoLMS"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"project_yandex_lms/lms-calc-with-gorutine/config"

	"project_yandex_lms/lms-calc-with-gorutine/grpc/taskserv"
)

func main() {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()
	if err != nil {
		fmt.Errorf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	lis, err := net.Listen("tcp", config.DefaultGrpcConfig().TaskServ)
	if err != nil {
		logger.Warn("ошибка при создании tcp сервера", zap.Error(err), zap.String("serverName", "TaskServer"))
	}
	defer lis.Close()
	grpcServ := grpc.NewServer()

	TaskServ := taskserv.NewTaskService()

	pb.RegisterTaskServiceServer(grpcServ, TaskServ)
	if err := grpcServ.Serve(lis); err != nil {
		log.Println("ошибка при Серве grpc: ", err)
		os.Exit(1)
	}
}

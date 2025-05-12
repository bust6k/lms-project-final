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
	"project_yandex_lms/lms-calc-with-gorutine/grpc/processedexpserv"
)

func main() {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()
	if err != nil {
		fmt.Errorf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	lis, err := net.Listen("tcp", config.DefaultGrpcConfig().ProcessedExpServ)
	if err != nil {
		logger.Warn("ошибка при создании tcp сервера", zap.Error(err), zap.String("serverName", "processed expression server"))
	}

	grpcServ := grpc.NewServer()

	internalServ := processedexpserv.New()

	pb.RegisterProcessedExpressionsServiceServer(grpcServ, internalServ)
	if err := grpcServ.Serve(lis); err != nil {
		log.Println("ошибка при Серве grpc: ", err)
		os.Exit(1)
	}
}

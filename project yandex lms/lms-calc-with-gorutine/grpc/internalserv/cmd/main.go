package main

import (
	pb "github.com/bust6k/protoLMS"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"project_yandex_lms/lms-calc-with-gorutine/config"
	"project_yandex_lms/lms-calc-with-gorutine/grpc/internalserv"
)

func main() {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()
	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	lis, err := net.Listen("tcp", config.DefaultGrpcConfig().InternalServ)
	if err != nil {
		logger.Warn("ошибка при создании tcp сервера", zap.Error(err), zap.String("serverName", "InternalServer"))
		return
	}

	grpcServ := grpc.NewServer()

	internalServ := internalserv.New()

	pb.RegisterInternalServiceServer(grpcServ, internalServ)
	if err := grpcServ.Serve(lis); err != nil {
		logger.Warn("ошибка при сервер gRPC", zap.Error(err))
	}
}

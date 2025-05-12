package agent

import (
	"context"
	"fmt"
	pb "github.com/bust6k/protoLMS"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"project_yandex_lms/lms-calc-with-gorutine/calc"
	"project_yandex_lms/lms-calc-with-gorutine/config"
	"project_yandex_lms/lms-calc-with-gorutine/entites"
	"strconv"
	"time"
)

func (A *Agent) CalculateTasks() {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	A.Wg.Add(A.Computing_power)
	for i := 0; i < A.Computing_power; i++ {
		go func() {
			defer A.Wg.Done()

			conn, err := grpc.Dial(config.DefaultGrpcConfig().TaskServ, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				logger.Warn("ошибка при подключении к gRPC серверу", zap.Error(err))

			}
			defer conn.Close()
			client := pb.NewTaskServiceClient(conn)

			req, err := client.GetTask(context.Background(), nil)
			if err != nil {
				A.chanErrors <- fmt.Errorf("ошибка при получении задачи с помощью gRPC")
				logger.Debug("error", zap.Error(err))
			}

			ParsedTime, err := time.ParseDuration(req.OperationTime)
			if err != nil {
				A.chanErrors <- fmt.Errorf("ошибка при парсинге string в time.Duration:%v", err)
			}

			task := entites.Task{Id: int(req.Id), Arg1: float64(req.Arg1), Arg2: float64(req.Arg2), Operation: req.Operation, Operation_time: ParsedTime}

			if err != nil {
				A.chanErrors <- fmt.Errorf("ошибка при получении новой задачи: %v", err)
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), task.Operation_time)
			defer cancel()

			select {
			case <-ctx.Done():
				A.chanErrors <- fmt.Errorf("ошибка, горутина слишком долго выполняет свою работу")
				return
			default:
				A.chanTasks <- task

				expression := strconv.FormatFloat(task.Arg1, 'f', -1, 64) + task.Operation + strconv.FormatFloat(task.Arg2, 'f', -1, 64)

				result, err := calc.Calc(expression)
				if err != nil {
					A.chanErrors <- fmt.Errorf("ошибка при расчете выражения: %v", err)
					return
				}

				A.chanResults <- result

				logger.Debug("вот один из результатов вычислений  и то, что он вычислял", zap.Float64("result", result), zap.Reflect("expression", expression))

			}
		}()
	}
	go func() {
		A.Wg.Wait()
		close(A.chanTasks)
		close(A.chanResults)
		close(A.chanErrors)
	}()

}

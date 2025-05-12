package agent

import (
	"context"
	"fmt"
	pb "github.com/bust6k/protoLMS"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"project_yandex_lms/lms-calc-with-gorutine/config"
	"sync"
)

func (A *Agent) SumTasksResultsInResultsChan() (float64, error) {
	var stackMutex sync.Mutex
	var stack []float64

	for taksi := range A.chanTasks {
		result := <-A.chanResults

		stackMutex.Lock()

		switch taksi.Operation {
		case "*", "/":
			if len(stack) > 0 {
				last := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if taksi.Operation == "*" {
					stack = append(stack, last*result)
				} else {
					if result == 0 {
						A.chanErrors <- fmt.Errorf("деление на ноль")
						stackMutex.Unlock()
						continue
					}
					stack = append(stack, last/result)
				}
			} else {
				stack = append(stack, result)
			}
		case "+", "-":
			stack = append(stack, result)
		}

		stackMutex.Unlock()
	}
	if len(stack) == 0 {
		return 0, fmt.Errorf("ошибка стэк пуст, результат не вычислен")
	}
	finalresult := stack[len(stack)-1]
	return finalresult, nil
}

func (A *Agent) PrintErrorsInChanErrors() {
	for err := range A.chanErrors {
		logger, _ := zap.NewDevelopment()

		defer logger.Sync()

		logger.Error("ошибка", zap.Error(err))
	}

}

func (A *Agent) PostTaskResultsToServer(userId string) error {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	finalresult, err := A.SumTasksResultsInResultsChan()
	if err != nil {
		logger.Warn("ошибка при вычислении результата задач", zap.Error(err))
		return fmt.Errorf("ошибка при вычислении результата задач")

	}
	logger.Debug("вот финальный результат", zap.Float64("final result", finalresult))

	A.PrintErrorsInChanErrors()

	conn, err := grpc.Dial(config.DefaultGrpcConfig().ProcessedExpServ, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Warn("ошибка при подключении у gRPC серверу", zap.Error(err))
		return fmt.Errorf("не получилось подключиться к серверу: ", err)

	}
	defer conn.Close()
	client := pb.NewProcessedExpressionsServiceClient(conn)
	_, err = client.PushNewProcessedExpression(context.Background(), &pb.ResultRequest{Result: float32(finalresult), UserId: userId})
	if err != nil {
		logger.Warn("ошибка при отправке выражения на сервер", zap.Error(err))
		return fmt.Errorf("ошибка при попытке отправить выражения на сервер:%v", err)
	}
	return nil

}

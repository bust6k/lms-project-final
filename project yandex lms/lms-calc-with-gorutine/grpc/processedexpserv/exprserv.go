package processedexpserv

import (
	"context"
	"fmt"
	pb "github.com/bust6k/protoLMS"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"project_yandex_lms/lms-calc-with-gorutine/database"
	"project_yandex_lms/lms-calc-with-gorutine/models"
	"project_yandex_lms/lms-calc-with-gorutine/variables"
)

type ExprServer struct {
	pb.UnimplementedProcessedExpressionsServiceServer
}

func New() *ExprServer {
	return &ExprServer{}
}

func (e *ExprServer) PushNewProcessedExpression(ctx context.Context, res *pb.ResultRequest) (*emptypb.Empty, error) {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	var expression models.ProcessedExpression
	expression.Id = variables.CurrentCountOfUnprocessedUserExpressions
	expression.Status = "ready"
	expression.Result = float64(res.Result)
	expression.UserId = res.UserId
	logger.Debug("processedExpression", zap.Reflect("processed expression", expression))

	if err := database.SaveProcessedExprToDB(expression); err != nil {
		return nil, fmt.Errorf("ошибка при сохранении уже обработанного пользовательского выражения в  базу данных:%v", err)
	}

	return nil, nil
}

func (e *ExprServer) GetAllProcessedExpressions(ctx context.Context, ec *emptypb.Empty) (*pb.ListProcesedExpressionRequest, error) {
	var sliceOfProcessedExpressionsRequest pb.ListProcesedExpressionRequest
	processedexprs, err := database.GetProcessedExprsInDB()

	if err != nil {
		return nil, err
	}

	for _, expr := range processedexprs {
		var newExpr pb.ProcessedExpressionRequest
		newExpr.Status = expr.Status
		newExpr.Id = int32(expr.Id)
		newExpr.Result = float32(expr.Result)
		newExpr.UserId = expr.UserId

		sliceOfProcessedExpressionsRequest.List = append(sliceOfProcessedExpressionsRequest.List, &newExpr)
	}

	return &sliceOfProcessedExpressionsRequest, nil
}

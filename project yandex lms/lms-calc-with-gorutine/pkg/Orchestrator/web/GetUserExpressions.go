package web

import (
	"context"
	pb "github.com/bust6k/protoLMS"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"project_yandex_lms/lms-calc-with-gorutine/config"
	"project_yandex_lms/lms-calc-with-gorutine/models"
)



var (
	grpcConn *grpc.ClientConn
)

func init() {
	var err error
	grpcConn, err = grpc.Dial(
		config.DefaultGrpcConfig().ProcessedExpServ,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic("failed to connect to gRPC server")
	}
}

func GetUserExpressions(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodPost:
		handlePostExpression(c)
	case http.MethodGet:
		handleGetExpressions(c)
	default:
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	}
}

func handlePostExpression(c *gin.Context) {
	userId := c.MustGet("user_id").(string)
	var expr models.ProcessedExpression
	if err := c.ShouldBindJSON(&expr); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	client := pb.NewProcessedExpressionsServiceClient(grpcConn)
	_, err := client.PushNewProcessedExpression(context.Background(), &pb.ResultRequest{
		Result: float32(expr.Result), UserId: userId,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "service error"})
		return
	}

	c.Status(http.StatusOK)
}

func handleGetExpressions(c *gin.Context) {
	userId := c.MustGet("user_id").(string)
	client := pb.NewProcessedExpressionsServiceClient(grpcConn)

	exprs, err := client.GetAllProcessedExpressions(context.Background(), nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "service error"})
		return
	}

	result := make([]models.ProcessedExpression, 0)
	for _, expr := range exprs.List {

		if expr.UserId == userId {
			result = append(result, models.ProcessedExpression{
				Id:     int(expr.Id),
				Status: expr.Status,
				Result: float64(expr.Result),
				UserId: userId,
			})
		}
	}

	c.JSON(http.StatusOK, result)
}

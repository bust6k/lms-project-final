package processedexpserv_test

import (
	"context"
	"project_yandex_lms/lms-calc-with-gorutine/grpc/processedexpserv"

	"fmt"
	pb "github.com/bust6k/protoLMS"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"

	"project_yandex_lms/lms-calc-with-gorutine/models"

	"testing"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) SaveProcessedExprToDB(expr models.ProcessedExpression) error {
	args := m.Called(expr)
	return args.Error(0)
}

func (m *MockDB) GetProcessedExprsInDB() ([]models.ProcessedExpression, error) {
	args := m.Called()
	return args.Get(0).([]models.ProcessedExpression), args.Error(1)
}

func setupGRPCServer(db *MockDB) (*grpc.Server, *bufconn.Listener) {
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	exprServer := processedexpserv.New()
	pb.RegisterProcessedExpressionsServiceServer(srv, exprServer)
	go func() {

		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return srv, lis
}

func bufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestExprServer_PushNewProcessedExpression(t *testing.T) {
	mockDB := new(MockDB)
	srv, lis := setupGRPCServer(mockDB)
	defer srv.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewProcessedExpressionsServiceClient(conn)

	t.Run("successful save", func(t *testing.T) {
		mockDB.On("SaveProcessedExprToDB", mock.Anything).Return(nil)

		_, err := client.PushNewProcessedExpression(ctx, &pb.ResultRequest{
			Result: 42.0,
			UserId: "user123",
		})

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		mockDB.On("SaveProcessedExprToDB", mock.Anything).Return(fmt.Errorf("db error"))

		_, err := client.PushNewProcessedExpression(ctx, &pb.ResultRequest{
			Result: 42.0,
			UserId: "user123",
		})

		assert.Error(t, err)
		mockDB.AssertExpectations(t)
	})
}

func TestExprServer_GetAllProcessedExpressions(t *testing.T) {
	mockDB := new(MockDB)
	srv, lis := setupGRPCServer(mockDB)
	defer srv.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewProcessedExpressionsServiceClient(conn)

	t.Run("successful get", func(t *testing.T) {
		testData := []models.ProcessedExpression{
			{Id: 1, Status: "ready", Result: 10.5, UserId: "user1"},
			{Id: 2, Status: "error", Result: 0, UserId: "user2"},
		}
		mockDB.On("GetProcessedExprsInDB").Return(testData, nil)

		resp, err := client.GetAllProcessedExpressions(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Len(t, resp.List, 2)
		mockDB.AssertExpectations(t)
	})

	t.Run("empty result", func(t *testing.T) {
		mockDB.On("GetProcessedExprsInDB").Return([]models.ProcessedExpression{}, nil)

		resp, err := client.GetAllProcessedExpressions(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Empty(t, resp.List)
		mockDB.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		mockDB.On("GetProcessedExprsInDB").Return([]models.ProcessedExpression{}, fmt.Errorf("db error"))

		_, err := client.GetAllProcessedExpressions(ctx, &emptypb.Empty{})
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
	})
}

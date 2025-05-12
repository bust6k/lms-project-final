package web

import (
	"context"
	pb "github.com/bust6k/protoLMS"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockProcessedExpressionsClient struct {
	mock.Mock
	pb.ProcessedExpressionsServiceClient
}

func (m *mockProcessedExpressionsClient) PushNewProcessedExpression(ctx context.Context, in *pb.ResultRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(in)
	return &empty.Empty{}, args.Error(0)
}

func (m *mockProcessedExpressionsClient) GetAllProcessedExpressions(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*pb.ListProcesedExpressionRequest, error) {
	args := m.Called()
	return args.Get(0).(*pb.ListProcesedExpressionRequest), args.Error(1)
}

func TestHandlePostExpression(t *testing.T) {
	originalConn := grpcConn
	defer func() { grpcConn = originalConn }()

	mockConn, _ := grpc.Dial("", grpc.WithInsecure())
	grpcConn = mockConn

	tests := []struct {
		name           string
		inputData      string
		expectedStatus int
		mockSetup      func(*mockProcessedExpressionsClient)
	}{
		{
			name:           "successful expression post",
			inputData:      `{"Id":1,"Status":"ready","Result":5.5,"UserId":"user1"}`,
			expectedStatus: http.StatusOK,
			mockSetup: func(m *mockProcessedExpressionsClient) {
				m.On("PushNewProcessedExpression", &pb.ResultRequest{
					Result: 5.5,
					UserId: "user1",
				}).Return(nil)
			},
		},
		{
			name:           "invalid json",
			inputData:      `invalid json`,
			expectedStatus: http.StatusBadRequest,
			mockSetup:      func(m *mockProcessedExpressionsClient) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mockProcessedExpressionsClient)
			tt.mockSetup(mockClient)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/", strings.NewReader(tt.inputData))
			c.Set("user_id", "user1")

			handlePostExpression(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestHandleGetExpressions(t *testing.T) {
	originalConn := grpcConn
	defer func() { grpcConn = originalConn }()

	mockConn, _ := grpc.Dial("", grpc.WithInsecure())
	grpcConn = mockConn

	mockClient := new(mockProcessedExpressionsClient)
	mockClient.On("GetAllProcessedExpressions").Return(&pb.ListProcesedExpressionRequest{
		List: []*pb.ProcessedExpressionRequest{
			{Id: 1, Status: "ready", Result: 5.5, UserId: "user1"},
		},
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("user_id", "user1")

	handleGetExpressions(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	mockClient.AssertExpectations(t)
}

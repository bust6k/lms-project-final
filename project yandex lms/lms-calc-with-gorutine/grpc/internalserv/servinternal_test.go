package internalserv

import (
	"context"
	pb "github.com/bust6k/protoLMS"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var internalSrv = New()

func init() {
	lis = bufconn.Listen(bufSize)
	srv := grpc.NewServer()

	pb.RegisterInternalServiceServer(srv, internalSrv)
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestPushTasks(t *testing.T) {
	tests := []struct {
		name    string
		list    *pb.ListTaskRequest
		wantErr bool
	}{
		{
			name: "valid request",
			list: &pb.ListTaskRequest{
				List: []*pb.TaskRequest{
					{
						Id:            1,
						Arg1:          10.5,
						Arg2:          2.5,
						Operation:     "+",
						OperationTime: "10000",
					},
					{
						Id:            2,
						Arg1:          10.5,
						Arg2:          7.5,
						Operation:     "+",
						OperationTime: "10000",
					},
					{
						Id:            3,
						Arg1:          10.5,
						Arg2:          10.5,
						Operation:     "-",
						OperationTime: "10000",
					},
				},
			},
			wantErr: false,
		}, {
			name: "invalid operation",
			list: &pb.ListTaskRequest{
				List: []*pb.TaskRequest{
					{
						Id:            1,
						Arg1:          10.5,
						Arg2:          2.5,
						Operation:     "+",
						OperationTime: "10000",
					},
					{
						Id:            2,
						Arg1:          10.5,
						Arg2:          7.5,
						Operation:     "%",
						OperationTime: "10000",
					},
					{
						Id:            3,
						Arg1:          10.5,
						Arg2:          10.5,
						Operation:     "^",
						OperationTime: "10000",
					},
				},
			},
			wantErr: true,
		}, {
			name: "invalid request",
			list: &pb.ListTaskRequest{
				List: []*pb.TaskRequest{
					{
						Id:            111110,
						Arg1:          -0000.5,
						Arg2:          2.5,
						Operation:     "add",
						OperationTime: "10000",
					},
					{
						Id:            2,
						Arg1:          10.5,
						Arg2:          7.5,
						Operation:     "+",
						OperationTime: "10000",
					},
					{
						Id:            3,
						Arg1:          10.5,
						Arg2:          10.5,
						Operation:     "-",
						OperationTime: "1skib000",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())

			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()

			client := pb.NewInternalServiceClient(conn)

			_, err = client.PushTasks(ctx, tt.list)

			if err != nil && !tt.wantErr || err == nil && tt.wantErr {
				t.Errorf("PushTasks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

.PHONY: all run
all : run

run:
	go run lms-calc-with-gorutine/grpc/internalserv/cmd/main.go &
	go run lms-calc-with-gorutine/grpc/taskserv/cmd/main.go &
	go run lms-calc-with-gorutine/grpc/processedexpserv/cmd/main.go &
	go run main.go
package useful

import (
	pb "github.com/bust6k/protoLMS"
	"project_yandex_lms/lms-calc-with-gorutine/entites"
)

func FromTaskToTaskRequest(task entites.Task) *pb.TaskRequest {

	taskReqInWiewTask := pb.TaskRequest{Id: int32(task.Id), Arg1: float32(task.Arg1), Arg2: float32(task.Arg2), Operation: task.Operation, OperationTime: task.Operation_time.String()}
	return &taskReqInWiewTask

}

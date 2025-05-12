package useful

import (
	"fmt"
	pb "github.com/bust6k/protoLMS"
	"project_yandex_lms/lms-calc-with-gorutine/entites"
	"time"
)

func FromTaskRequestToTask(task *pb.TaskRequest) (entites.Task, error) {

	dur, err := time.ParseDuration(task.OperationTime)
	if err != nil {
		return entites.Task{}, fmt.Errorf("ошибка при парсинге")
	}
	taskReqInWiewTask := entites.Task{Id: int(task.Id), Arg1: float64(task.Arg1), Arg2: float64(task.Arg2), Operation: task.Operation, Operation_time: dur}
	return taskReqInWiewTask, nil
}

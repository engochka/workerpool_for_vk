package workerpool

import (
	"fmt"
	"time"
)

type Task struct {
	Data   string
	TaskID int
}

// NewTask возвращает экземпляр Task
func NewTask(data string, taskID int) (*Task, error) {
	return &Task{Data: data, TaskID: taskID}, nil
}

// process выполняет задачу
func process(workerID int, task *Task) {
	time.Sleep(500 * time.Millisecond) // имитация бурной деятельности
	fmt.Printf("Worker %d выполнил задачу #%d:%v\n", workerID, task.TaskID, task.Data)
}

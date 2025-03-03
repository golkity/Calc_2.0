package store

import (
	"sync"

	"lms/internal/task"
)

var (
	DB = &memoryStore{
		Expressions: map[string]*task.Expression{},
		Tasks:       map[int]*task.Task{},
	}
	StoreMu    sync.Mutex
	NextExprID int
	NextTaskID int
)

type memoryStore struct {
	Expressions map[string]*task.Expression
	Tasks       map[int]*task.Task
}

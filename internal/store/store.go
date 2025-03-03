package store

import (
	"sync"

	"github.com/golkity/Calc_2.0/internal/task"
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

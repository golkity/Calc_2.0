package orchestrator

import (
	"strconv"

	"github.com/golkity/Calc_2.0/internal/custom_errors"
	"github.com/golkity/Calc_2.0/internal/store"
	"github.com/golkity/Calc_2.0/internal/task"
)

func AddNewExpression(expr string) (string, error) {
	store.StoreMu.Lock()
	defer store.StoreMu.Unlock()

	store.NextExprID++
	id := strconv.Itoa(store.NextExprID)

	e := &task.Expression{
		ID:     id,
		Status: "pending",
		Result: nil,
	}
	store.DB.Expressions[id] = e

	store.NextTaskID++
	tid := store.NextTaskID
	t := &task.Task{
		ID:            tid,
		Arg1:          expr,
		Arg2:          "",
		Operation:     "",
		OperationTime: 0,
		Status:        task.StatusPending,
		ExpressionID:  id,
	}
	store.DB.Tasks[tid] = t

	e.Tasks = append(e.Tasks, tid)

	return id, nil
}

func GetAllExpressions() []*task.Expression {
	store.StoreMu.Lock()
	defer store.StoreMu.Unlock()

	var out []*task.Expression
	for _, e := range store.DB.Expressions {
		out = append(out, e)
	}
	return out
}

func GetExpressionByID(id string) (*task.Expression, error) {
	store.StoreMu.Lock()
	defer store.StoreMu.Unlock()

	e, ok := store.DB.Expressions[id]
	if !ok {
		return nil, custom_errors.ErrNotFound
	}
	return e, nil
}

func GetNextTask() *task.Task {
	store.StoreMu.Lock()
	defer store.StoreMu.Unlock()

	for _, t := range store.DB.Tasks {
		if t.Status != task.StatusCompleted {
			return t
		}
	}
	return nil
}

func CompleteTask(taskID int, result float64) error {
	store.StoreMu.Lock()
	defer store.StoreMu.Unlock()

	t, ok := store.DB.Tasks[taskID]
	if !ok {
		return custom_errors.ErrNotFound
	}
	if t.Status == task.StatusCompleted {
		return nil
	}
	t.Status = task.StatusCompleted

	t.Result = &result

	exprID := t.ExpressionID
	e, ok2 := store.DB.Expressions[exprID]
	if !ok2 {
		// Теоретически если это будет, то я буду плакать
		return nil
	}

	allDone := true
	var lastVal *float64
	for _, tid := range e.Tasks {
		t := store.DB.Tasks[tid]
		if t.Status != task.StatusCompleted || t.Result == nil {
			allDone = false
			break
		}
		lastVal = t.Result
	}
	if allDone {
		e.Status = "done"
		e.Result = lastVal
	}
	return nil
}

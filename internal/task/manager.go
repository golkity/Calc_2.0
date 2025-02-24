package task

import (
	"fmt"
	"sync"
	"time"
)

var taskIDGenerator sync.Mutex

func GenerateID() string {
	taskIDGenerator.Lock()
	defer taskIDGenerator.Unlock()
	return fmt.Sprintf("task-%d", time.Now().UnixNano())
}

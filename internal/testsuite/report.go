package testsuite

import (
	_ "embed"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

var goose = `░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░░ЗАПУСКАЕМ░ГУСЕЙ-РАЗВЕДЧИКОВ░░░░
░░░░░▄▀▀▀▄░░░▄▀▀▀▀▄░░░▄▀▀▀▄░░░░░
▄███▀░◐░░░▌░▐0░░░░0▌░▐░░░◐░▀███▄
░░░░▌░░░░░▐░▌░▐▀▀▌░▐░▌░░░░░▐░░░░
░░░░▐░░░░░▐░▌░▌▒▒▐░▐░▌░░░░░▌░░░░`

type Result struct {
	name    string
	ok      bool
	elapsed time.Duration
}

var (
	mu   sync.Mutex
	logs []Result
)

func Record(name string, ok bool, d time.Duration) {
	mu.Lock()
	logs = append(logs, Result{name, ok, d})
	mu.Unlock()
}

func Spinner(d time.Duration) {
	frames := []rune{'|', '/', '-', '\\'}
	next := time.NewTicker(120 * time.Millisecond)
	defer next.Stop()

	deadline := time.Now().Add(d)
	i := 0
	for time.Now().Before(deadline) {
		fmt.Printf("\rЗагрузка %c", frames[i%len(frames)])
		i++
		<-next.C
	}
	fmt.Print("\rЗагрузка завершена\n\n")
}

func Report() {
	if len(logs) == 0 {
		return
	}

	sort.Slice(logs, func(i, j int) bool { return logs[i].name < logs[j].name })

	nameW := len("TEST NAME")
	timeW := len("TIME, ms")
	for _, r := range logs {
		if l := len(r.name); l > nameW {
			nameW = l
		}
	}
	border := func(left, fill, sep, right string) {
		fmt.Print(left, strings.Repeat(fill, nameW+2), sep, strings.Repeat(fill, 8), sep,
			strings.Repeat(fill, timeW+2), right, "\n")
	}

	border("╔", "═", "╦", "╗")
	fmt.Printf("║ %-*s ║ %-6s ║ %-*s ║\n", nameW, "TEST NAME", "STATUS", timeW, "TIME, ms")
	border("╠", "═", "╬", "╣")

	var pass int
	for _, r := range logs {
		status := "PASS"
		if !r.ok {
			status = "FAIL"
		} else {
			pass++
		}
		fmt.Printf("║ %-*s ║ %-6s ║ %*d ║\n",
			nameW, r.name, status, timeW, r.elapsed.Milliseconds())
	}
	border("╚", "═", "╩", "╝")

	pct := float64(pass) / float64(len(logs)) * 100
	bar := strings.Repeat("█", int(pct/5)) + strings.Repeat("░", 20-int(pct/5))
	fmt.Printf("\nШкала: [%s] (%.2f%%).  Итог: %d/%d\n\n", bar, pct, pass, len(logs))
}

func TestWarps(m *testing.M) {
	fmt.Println(goose)

	Spinner(5 * time.Second)

	code := m.Run()

	Report()
	os.Exit(code)
}

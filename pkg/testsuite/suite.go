package testsuite

import (
	"os"
	"strconv"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

type Suite[T any] struct {
	Tests []T `yaml:"tests"`
}

func Load[T any](path string) (*Suite[T], error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s Suite[T]
	if err := yaml.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func AlmostEq(got, want, eps float64) bool {
	d := got - want
	return d < eps && d > -eps
}

func Run[T any](t *testing.T, cases []T, test func(tc T, t *testing.T)) {
	for i, c := range cases {
		tc := c
		name := t.Name() + "#" + strconv.Itoa(i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			start := time.Now()
			defer func() { Record(name, !t.Failed(), time.Since(start)) }()
			test(tc, t)
		})
	}
}

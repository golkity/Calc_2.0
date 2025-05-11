package calc

import (
	"testing"

	"github.com/golkity/Calc_2.0/testsuite"
)

type Case struct {
	Expression string  `yaml:"expression"`
	Expected   float64 `yaml:"expected"`
}

func TestMain(m *testing.M) { testsuite.WrapMain(m) }

func TestCalc(t *testing.T) {
	s, err := testsuite.Load[Case]("../../test/calc_test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	testsuite.Run(t, s.Tests, func(tc Case, t *testing.T) {
		got, err := Calc(tc.Expression)
		if err != nil {
			t.Fatalf("Calc(%q) error: %v", tc.Expression, err)
		}
		if !testsuite.AlmostEq(got, tc.Expected, 1e-4) {
			t.Errorf("want %v, got %v", tc.Expected, got)
		}
	})
}

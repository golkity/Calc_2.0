package stuct

type CalcTestCase struct {
	Expression string  `yaml:"expression"`
	Expected   float64 `yaml:"expected"`
}

type CalcTestSuite struct {
	Tests []CalcTestCase `yaml:"tests"`
}

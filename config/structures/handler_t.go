package stuct

type HandlerTestCase struct {
	Name   string `yaml:"name"`
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
	Body   string `yaml:"body"`
	Status int    `yaml:"status"`
}

type HandlerTestCases struct {
	Tests []HandlerTestCase `yaml:"tests"`
}

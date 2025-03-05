package stuct

type ConfigTestCases struct {
	Tests []ConfigTestCase `yaml:"tests"`
}

type ConfigTestCase struct {
	Name          string          `yaml:"name"`
	FileContent   *string         `yaml:"fileContent"`
	Expected      *ExpectedConfig `yaml:"expected,omitempty"`
	ExpectedError string          `yaml:"expectedError,omitempty"`
}

type ExpectedConfig struct {
	ServerPort string `yaml:"ServerPort"`
	LogFile    string `yaml:"LogFile"`
}

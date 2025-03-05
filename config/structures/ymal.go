package stuct

import (
	"os"

	"gopkg.in/yaml.v2"
)

func LoadYML[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result T
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

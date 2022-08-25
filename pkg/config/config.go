package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	URL    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Data   string            `json:"data"`
}

type Config struct {
	BaseURL string
}

func ReadYAML(fileName string, data interface{}) error {
	file, err := os.ReadFile(fileName)

	if err == nil {
		err = yaml.Unmarshal(file, &data)
	}

	return err
}

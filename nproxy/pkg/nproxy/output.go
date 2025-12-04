package nproxy

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type ErrorOutput struct {
	Error ErrorDetail `yaml:"error"`
}

type ErrorDetail struct {
	Code    string `yaml:"code"`
	Message string `yaml:"message"`
}

func PrintYAML(data interface{}) error {
	return writeYAML(os.Stdout, data)
}

func writeYAML(w io.Writer, data interface{}) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return APIError(fmt.Sprintf("failed to serialize: %s", err))
	}
	return enc.Close()
}

func PrintError(err error) {
	ne, ok := err.(*NproxyError)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}

	output := ErrorOutput{
		Error: ErrorDetail{
			Code:    string(ne.Code),
			Message: ne.Message,
		},
	}

	if err := writeYAML(os.Stderr, output); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", ne.Message)
	}
}

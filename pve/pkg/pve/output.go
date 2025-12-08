package pve

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

func PrintYAMLTo(w io.Writer, data interface{}) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return APIError(fmt.Sprintf("failed to serialize: %s", err))
	}
	return nil
}

func PrintYAML(data interface{}) error {
	return PrintYAMLTo(os.Stdout, data)
}

func PrintError(err error) {
	pe, ok := err.(*PveError)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}

	output := ErrorOutput{
		Error: ErrorDetail{
			Code:    string(pe.Code),
			Message: pe.Message,
		},
	}

	enc := yaml.NewEncoder(os.Stderr)
	enc.SetIndent(2)
	if err := enc.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", pe.Message)
	}
}

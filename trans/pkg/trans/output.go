package trans

import (
	"fmt"
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
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return APIError(fmt.Sprintf("failed to serialize: %s", err))
	}
	return nil
}

func PrintError(err error) {
	te, ok := err.(*TransError)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}

	output := ErrorOutput{
		Error: ErrorDetail{
			Code:    string(te.Code),
			Message: te.Message,
		},
	}

	enc := yaml.NewEncoder(os.Stderr)
	enc.SetIndent(2)
	if err := enc.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", te.Message)
	}
}

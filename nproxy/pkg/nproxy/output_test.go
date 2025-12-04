package nproxy

import "testing"

func TestErrorOutputFormat(t *testing.T) {
	err := ConfigError("test message")
	output := ErrorOutput{
		Error: ErrorDetail{
			Code:    string(err.Code),
			Message: err.Message,
		},
	}

	if output.Error.Code != "CONFIG_ERROR" {
		t.Errorf("Code = %s, want CONFIG_ERROR", output.Error.Code)
	}
	if output.Error.Message != "test message" {
		t.Errorf("Message = %s, want 'test message'", output.Error.Message)
	}
}

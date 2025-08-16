package tools

import (
	"fmt"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestToolError(t *testing.T) {
	err := fmt.Errorf("test error")
	result := toolError("Test operation failed", err)

	if result == nil {
		t.Fatal("Expected result to not be nil")
	}

	if !result.IsError {
		t.Error("Expected IsError to be true")
	}

	if len(result.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(result.Content))
	}

	content, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Error("Expected TextContent")
	}

	expected := "Test operation failed: test error"
	if content.Text != expected {
		t.Errorf("Expected text '%s', got '%s'", expected, content.Text)
	}
}

func TestToolSuccess(t *testing.T) {
	data := "test success data"
	result := toolSuccess(data)

	if result == nil {
		t.Fatal("Expected result to not be nil")
	}

	if result.IsError {
		t.Error("Expected IsError to be false")
	}

	if len(result.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(result.Content))
	}

	content, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Error("Expected TextContent")
	}

	if content.Text != data {
		t.Errorf("Expected text '%s', got '%s'", data, content.Text)
	}
}

func TestValidateDeviceID_Valid(t *testing.T) {
	validIDs := []string{
		"abc123",
		"device-123",
		"node_456",
		"test.device",
	}

	for _, id := range validIDs {
		t.Run(id, func(t *testing.T) {
			err := validateDeviceID(id)
			if err != nil {
				t.Errorf("Expected no error for valid ID '%s', got %v", id, err)
			}
		})
	}
}

func TestValidateDeviceID_Invalid(t *testing.T) {
	testCases := []struct {
		id          string
		expectedErr string
	}{
		{"", "device ID cannot be empty"},
		{"ab", "device ID length should be between 3 and 50 characters"},
		{string(make([]byte, 51)), "device ID length should be between 3 and 50 characters"},
		{"abc def", "device ID contains invalid characters"},
		{"abc\tdef", "device ID contains invalid characters"},
		{"abc\ndef", "device ID contains invalid characters"},
		{"abc\rdef", "device ID contains invalid characters"},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			err := validateDeviceID(tc.id)
			if err == nil {
				t.Errorf("Expected error for invalid ID '%s'", tc.id)
			}
			if err.Error() != tc.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tc.expectedErr, err.Error())
			}
		})
	}
}

func TestGetStringParam_Success(t *testing.T) {
	params := map[string]any{
		"deviceID": "test-device",
		"name":     "test-name",
	}

	value, err := getStringParam(params, "deviceID")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if value != "test-device" {
		t.Errorf("Expected 'test-device', got '%s'", value)
	}
}

func TestGetStringParam_Missing(t *testing.T) {
	params := map[string]any{
		"name": "test-name",
	}

	_, err := getStringParam(params, "deviceID")
	if err == nil {
		t.Error("Expected error for missing parameter")
	}

	expected := "deviceID parameter is required"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}

func TestGetStringParam_WrongType(t *testing.T) {
	params := map[string]any{
		"deviceID": 123,
	}

	_, err := getStringParam(params, "deviceID")
	if err == nil {
		t.Error("Expected error for wrong parameter type")
	}

	expected := "deviceID parameter must be a string"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}

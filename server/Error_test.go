package server

import (
	"encoding/json"
	"testing"
)

func TestError(t *testing.T) {
	str, _ := json.Marshal(Error{"code", "description"})

	expected := `{"code":"code","description":"description"}`
	if string(str) != expected {
		t.Errorf("Error marshaled incorrectly: got %v want %v",
			string(str), expected)
	}
}

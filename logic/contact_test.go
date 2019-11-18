package logic

import (
	"fmt"
	"testing"
)

func TestParseContact(t *testing.T) {
	t.Run("invalid json not parsed", func(t *testing.T) {
		data := "INVALID JSON"
		c, err := ParseContact([]byte(data))
		if err == nil {
			t.Errorf("Invalid JSON is parsed")
		}
		if len(c.ID) > 0 {
			t.Errorf("Invalid JSON is parsed to unempty data")
		}
	})

	t.Run("correct json is parsed", func(t *testing.T) {
		id := "test_id"
		data := fmt.Sprintf("{\"contact_id\": \"%v\"}", id)
		c, err := ParseContact([]byte(data))
		if err != nil {
			t.Errorf("Failed to parse with error: %v", err)
		}
		if c.ID != id {
			t.Errorf("ID is invalid. Expected: %v. Got: %v", id, c.ID)
		}
	})
}

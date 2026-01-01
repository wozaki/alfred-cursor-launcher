package internal

import (
	"encoding/json"
	"testing"
)

func TestAddItem(t *testing.T) {
	sf := NewScriptFilter()
	item := Item{
		Title:    "Test",
		Subtitle: "Test subtitle",
		Arg:      "test-arg",
	}

	sf.AddItem(item)

	if len(sf.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(sf.Items))
	}

	if sf.Items[0].Title != "Test" {
		t.Errorf("Expected title 'Test', got '%s'", sf.Items[0].Title)
	}
}

func TestAddErrorItem(t *testing.T) {
	sf := NewScriptFilter()
	sf.AddErrorItem("Test error")

	if len(sf.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(sf.Items))
	}

	if sf.Items[0].Title != "Error: Test error" {
		t.Errorf("Expected title 'Error: Test error', got '%s'", sf.Items[0].Title)
	}

	if sf.Items[0].Valid == nil || *sf.Items[0].Valid != false {
		t.Error("Expected Valid to be false")
	}
}

func TestJSONSerialization(t *testing.T) {
	sf := NewScriptFilter()
	sf.AddItem(Item{
		Title:    "Test",
		Subtitle: "Test subtitle",
		Arg:      "test-arg",
	})

	// Verify ScriptFilter can be serialized to valid JSON
	data, err := json.Marshal(sf)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	var parsed ScriptFilter
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	if len(parsed.Items) != 1 {
		t.Errorf("Expected 1 item in parsed JSON, got %d", len(parsed.Items))
	}
}

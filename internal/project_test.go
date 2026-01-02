package internal

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestProject_ToAlfredItem_Local(t *testing.T) {
	project := Project{
		FolderURI: "file:///path/to/test-project",
		Label:     "test-project",
	}
	result := project.ToAlfredItem()

	if result.Title != "📁 test-project" {
		t.Errorf("Expected title '📁 test-project', got '%v'", result.Title)
	}

	if result.Arg != "file:///path/to/test-project" {
		t.Errorf("Expected arg 'file:///path/to/test-project', got '%v'", result.Arg)
	}
}

func TestProject_ToAlfredItem_Remote(t *testing.T) {
	project := Project{
		FolderURI: "vscode-remote://test-remote/path/to/remote-project",
		Label:     "",
	}
	result := project.ToAlfredItem()

	if result.Title == "" {
		t.Error("Expected title to be set")
	}

	if result.Arg == "" {
		t.Error("Expected arg to be set")
	}
}

func TestProject_ToAlfredItem_Empty(t *testing.T) {
	project := Project{
		FolderURI: "",
		Label:     "",
	}
	result := project.ToAlfredItem()

	if result.Title != "" {
		t.Errorf("Expected empty title, got '%v'", result.Title)
	}
}

func TestRemoveSuffixPattern(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"project-name in project-name-0 (undefined)", "project-name"},
		{"normal project name", "normal project name"},
		{"project in test-123 (undefined) extra", "project extra"},
	}

	for _, tt := range tests {
		result := removeSuffixPattern(tt.input)
		if result != tt.expected {
			t.Errorf("removeSuffixPattern(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"~/test", filepath.Join(homeDir, "test")},
		{"/absolute/path", "/absolute/path"},
	}

	for _, tt := range tests {
		result, err := expandPath(tt.input)
		if err != nil {
			t.Errorf("expandPath(%q) returned error: %v", tt.input, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("expandPath(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestProjectStore_List requires a test database
// This is a placeholder - actual test would require creating a test database
func TestProjectStore_List_WithTestDB(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment (requires Cursor app)")
	}

	// Create a temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "state.vscdb")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ItemTable (
			key TEXT PRIMARY KEY,
			value TEXT
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert test data
	testData := map[string]interface{}{
		"entries": []map[string]interface{}{
			{
				"folderUri": "file:///test/project",
				"label":     "test-project",
			},
		},
	}
	jsonData, _ := json.Marshal(testData)
	_, err = db.Exec(
		"INSERT INTO ItemTable (key, value) VALUES (?, ?)",
		"history.recentlyOpenedPathsList",
		string(jsonData),
	)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Test ProjectStore.List() with the test database
	store := &ProjectStore{dbPath: dbPath}
	projectList, err := store.List()
	if err != nil {
		t.Fatalf("Failed to list projects: %v", err)
	}

	if len(projectList) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projectList))
	}

	if projectList[0].FolderURI != "file:///test/project" {
		t.Errorf("Expected folderUri 'file:///test/project', got '%s'", projectList[0].FolderURI)
	}
}


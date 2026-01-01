package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	stateDBPath        = "~/Library/Application Support/Cursor/User/globalStorage/state.vscdb"
	vscodeRemotePrefix = "vscode-remote://"
)

// Project represents a Cursor project entry
type Project struct {
	FolderURI string `json:"folderUri"`
	Label     string `json:"label"`
}

// ToAlfredItem converts a Project to an Alfred Item
func (p *Project) ToAlfredItem() Item {
	if p.FolderURI == "" {
		return Item{}
	}

	if strings.HasPrefix(p.FolderURI, vscodeRemotePrefix) {
		return p.formatRemote()
	}
	return p.formatLocal()
}

// formatLocal formats a local project as an Alfred Item
func (p *Project) formatLocal() Item {
	path := strings.TrimPrefix(p.FolderURI, "file://")
	folderName := filepath.Base(path)
	displayName := p.Label
	if displayName == "" {
		displayName = folderName
	}

	return Item{
		Title:    "📁 " + displayName,
		Subtitle: p.FolderURI,
		Arg:      p.FolderURI,
	}
}

// formatRemote formats a remote project as an Alfred Item
func (p *Project) formatRemote() Item {
	// Parse vscode-remote://authority/path format
	trimmedURI := strings.TrimPrefix(p.FolderURI, vscodeRemotePrefix)
	parts := strings.SplitN(trimmedURI, "/", 2)
	if len(parts) < 2 {
		parts = append(parts, "")
	}

	authority := parts[0]
	remotePath := "/" + parts[1]

	// URL decode
	decodedAuthority, err := url.QueryUnescape(authority)
	if err == nil {
		authority = decodedAuthority
	}

	folderName := filepath.Base(remotePath)
	displayName := p.Label
	if displayName == "" {
		displayName = fmt.Sprintf("%s [%s]", folderName, authority)
	}

	// Remove unnecessary suffix pattern (e.g., " in project-name-0 (undefined)")
	displayName = removeSuffixPattern(displayName)

	fullURI := vscodeRemotePrefix + authority + remotePath

	return Item{
		Title:    "🌐 " + displayName,
		Subtitle: displayName,
		Arg:      fullURI,
	}
}

// ProjectStore handles fetching recent projects from Cursor's database
type ProjectStore struct {
	dbPath string
}

// NewProjectStore creates a new ProjectStore instance
func NewProjectStore() (*ProjectStore, error) {
	expandedPath, err := expandPath(stateDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand path: %w", err)
	}

	return &ProjectStore{
		dbPath: expandedPath,
	}, nil
}

// List fetches and returns recent projects
func (s *ProjectStore) List() ([]Project, error) {
	if err := s.validateEnvironment(); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", s.dbPath+"?mode=ro")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var value string
	err = db.QueryRow("SELECT value FROM ItemTable WHERE key = ?", "history.recentlyOpenedPathsList").Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no recently opened projects found")
		}
		return nil, fmt.Errorf("database query error: %w", err)
	}

	var data struct {
		Entries []Project `json:"entries"`
	}
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Filter entries that have folderUri
	var projects []Project
	for _, entry := range data.Entries {
		if entry.FolderURI != "" {
			projects = append(projects, entry)
		}
	}

	return projects, nil
}

func (s *ProjectStore) validateEnvironment() error {
	cursorAppPath := "/Applications/Cursor.app"
	if _, err := os.Stat(cursorAppPath); os.IsNotExist(err) {
		return fmt.Errorf("Cursor app not found")
	}

	if _, err := os.Stat(s.dbPath); os.IsNotExist(err) {
		return fmt.Errorf("Cursor database file not found: %s", s.dbPath)
	}

	return nil
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return strings.Replace(path, "~", homeDir, 1), nil
	}
	return path, nil
}

func removeSuffixPattern(s string) string {
	// Remove patterns like " in project-name-0 (undefined)"
	// Regex matches: " in " followed by any string, then "-\d+ (undefined)" pattern
	re := regexp.MustCompile(` in [^\[\]\(\)]+-\d+ \(undefined\)`)
	return re.ReplaceAllString(s, "")
}

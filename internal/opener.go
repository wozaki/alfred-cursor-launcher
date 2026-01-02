package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	cursorAppPath = "/Applications/Cursor.app"
	cursorBinPath = "/Applications/Cursor.app/Contents/MacOS/Cursor"
)

// Opener handles opening projects in Cursor
type Opener struct{}

// NewOpener creates a new Opener instance
func NewOpener() *Opener {
	return &Opener{}
}

// Open opens a project in Cursor
func (o *Opener) Open(uri string) error {
	if err := o.validateEnvironment(); err != nil {
		return err
	}

	if strings.HasPrefix(uri, "file://") {
		return o.openLocal(uri)
	}
	return o.openRemote(uri)
}

func (o *Opener) validateEnvironment() error {
	if _, err := os.Stat(cursorAppPath); os.IsNotExist(err) {
		return fmt.Errorf("cursor app not found")
	}

	if _, err := os.Stat(cursorBinPath); os.IsNotExist(err) {
		return fmt.Errorf("cursor binary not found")
	}

	return nil
}

func (o *Opener) openLocal(uri string) error {
	path := strings.TrimPrefix(uri, "file://")

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("specified path does not exist: %s", path)
	}

	// Try to execute directly with binary first
	escapedPath := filepath.Clean(path)
	cmd := exec.Command(cursorBinPath, "--new-window", escapedPath)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Fallback to osascript if binary execution fails
	if err := exec.Command("open", "-a", "Cursor").Run(); err != nil {
		return fmt.Errorf("failed to launch Cursor app: %w", err)
	}

	// Wait for app to launch
	time.Sleep(2 * time.Second)

	// Open project with osascript
	escapedPath = strings.ReplaceAll(escapedPath, `"`, `\"`)
	script := fmt.Sprintf(`tell application "Cursor" to open location "file://%s"`, escapedPath)
	cmd = exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open project: %w", err)
	}

	return nil
}

func (o *Opener) openRemote(uri string) error {
	// Try with binary first
	escapedURI := strings.ReplaceAll(uri, `"`, `\"`)
	cmd := exec.Command(cursorBinPath, "--new-window", "--folder-uri", escapedURI)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Fallback to osascript if binary execution fails
	if err := exec.Command("open", "-a", "Cursor").Run(); err != nil {
		return fmt.Errorf("failed to launch Cursor app: %w", err)
	}

	// Wait for app to launch
	time.Sleep(2 * time.Second)

	// Open remote project with osascript
	escapedURI = strings.ReplaceAll(uri, `"`, `\"`)
	script := fmt.Sprintf(`tell application "Cursor" to open location "%s"`, escapedURI)
	cmd = exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open remote project: %w", err)
	}

	return nil
}

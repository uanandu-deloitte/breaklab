package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"breaklab/internal/catalog"
)

// Session is the state kept between commands in .breaklab/session.json.
type Session struct {
	Lab       string `json:"lab"`
	Container string `json:"container"`
	Variant   string `json:"variant"`
	HintsUsed int    `json:"hints_used"`

	// AttemptIndex is the index into progress.json's Attempts this
	// session belongs to, so each command's own process knows which
	// record to update without having to search for it.
	AttemptIndex int `json:"attempt_index"`
}

// LoadSession reads the session file. A missing file means no lab is active
// and returns (nil, nil); a corrupt file is an error naming the path.
func LoadSession(path string) (*Session, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("%s is corrupt (%v); delete it to clear the session", path, err)
	}
	if s.Lab == "" || s.Container == "" {
		return nil, fmt.Errorf("%s is corrupt (missing lab or container); delete it to clear the session", path)
	}
	// session.json is hand-editable, and Container is passed to docker as
	// an argument: anything other than the exact name Start would have
	// written (e.g. "-f", parsed by docker as a flag) is rejected here.
	// Variant is the one field legitimately empty — a lab with no
	// variants at all (selectVariant) never has one to record.
	validVariant := s.Variant == "" || catalog.ValidID(s.Variant)
	if !catalog.ValidID(s.Lab) || !validVariant || s.Container != containerName(s.Lab) {
		return nil, fmt.Errorf("%s is corrupt (invalid lab, variant, or container); delete it to clear the session", path)
	}
	return &s, nil
}

// SaveSession writes the session atomically (temp file + rename) with mode 0600.
func SaveSession(path string, s Session) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, "session-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// ClearSession deletes the session file; a missing file is fine.
func ClearSession(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

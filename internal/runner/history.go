package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// LastVariant returns the variant most recently recorded for labID in the
// history file at path, or "" if none has been recorded yet. A missing
// file is not an error; a corrupt file is.
func LastVariant(path, labID string) (string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", path, err)
	}

	var history map[string]string
	if err := json.Unmarshal(data, &history); err != nil {
		return "", fmt.Errorf("%s is corrupt (%v); delete it to clear variant history", path, err)
	}
	return history[labID], nil
}

// RecordVariant records variantID as the most recent variant used for
// labID, merging into any existing history and writing atomically (temp
// file + rename, mode 0600) like session.json. A corrupt existing file is
// discarded rather than blocking the write.
func RecordVariant(path, labID, variantID string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}

	history := map[string]string{}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &history)
	}
	history[labID] = variantID

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, "history-*.tmp")
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

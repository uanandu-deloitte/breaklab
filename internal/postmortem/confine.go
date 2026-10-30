package postmortem

import (
	"fmt"
	"os"
	"path/filepath"
)

// Confined returns recorded, cleaned and made absolute, only if it names a
// regular .md file directly inside dir — not in a subdirectory, not
// outside dir via "..", and not a symlink. Every postmortem path that
// comes from progress.json (a hand-editable file) goes through this
// before it's read, so an edited path can't make BreakLab read or serve
// an arbitrary file.
func Confined(dir, recorded string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(recorded)
	if err != nil {
		return "", err
	}

	if filepath.Dir(absPath) != absDir || filepath.Ext(absPath) != ".md" {
		return "", fmt.Errorf("postmortem path %q is not a .md file directly inside %s", recorded, dir)
	}
	info, err := os.Lstat(absPath)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("postmortem path %q is not a regular file", recorded)
	}
	return absPath, nil
}

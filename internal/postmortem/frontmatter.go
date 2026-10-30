package postmortem

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// SplitFrontMatter separates a postmortem file's leading "---"-delimited
// YAML block from the markdown body that follows it.
func SplitFrontMatter(content string) (raw, body string, err error) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return "", "", fmt.Errorf("no front matter: file must start with a %q line", "---")
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			return strings.Join(lines[1:i], "\n"), strings.Join(lines[i+1:], "\n"), nil
		}
	}
	return "", "", fmt.Errorf("no front matter: missing closing %q line", "---")
}

// ParseFrontMatter decodes a front matter block's raw YAML.
func ParseFrontMatter(raw string) (FrontMatter, error) {
	var fm FrontMatter
	if err := yaml.Unmarshal([]byte(raw), &fm); err != nil {
		return FrontMatter{}, fmt.Errorf("parsing front matter: %w", err)
	}
	return fm, nil
}

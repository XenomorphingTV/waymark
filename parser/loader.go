package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Load(path string) (string, error) {
	return load(path, make(map[string]bool))
}

func load(path string, seen map[string]bool) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("include %q: %w", path, err)
	}

	if seen[abs] {
		return "", fmt.Errorf("circular include: %s", abs)
	}
	seen[abs] = true

	data, err := os.ReadFile(abs)
	if err != nil {
		return "", fmt.Errorf("include %q: %w", path, err)
	}

	dir := filepath.Dir(abs)
	var out strings.Builder

	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "include ") {
			quoted := strings.TrimPrefix(trimmed, "include ")
			quoted = strings.Trim(quoted, `"`)

			included, err := load(filepath.Join(dir, quoted), seen)
			if err != nil {
				return "", err
			}
			out.WriteString(included)
			out.WriteString("\n")
		} else {
			out.WriteString(line)
			out.WriteString("\n")
		}
	}

	return out.String(), nil
}

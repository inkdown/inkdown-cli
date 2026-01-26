package validate

import (
	"encoding/json"
	"fmt"
	"inkdown-cli/utils"
	"os"
	"path/filepath"
	"strings"
)

var allowedExtensions = map[string]bool{
	".ts": true, ".js": true, ".json": true, ".md": true,
	".css": true, ".png": true, ".jpg": true, ".jpeg": true,
	".svg": true, ".gitignore": true, ".ym": true, ".yaml": true,
	".yml": true, ".mjs": true,
}

var forbiddenTokens = []struct {
	Token   string
	Message string
}{
	{"window.", "Direct access to \"window\" is forbidden. Use platform-agnostic abstractions."},
	{"document.", "Direct access to \"document\" is forbidden. Use platform-agnostic abstractions."},
	{"innerHTML", "Usage of \"innerHTML\" is forbidden."},
	{"outerHTML", "Usage of \"outerHTML\" is forbidden."},
	{"@codemirror/", "Direct imports from \"@codemirror/\" are forbidden. Use @inkdown/core editor abstractions."},
	{"@tauri-apps/", "Direct imports from \"@tauri-apps/\" are forbidden. Use @inkdown/core native abstractions."},
}

type PackageJSON interface {
	GetName() string
	GetVersion() string
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func ValidatePlugin(dir string) error {
	utils.Info("Validating plugin in: %s", dir)

	hasErrors := false

	if _, err := os.Stat(filepath.Join(dir, "manifest.json")); os.IsNotExist(err) {
		utils.Error("Missing 'manifest.json'")
		hasErrors = true
	} else {
		content, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
		if err == nil {
			var pkg Package
			if err := json.Unmarshal(content, &pkg); err != nil {
				utils.Error("Invalid 'manifest.json': %v", err)
				hasErrors = true
			} else {
				if pkg.Name == "" {
					utils.Error("manifest.json missing 'name'")
					hasErrors = true
				}
				if pkg.Version == "" {
					utils.Error("manifest.json missing 'version'")
					hasErrors = true
				}
			}
		}
	}

	srcDir := filepath.Join(dir, "src")
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		utils.Error("Missing 'src' directory. Plugin source must be in 'src'.")
		return fmt.Errorf("plugin validation failed")
	}

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		name := info.Name()
		if !allowedExtensions[ext] && name != "LICENSE" && name != "README" && name != "README.md" {
			utils.Error("Forbidden file type found: %s (Extension \"%s\" is not in whitelist)", path, ext)
			hasErrors = true
		}

		if ext == ".ts" || ext == ".js" || ext == ".mjs" {
			content, err := os.ReadFile(path)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				for i, line := range lines {
					trimmed := strings.TrimSpace(line)
					if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") {
						continue
					}
					for _, rule := range forbiddenTokens {
						if strings.Contains(line, rule.Token) {
							utils.Error("Forbidden token \"%s\" found in %s:%d", rule.Token, path, i+1)
							utils.Note(rule.Message)
							hasErrors = true
						}
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error scanning files: %v", err)
	}

	if hasErrors {
		return fmt.Errorf("plugin validation failed")
	}

	utils.Success("Plugin validation passed!")
	return nil
}

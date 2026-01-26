package validate

import (
	"encoding/json"
	"fmt"
	"inkdown-cli/utils"
	"os"
	"path/filepath"
)

type ThemeManifest struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Modes   []string `json:"modes"`
}

func ValidateTheme(dir string) error {
	utils.Info("Validating theme in: %s", dir)

	hasErrors := false

	// 1. Check theme.json
	manifestPath := filepath.Join(dir, "theme.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		utils.Error("Missing 'theme.json'")
		return fmt.Errorf("theme validation failed")
	}

	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("could not read theme.json: %v", err)
	}

	var theme ThemeManifest
	if err := json.Unmarshal(content, &theme); err != nil {
		utils.Error("Invalid 'theme.json': %v", err)
		hasErrors = true
	} else {
		if theme.Name == "" {
			utils.Error("theme.json missing 'name'")
			hasErrors = true
		}
		if theme.Version == "" {
			utils.Error("theme.json missing 'version'")
			hasErrors = true
		}

		// 2. Check CSS files based on modes
		modes := theme.Modes
		if len(modes) == 0 {
			modes = []string{"dark"} // default if not specified? Assuming similar logic to registry
		}

		for _, mode := range modes {
			cssFile := filepath.Join(dir, mode+".css")
			if _, err := os.Stat(cssFile); os.IsNotExist(err) {
				utils.Error("Missing required CSS file for mode '%s': %s", mode, mode+".css")
				hasErrors = true
			}
		}
	}

	if hasErrors {
		return fmt.Errorf("theme validation failed")
	}

	utils.Success("Theme validation passed!")
	return nil
}

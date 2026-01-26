package generator

import (
	"inkdown-cli/internal/templates"
	"os"
	"path/filepath"
	"strings"
)

func CopyPluginTemplate(dir *string, abs *string, name *string, desc *string) error {
	entries, err := templates.Templates.ReadDir(*dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(*dir, entry.Name())
		destPath := filepath.Join(*abs, entry.Name())

		if entry.IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
			CopyPluginTemplate(&srcPath, &destPath, name, desc)
			continue
		}

		data, err := templates.Templates.ReadFile(srcPath)
		if err != nil {
			return err
		}

		if entry.Name() == "manifest.json" || entry.Name() == "package.json" {
			content := string(data)

			if *name != "" {
				content = strings.ReplaceAll(content, "My plugin", *name)
				if entry.Name() == "manifest.json" {
					content = strings.ReplaceAll(content, "my-plugin", strings.ToLower(strings.ReplaceAll(*name, " ", "-")))
				}
			}

			if *desc != "" {
				content = strings.ReplaceAll(content, "A custom plugin made for inkdown", *desc)
			}

			data = []byte(content)
		}

		os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		os.WriteFile(destPath, data, 0644)
	}

	return nil
}

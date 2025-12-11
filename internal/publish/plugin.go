package publish

import (
	"encoding/json"
	"fmt"
	"os"
)

func PublishPlugin(dir *string) (string, error) {
	entries, err := os.ReadDir(*dir)

	if err != nil {
		fmt.Printf("Error while reading the plugin directory: %s\n", err)

		return "", err
	}

	for _, entry := range entries {
		if entry.Name() == "package.json" {
			content, err := os.ReadFile(entry.Name())

			if err != nil {
				fmt.Printf("Error while reading the package.json file: %s\n", err)

				return "", err
			}

			j, err := json.Marshal(content)

			if err != nil {
				fmt.Printf("Error while parsing the package.json file: %s\n", err)

				return "", err
			}

			fmt.Println(j)
		}
	}

	return "link", nil
}

package plugin

import (
	"fmt"
	"path/filepath"

	"inkdown-cli/internal/publish"

	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a plugin",
	RunE: func(cmd *cobra.Command, args []string) error {
		if path == "" {
			path = "."
		}

		fmt.Printf("Initializing Inkdown publisher\n")
		fmt.Printf("  Path: %s\n", path)
		if name != "" {
			fmt.Printf("  Name: %s\n", name)
		}

		_, err := filepath.Abs(path)

		if err != nil {
			fmt.Printf("This path does not exist: %s\n", path)
			return err
		}

		link, err := publish.PublishPlugin(&path)

		if err != nil {
			fmt.Printf("Failed to publish plugin: %v\n", err)
			return err
		}
		fmt.Printf("Plugin published successfully: %s\n", link)

		return nil
	},
}

func init() {

	publishCmd.Flags().StringVarP(&path, "path", "p", ".", "Path to the plugin")
	publishCmd.Flags().StringVarP(&name, "name", "n", "", "Plugin name")
	publishCmd.Flags().StringVarP(&description, "description", "d", "", "Plugin description")

	PluginCmd.AddCommand(publishCmd)
}

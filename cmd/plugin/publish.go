package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"inkdown-cli/internal/publish"

	"github.com/spf13/cobra"
)

var pluginPath string

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a plugin",
	RunE: func(cmd *cobra.Command, args []string) error {
		if pluginPath == "" {
			var err error
			pluginPath, err = os.Getwd()
			if err != nil {
				pluginPath = "."
			}
		}

		if name != "" {
			fmt.Printf("  Name: %s\n", name)
		}

		_, err := filepath.Abs(pluginPath)

		if err != nil {
			fmt.Printf("This path does not exist: %s\n", pluginPath)
			return err
		}

		link, err := publish.PublishPlugin(&pluginPath)

		if err != nil {
			fmt.Printf("Failed to publish plugin: %v\n", err)
			return err
		}

		fmt.Printf("Plugin PR published successfully, you can check in this link: %s\n", link)

		return nil
	},
}

func init() {

	publishCmd.Flags().StringVarP(&pluginPath, "path", "d", ".", "Path to the plugin")

	PluginCmd.AddCommand(publishCmd)
}

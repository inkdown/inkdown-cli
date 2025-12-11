package plugin

import (
	"fmt"
	"inkdown-cli/internal/generator"

	"github.com/spf13/cobra"
)

var (
	path        string
	name        string
	description string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new inkdown plugin",
	RunE: func(cmd *cobra.Command, args []string) error {
		if path == "" {
			path = "."
		}

		fmt.Printf("Initializing Inkdown plugin project...\n")
		fmt.Printf("  Path: %s\n", path)
		if name != "" {
			fmt.Printf("  Name: %s\n", name)
		}

		templateDir := "plugin"

		generator.CopyPluginTemplate(&templateDir, &path, &name, &description)
		return nil
	},
}

func init() {

	initCmd.Flags().StringVarP(&path, "path", "p", ".", "Path to initialize the plugin")
	initCmd.Flags().StringVarP(&name, "name", "n", "", "Plugin name")
	initCmd.Flags().StringVarP(&description, "description", "d", "", "Plugin description")

	PluginCmd.AddCommand(initCmd)
}

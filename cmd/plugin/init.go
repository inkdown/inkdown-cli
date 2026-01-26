package plugin

import (
	"fmt"
	"path/filepath"

	"inkdown-cli/internal/generator"

	"github.com/spf13/cobra"
)

var (
	initPath    string
	name        string
	description string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new inkdown plugin",
	RunE: func(cmd *cobra.Command, args []string) error {
		if initPath == "" {
			initPath = "."
		}

		abs, err := filepath.Abs(initPath)

		if err != nil {
			fmt.Printf("Could not initialize the plugin: %s\n", err)
		}

		fmt.Printf("Initializing Inkdown plugin project in the path: %s\n", abs)

		if name != "" {
			fmt.Printf("  Name: %s\n", name)
		}

		templateDir := "plugin"

		generator.CopyPluginTemplate(&templateDir, &abs, &name, &description)
		return nil
	},
}

func init() {

	initCmd.Flags().StringVarP(&initPath, "path", "p", ".", "Path to initialize the plugin")
	initCmd.Flags().StringVarP(&name, "name", "n", "", "Plugin name")
	initCmd.Flags().StringVarP(&description, "description", "d", "", "Plugin description")

	PluginCmd.AddCommand(initCmd)
}

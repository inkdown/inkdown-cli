package cmd

import (
	"fmt"
	"inkdown-cli/internal/generator"

	"github.com/spf13/cobra"
)

var (
	pluginPath string
	pluginName string
	pluginDesc string
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Initialize a new Inkdown plugin project",
	Run: func(cmd *cobra.Command, args []string) {

		if pluginPath == "" {
			pluginPath = "."
		}

		fmt.Printf("Initializing Inkdown plugin project...\n")
		fmt.Printf("  Path: %s\n", pluginPath)
		if pluginName != "" {
			fmt.Printf("  Name: %s\n", pluginName)
		}

		templateDir := "plugin"

		generator.CopyPluginTemplate(&templateDir, &pluginPath, &pluginName, &pluginDesc)
	},
}

func init() {
	pluginCmd.Flags().StringVarP(&pluginPath, "path", "p", "", "Path to initialize the plugin project (default: current directory)")
	pluginCmd.Flags().StringVar(&pluginName, "name", "", "Name of the plugin")
	pluginCmd.Flags().StringVar(&pluginDesc, "desc", "", "Description of the plugin")

	initCmd.AddCommand(pluginCmd)
}

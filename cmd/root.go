package cmd

import (
	"inkdown-cli/cmd/plugin"
	"inkdown-cli/cmd/theme"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ink",
	Short: "Inkdown cli for publishing plugins and themes easily",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(plugin.PluginCmd)
	rootCmd.AddCommand(theme.ThemeCmd)

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(logoutCmd)
}

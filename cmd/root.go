package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "ink",
	Short: "Inkdown cli for publishing plugins and themes",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(logoutCmd)
}

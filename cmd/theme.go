package cmd

import (
	"fmt"

	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	themePath string
	themeName string
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Initialize a new Inkdown theme project",
	Run: func(cmd *cobra.Command, args []string) {

		if themePath == "" {
			themePath = "."
		}

		fmt.Printf("Initializing Inkdown theme project...\n")
		fmt.Printf("  Path: %s\n", themePath)
		if themeName != "" {
			fmt.Printf("  Name: %s\n", themeName)
		}

		_, err := filepath.Abs(themePath)

		if err != nil {
			fmt.Printf("This path does not exist: %s", themePath)
		}

	},
}

func init() {
	themeCmd.Flags().StringVarP(&themePath, "path", "p", "", "Path to initialize the theme project (default: current directory)")
	themeCmd.Flags().StringVar(&themeName, "name", "", "Name of the theme")

	initCmd.AddCommand(themeCmd)
}

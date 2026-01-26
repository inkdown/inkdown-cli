package theme

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	initPath string
	name     string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Inkdown theme project",
	Run: func(cmd *cobra.Command, args []string) {

		if initPath == "" {
			initPath = "."
		}

		fmt.Printf("  Path: %s\n", initPath)
		if name != "" {
			fmt.Printf("  Name: %s\n", name)
		}

		abs, err := filepath.Abs(initPath)

		if err != nil {
			fmt.Printf("This path does not exist: %s", initPath)
		}

		fmt.Printf("Initializing Inkdown theme project in path: %s\n", abs)

	},
}

func init() {
	initCmd.Flags().StringVarP(&initPath, "path", "p", ".", "Path to initialize the theme")
	initCmd.Flags().StringVarP(&name, "name", "n", "", "Theme name")

	ThemeCmd.AddCommand(initCmd)
}

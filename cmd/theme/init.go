package theme

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	path string
	name string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Inkdown theme project",
	Run: func(cmd *cobra.Command, args []string) {

		if path == "" {
			path = "."
		}

		fmt.Printf("Initializing Inkdown theme project...\n")
		fmt.Printf("  Path: %s\n", path)
		if name != "" {
			fmt.Printf("  Name: %s\n", name)
		}

		_, err := filepath.Abs(path)

		if err != nil {
			fmt.Printf("This path does not exist: %s", path)
		}

	},
}

func init() {
	initCmd.Flags().StringVarP(&path, "path", "p", ".", "Path to initialize the theme")
	initCmd.Flags().StringVarP(&name, "name", "n", "", "Theme name")

	ThemeCmd.AddCommand(initCmd)
}

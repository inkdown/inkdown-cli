package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"inkdown-cli/internal/validate"

	"github.com/spf13/cobra"
)

var validatePath string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a plugin project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if validatePath == "" {
			validatePath = "."
		}

		abs, err := filepath.Abs(validatePath)
		if err != nil {
			return err
		}

		// Check if it's a directory
		info, err := os.Stat(abs)
		if err != nil {
			return fmt.Errorf("path not found: %s", abs)
		}
		if !info.IsDir() {
			return fmt.Errorf("path must be a directory: %s", abs)
		}

		if err := validate.ValidatePlugin(abs); err != nil {
			// Error is already printed by the validator
			return err
		}

		return nil
	},
}

func init() {
	validateCmd.Flags().StringVarP(&validatePath, "path", "p", ".", "Path to the plugin")
	PluginCmd.AddCommand(validateCmd)
}

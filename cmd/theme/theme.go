package theme

import (
	"github.com/spf13/cobra"
)

var ThemeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Initialize or publish a new Inkdown theme project",
}

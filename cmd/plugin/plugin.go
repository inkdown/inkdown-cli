package plugin

import (
	"github.com/spf13/cobra"
)

var PluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Initialize or publish a new Inkdown plugin project",
}

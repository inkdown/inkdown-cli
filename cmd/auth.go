package cmd

import (
	"inkdown-cli/internal/auth"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Run: func(cmd *cobra.Command, args []string) {
		auth.Auth()
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Inkdown CLI",
	Run: func(cmd *cobra.Command, args []string) {
		auth.Logout()
	},
}

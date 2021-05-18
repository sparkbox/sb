package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sb",
	Short: "A cli to automate various things around Sparkbox.",
	Long: `For now we use Slack as our Authentication service in lieu of something like LDAP.
We currently only support generating SSH certificates for the various hosts we manage.
In the future we could handle things like onboarding here.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

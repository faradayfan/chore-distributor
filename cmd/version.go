package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version, commit SHA, and build date of chore-distributor.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("chore-distributor %s\n", Version)
		fmt.Printf("  Commit: %s\n", CommitSHA)
		fmt.Printf("  Built:  %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chore-distributor",
	Short: "A fair chore distribution tool for families",
	Long: `Chore Distributor fairly distributes household chores among family members
based on earning potential, with optional effort capacity limits.

Features:
  - Fair Distribution: Balances chores by amount earned
  - Effort Capacity: Set maximum effort limits for individuals
  - Randomization: Shuffles assignments each run
  - JSON Configuration: Easy to modify chores and people`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

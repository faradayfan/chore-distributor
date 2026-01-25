package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/faradayfan/chore-distributor/internal/config"
	"github.com/faradayfan/chore-distributor/internal/distributor"
	"github.com/faradayfan/chore-distributor/internal/models"
	"github.com/faradayfan/chore-distributor/internal/notes"
	"github.com/faradayfan/chore-distributor/internal/sms"
	"github.com/spf13/cobra"
)

var (
	configPath string
	verbose    bool
	sendSMS    bool
	noteName   string
	dryRun     bool
	confirm    bool
)

var distributeCmd = &cobra.Command{
	Use:   "distribute",
	Short: "Distribute chores among family members",
	Long: `Distribute chores among family members based on the configuration file.

The distribution algorithm:
  1. Loads chores and people from the JSON configuration file
  2. Shuffles chores and sorts by earning amount (highest first)
  3. Assigns each chore to the person with the lowest current earnings
     who has available capacity
  4. Displays the final distribution
  5. Optionally sends iMessage notifications to each person (macOS only)
  6. Optionally saves to an Apple Note (macOS only)`,
	Example: `  # Use default config file (chores_config.json)
  chore-distributor distribute

  # Use a custom config file
  chore-distributor distribute --config /path/to/config.json

  # Show difficulty and capacity information
  chore-distributor distribute --verbose

  # Send iMessage notifications to everyone (macOS only)
  chore-distributor distribute --sms

  # Save to an Apple Note called "Chore History"
  chore-distributor distribute --note "Chore History"

  # Do everything: verbose output, SMS, and save to note
  chore-distributor distribute -v --sms --note "Chore History"

  # Confirm before sending messages and saving to notes
  chore-distributor distribute --sms --note "Chore History" --confirm

  # Preview all actions without sending/saving (dry run)
  chore-distributor distribute --sms --note "Chore History" --dry-run`,
	Run: func(cmd *cobra.Command, args []string) {
		runDistribute()
	},
}

func runDistribute() {
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	for {
		for i := range cfg.People {
			cfg.People[i].Chores = []models.Chore{}
			cfg.People[i].TotalDifficulty = 0
			cfg.People[i].TotalEarned = 0

			// Recalculate from pre-assigned chores
			for _, chore := range cfg.People[i].PreAssignedChores {
				cfg.People[i].TotalDifficulty += chore.Difficulty
				cfg.People[i].TotalEarned += chore.Earned
			}
		}

		cfg.People = distributor.Distribute(cfg.Chores, cfg.People)

		opts := distributor.PrintOptions{
			Verbose: verbose,
		}
		distributor.PrintDistribution(os.Stdout, cfg.People, opts)

		if confirm && (noteName != "" || sendSMS) && !dryRun {
			result := promptConfirmation()
			switch result {
			case "retry":
				fmt.Println("--- Retrying distribution ---")
				continue
			case "cancel":
				fmt.Println("Cancelled.")
				os.Exit(0)
			case "confirm":
			}
		}

		break
	}

	if noteName != "" {
		if !notes.IsSupported() && !dryRun {
			fmt.Fprintf(os.Stderr, "Error: Apple Notes is only supported on macOS\n")
			os.Exit(1)
		}

		fmt.Println("\n--- Saving to Apple Notes ---")
		writer := notes.NewWriter(noteName, dryRun)
		if err := writer.PrependChoreList(cfg.People, verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving to Notes: %v\n", err)
			os.Exit(1)
		}
	}

	if sendSMS {
		if !sms.IsSupported() && !dryRun {
			fmt.Fprintf(os.Stderr, "Error: iMessage is only supported on macOS\n")
			os.Exit(1)
		}

		fmt.Println("\n--- Sending iMessage Notifications ---")
		sender := sms.NewSender(dryRun)
		if err := sender.SendChoreAssignments(cfg.People, verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending messages: %v\n", err)
			os.Exit(1)
		}
	}
}

func promptConfirmation() string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n[C]onfirm, [R]etry, or [A]bort? ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			return "cancel"
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "c", "confirm":
			return "confirm"
		case "r", "retry":
			return "retry"
		case "a", "abort", "cancel":
			return "cancel"
		default:
			fmt.Println("Please enter C (confirm), R (retry), or A (abort)")
		}
	}
}

func init() {
	rootCmd.AddCommand(distributeCmd)

	distributeCmd.Flags().StringVarP(&configPath, "config", "c", "chores_config.json",
		"Path to the JSON configuration file")
	distributeCmd.Flags().BoolVarP(&verbose, "verbose", "v", false,
		"Show difficulty and capacity information")
	distributeCmd.Flags().BoolVarP(&sendSMS, "sms", "s", false,
		"Send iMessage notifications (macOS only). Contact can be phone number or Apple ID email.")
	distributeCmd.Flags().StringVarP(&noteName, "note", "o", "",
		"Save chore list to an Apple Note with this name (macOS only). Creates note if it doesn't exist.")
	distributeCmd.Flags().BoolVarP(&confirm, "confirm", "i", false,
		"Prompt for confirmation before sending messages and saving to notes")
	distributeCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false,
		"Preview actions without actually sending messages or saving to notes")
}

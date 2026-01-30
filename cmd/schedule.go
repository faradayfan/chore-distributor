package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const (
	launchAgentLabel = "com.faradayfan.chore-distributor"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage automatic scheduling (macOS only)",
	Long: `Set up and manage automatic chore distribution using macOS launchd.

This command helps you schedule the chore distributor to run automatically
on a specific day and time each week.`,
}

var scheduleInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and enable automatic scheduling",
	Long: `Creates a launchd agent to automatically run chore distribution.

You'll be prompted to configure:
  - Day of week (Sunday-Saturday)
  - Time of day
  - Config file path
  - Whether to send SMS
  - Whether to save to Notes`,
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "darwin" {
			fmt.Fprintf(os.Stderr, "Error: Automatic scheduling is only supported on macOS\n")
			os.Exit(1)
		}
		installSchedule()
	},
}

var scheduleUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove automatic scheduling",
	Long:  `Unloads and removes the launchd agent, stopping automatic chore distribution.`,
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "darwin" {
			fmt.Fprintf(os.Stderr, "Error: Automatic scheduling is only supported on macOS\n")
			os.Exit(1)
		}
		uninstallSchedule()
	},
}

var scheduleStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check scheduling status",
	Long:  `Displays whether automatic scheduling is enabled and shows the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "darwin" {
			fmt.Fprintf(os.Stderr, "Error: Automatic scheduling is only supported on macOS\n")
			os.Exit(1)
		}
		showScheduleStatus()
	},
}

var scheduleTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Run scheduled task immediately",
	Long:  `Manually triggers the scheduled task once for testing purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "darwin" {
			fmt.Fprintf(os.Stderr, "Error: Automatic scheduling is only supported on macOS\n")
			os.Exit(1)
		}
		testSchedule()
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.AddCommand(scheduleInstallCmd)
	scheduleCmd.AddCommand(scheduleUninstallCmd)
	scheduleCmd.AddCommand(scheduleStatusCmd)
	scheduleCmd.AddCommand(scheduleTestCmd)
}

func getLaunchAgentPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(home, "Library", "LaunchAgents", launchAgentLabel+".plist")
}

func getExecutablePath() string {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting executable path: %v\n", err)
		os.Exit(1)
	}
	// Resolve symlinks
	realExe, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return exe
	}
	return realExe
}

func promptScheduleConfig() map[string]string {
	reader := bufio.NewReader(os.Stdin)
	config := make(map[string]string)

	fmt.Println("=== Chore Distributor Scheduling Setup ===")

	// Day of week
	fmt.Println("Select day of week:")
	fmt.Println("  0 - Sunday")
	fmt.Println("  1 - Monday")
	fmt.Println("  2 - Tuesday")
	fmt.Println("  3 - Wednesday")
	fmt.Println("  4 - Thursday")
	fmt.Println("  5 - Friday")
	fmt.Println("  6 - Saturday")
	for {
		fmt.Print("\nEnter day (0-6): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		day, err := strconv.Atoi(input)
		if err == nil && day >= 0 && day <= 6 {
			config["weekday"] = input
			break
		}
		fmt.Println("Invalid input. Please enter a number between 0 and 6.")
	}

	// Hour
	for {
		fmt.Print("\nEnter hour (0-23, e.g., 10 for 10 AM): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		hour, err := strconv.Atoi(input)
		if err == nil && hour >= 0 && hour <= 23 {
			config["hour"] = input
			break
		}
		fmt.Println("Invalid input. Please enter a number between 0 and 23.")
	}

	// Minute
	for {
		fmt.Print("Enter minute (0-59): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		minute, err := strconv.Atoi(input)
		if err == nil && minute >= 0 && minute <= 59 {
			config["minute"] = input
			break
		}
		fmt.Println("Invalid input. Please enter a number between 0 and 59.")
	}

	// Config file path
	fmt.Print("\nEnter config file path (press Enter for 'chores_config.json'): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		input = "chores_config.json"
	}
	config["configPath"] = input

	// SMS
	fmt.Print("\nSend iMessage notifications? (y/n): ")
	input, _ = reader.ReadString('\n')
	config["sms"] = strings.ToLower(strings.TrimSpace(input))

	// Notes
	fmt.Print("Save to Apple Notes? (y/n): ")
	input, _ = reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(input)) == "y" {
		config["notes"] = "y"
		fmt.Print("Enter note name (default: 'Chore History'): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			input = "Chore History"
		}
		config["noteName"] = input
	} else {
		config["notes"] = "n"
	}

	return config
}

func generatePlist(config map[string]string) string {
	exePath := getExecutablePath()
	exeDir := filepath.Dir(exePath)
	logDir := filepath.Join(exeDir, "logs")

	// Build program arguments
	args := []string{
		fmt.Sprintf("        <string>%s</string>", exePath),
		"        <string>distribute</string>",
		"        <string>--config</string>",
		fmt.Sprintf("        <string>%s</string>", config["configPath"]),
	}

	if config["sms"] == "y" {
		args = append(args, "        <string>--sms</string>")
	}

	if config["notes"] == "y" {
		args = append(args,
			"        <string>--note</string>",
			fmt.Sprintf("        <string>%s</string>", config["noteName"]),
		)
	}

	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>

    <key>ProgramArguments</key>
    <array>
%s
    </array>

    <key>StartCalendarInterval</key>
    <dict>
        <key>Weekday</key>
        <integer>%s</integer>
        <key>Hour</key>
        <integer>%s</integer>
        <key>Minute</key>
        <integer>%s</integer>
    </dict>

    <key>StandardOutPath</key>
    <string>%s/stdout.log</string>

    <key>StandardErrorPath</key>
    <string>%s/stderr.log</string>

    <key>RunAtLoad</key>
    <false/>
</dict>
</plist>
`, launchAgentLabel, strings.Join(args, "\n"), config["weekday"], config["hour"], config["minute"], logDir, logDir)

	return plist
}

func installSchedule() {
	config := promptScheduleConfig()

	plistPath := getLaunchAgentPath()
	plistContent := generatePlist(config)

	// Create LaunchAgents directory if it doesn't exist
	launchAgentsDir := filepath.Dir(plistPath)
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating LaunchAgents directory: %v\n", err)
		os.Exit(1)
	}

	// Create logs directory
	exePath := getExecutablePath()
	exeDir := filepath.Dir(exePath)
	logDir := filepath.Join(exeDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logs directory: %v\n", err)
		os.Exit(1)
	}

	// Check if already exists and loaded
	if _, err := os.Stat(plistPath); err == nil {
		fmt.Println("\nSchedule already exists. Unloading existing configuration...")
		runCommand("launchctl", "unload", plistPath)
	}

	// Write plist file
	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing plist file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Created schedule configuration at: %s\n", plistPath)

	// Load the agent
	if err := runCommand("launchctl", "load", plistPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading schedule: %v\n", err)
		fmt.Fprintf(os.Stderr, "You may need to manually load it with: launchctl load %s\n", plistPath)
		os.Exit(1)
	}

	fmt.Println("✓ Schedule loaded successfully")
	fmt.Printf("\nScheduled to run: %s at %s:%s\n",
		getDayName(config["weekday"]),
		config["hour"],
		config["minute"])
	fmt.Printf("Logs will be written to: %s\n", logDir)
	fmt.Println("\nUse 'chore-distributor schedule status' to check the schedule")
	fmt.Println("Use 'chore-distributor schedule test' to run it immediately")
}

func uninstallSchedule() {
	plistPath := getLaunchAgentPath()

	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		fmt.Println("No schedule is currently installed.")
		return
	}

	fmt.Print("Are you sure you want to remove the schedule? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(input)) != "y" {
		fmt.Println("Cancelled.")
		return
	}

	// Unload the agent
	if err := runCommand("launchctl", "unload", plistPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to unload schedule: %v\n", err)
	}

	// Remove the plist file
	if err := os.Remove(plistPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing plist file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Schedule removed successfully")
}

func showScheduleStatus() {
	plistPath := getLaunchAgentPath()

	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		fmt.Println("Status: Not scheduled")
		fmt.Println("\nRun 'chore-distributor schedule install' to set up automatic scheduling.")
		return
	}

	fmt.Println("Status: Scheduled")
	fmt.Printf("Configuration file: %s\n\n", plistPath)

	// Read and display plist content
	content, err := os.ReadFile(plistPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading schedule configuration: %v\n", err)
		return
	}

	// Parse and display key information
	contentStr := string(content)
	fmt.Println("Schedule Details:")

	// Extract weekday
	if strings.Contains(contentStr, "<key>Weekday</key>") {
		lines := strings.Split(contentStr, "\n")
		for i, line := range lines {
			if strings.Contains(line, "<key>Weekday</key>") && i+1 < len(lines) {
				weekday := extractInteger(lines[i+1])
				fmt.Printf("  Day: %s\n", getDayName(weekday))
			}
			if strings.Contains(line, "<key>Hour</key>") && i+1 < len(lines) {
				hour := extractInteger(lines[i+1])
				fmt.Printf("  Hour: %s\n", hour)
			}
			if strings.Contains(line, "<key>Minute</key>") && i+1 < len(lines) {
				minute := extractInteger(lines[i+1])
				fmt.Printf("  Minute: %s\n", minute)
			}
		}
	}

	// Check if it's loaded
	output, _ := runCommandOutput("launchctl", "list")
	if strings.Contains(output, launchAgentLabel) {
		fmt.Println("\n✓ Agent is loaded and active")
	} else {
		fmt.Println("\n⚠ Agent is not currently loaded")
		fmt.Printf("Run: launchctl load %s\n", plistPath)
	}

	// Show log location
	exePath := getExecutablePath()
	exeDir := filepath.Dir(exePath)
	logDir := filepath.Join(exeDir, "logs")
	fmt.Printf("\nLogs: %s\n", logDir)
}

func testSchedule() {
	plistPath := getLaunchAgentPath()

	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		fmt.Println("No schedule is currently installed.")
		fmt.Println("Run 'chore-distributor schedule install' first.")
		return
	}

	fmt.Println("Running scheduled task now...")
	if err := runCommand("launchctl", "start", launchAgentLabel); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting task: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Task started")
	fmt.Println("\nCheck logs for output:")
	exePath := getExecutablePath()
	exeDir := filepath.Dir(exePath)
	logDir := filepath.Join(exeDir, "logs")
	fmt.Printf("  tail -f %s/stdout.log\n", logDir)
	fmt.Printf("  tail -f %s/stderr.log\n", logDir)
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func getDayName(weekday string) string {
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	day, err := strconv.Atoi(weekday)
	if err != nil || day < 0 || day > 6 {
		return weekday
	}
	return days[day]
}

func extractInteger(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "<integer>")
	line = strings.TrimSuffix(line, "</integer>")
	return strings.TrimSpace(line)
}

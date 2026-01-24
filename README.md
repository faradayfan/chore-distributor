# Chore Distributor

A Go command-line tool that fairly distributes household chores among family members based on earning potential, with optional effort capacity limits. Supports sending notifications via iMessage and saving history to Apple Notes (macOS only).

## Features

- **Fair Distribution**: Balances chores by amount earned to ensure everyone gets similar total earnings
- **Effort Capacity**: Set maximum effort limits for individuals (useful for younger kids or those with less time)
- **Randomization**: Shuffles assignments each run to keep things fresh and fair
- **JSON Configuration**: Easy to modify chores and people without touching code
- **iMessage Notifications**: Send chore assignments directly to family members via iMessage (macOS only)
- **Apple Notes Integration**: Save chore history to an Apple Note for record keeping (macOS only)
- **Confirmation Prompt**: Review and retry distributions before committing

## Prerequisites

- **Go**: Version 1.21 or later
- **macOS**: Required for iMessage and Apple Notes features (the core distribution works on any platform)
- **Messages App**: Must be signed into iMessage for SMS notifications
- **Notes App**: Must be signed into iCloud or have local notes enabled

## Installation

1. Clone or download the repository:

   ```bash
   git clone https://github.com/faradayfan/chore-distributor.git
   cd chore-distributor
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the program:

   ```bash
   go build -o chore-distributor ./cmd/chore-distributor
   ```

   Or install directly to your Go bin:

   ```bash
   go install ./cmd/chore-distributor
   ```

## Configuration

Create a JSON configuration file with your chores and family members. See `example.json` for a template.

### Configuration File Format

```json
{
  "chores": [
    {
      "Name": "Kitchen",
      "Difficulty": 6,
      "Earned": 5
    },
    {
      "Name": "Bathroom",
      "Difficulty": 5,
      "Earned": 4
    }
  ],
  "people": [
    {
      "Name": "Alice",
      "Contact": "+15551234567",
      "EffortCapacity": 0
    },
    {
      "Name": "Bob",
      "Contact": "bob@icloud.com",
      "EffortCapacity": 15
    }
  ]
}
```

### Chore Properties

| Property     | Type   | Description                                                 |
| ------------ | ------ | ----------------------------------------------------------- |
| `Name`       | string | The name/description of the chore                           |
| `Difficulty` | int    | How much effort the chore requires (1-10 scale recommended) |
| `Earned`     | int    | How much money/points earned for completing this chore      |

### Person Properties

| Property         | Type   | Description                                                                                                                      |
| ---------------- | ------ | -------------------------------------------------------------------------------------------------------------------------------- |
| `Name`           | string | Person's name                                                                                                                    |
| `Contact`        | string | Phone number (e.g., `+15551234567`) or Apple ID email (e.g., `user@icloud.com`) for iMessage. Leave empty to skip notifications. |
| `EffortCapacity` | int    | Maximum total difficulty they can handle. Set to `0` for no limit.                                                               |

## Usage

### Basic Usage

Distribute chores using the default config file (`chores_config.json`):

```bash
./chore-distributor distribute
```

### Command Line Options

```bash
./chore-distributor distribute [flags]
```

| Flag        | Short | Description                                                          |
| ----------- | ----- | -------------------------------------------------------------------- |
| `--config`  | `-c`  | Path to the JSON configuration file (default: `chores_config.json`)  |
| `--verbose` | `-v`  | Show difficulty and capacity information in output                   |
| `--sms`     | `-s`  | Send iMessage notifications to each person (macOS only)              |
| `--note`    | `-o`  | Save chore list to an Apple Note with this name (macOS only)         |
| `--confirm` | `-i`  | Prompt for confirmation before sending messages and saving to notes  |
| `--dry-run` | `-n`  | Preview actions without actually sending messages or saving to notes |
| `--help`    | `-h`  | Show help information                                                |

### Examples

```bash
# Use a custom config file
./chore-distributor distribute --config /path/to/config.json

# Show detailed output with difficulty and capacity info
./chore-distributor distribute -c example.json --verbose

# Send iMessage notifications to everyone
./chore-distributor distribute -c example.json --sms

# Save to an Apple Note called "Chore History"
./chore-distributor distribute -c example.json --note "Chore History"

# Do everything with confirmation prompt
./chore-distributor distribute -c example.json -v --sms --note "Chore History" --confirm

# Preview what would happen without actually doing it
./chore-distributor distribute -c example.json --sms --note "Chore History" --dry-run
```

### Confirmation Prompt

When using `--confirm`, you'll be prompted after viewing the distribution:

```text
[C]onfirm, [R]etry, or [A]bort?
```

- **C / Confirm**: Proceed with sending messages and saving to notes
- **R / Retry**: Generate a new random distribution
- **A / Abort**: Cancel without sending or saving

This allows you to re-roll the distribution until you're happy with it.

## How the Distribution Algorithm Works

1. Chores are loaded from the configuration file
2. Chores are shuffled randomly, then sorted by earning amount (highest first)
3. Each chore is assigned to the person with:
   - The lowest current total earnings
   - Available capacity (if they have a limit set)
4. If multiple people are tied for lowest earnings, one is randomly selected
5. The final distribution is displayed

## Example Output

### Default Output

```
=== Chore Distribution ===

Alice:
  Chores:
    - Kitchen (Earns: $5)
    - Living Room (Earns: $3)
  Total Earned: $8

Bob:
  Chores:
    - Bathroom (Earns: $4)
    - Family Room (Earns: $2)
  Total Earned: $6
```

### Verbose Output (`--verbose`)

```
=== Chore Distribution ===

Alice:
  Chores:
    - Kitchen (Difficulty: 6, Earns: $5)
    - Living Room (Difficulty: 4, Earns: $3)
  Total Difficulty: 10
  Total Earned: $8

Bob (Effort Capacity: 15):
  Chores:
    - Bathroom (Difficulty: 5, Earns: $4)
    - Family Room (Difficulty: 3, Earns: $2)
  Total Difficulty: 8 / 15
  Total Earned: $6
```

## Apple Notes History

When using `--note`, each distribution is prepended to the note with the date. The note title is preserved, and new entries appear at the top:

```
Chore History

Saturday, January 25, 2026

Alice
• Kitchen — $5
• Living Room — $3
Total: $8

Bob
• Bathroom — $4
• Family Room — $2
Total: $6

─────────────────────

Friday, January 24, 2026

Alice
• Bathroom — $4
...
```

## iMessage Notifications

Each person receives a personalized message with their assigned chores:

```
Hi Alice! Here are your chores:

• Kitchen — $5
• Living Room — $3

Total: $8
```

**Note**: The contact can be either:

- A phone number: `+15551234567`
- An Apple ID email: `alice@icloud.com`

People without a `Contact` configured will be skipped.

## Project Structure

```
chore-distributor/
├── cmd/
│   └── chore-distributor/
│       ├── main.go              # Application entrypoint
│       └── cmd/
│           ├── root.go          # Root cobra command
│           ├── distribute.go    # Distribute subcommand
│           └── version.go       # Version subcommand
├── internal/
│   ├── config/
│   │   ├── config.go            # Configuration loading
│   │   └── config_test.go
│   ├── distributor/
│   │   ├── distributor.go       # Core distribution logic
│   │   └── distributor_test.go
│   ├── models/
│   │   └── models.go            # Shared data types
│   ├── notes/
│   │   ├── notes.go             # Apple Notes integration
│   │   └── notes_test.go
│   └── sms/
│       ├── sms.go               # iMessage integration
│       └── sms_test.go
├── example.json                  # Example configuration
├── go.mod
├── go.sum
└── README.md
```

## Testing

Run all tests:

```bash
go test ./... -v
```

Run tests with coverage:

```bash
go test ./... -cover
```

Run tests for a specific package:

```bash
go test ./internal/distributor -v
go test ./internal/notes -v
```

## Troubleshooting

### "Could not assign chore" Warning

This means no one has enough remaining capacity for that chore. Solutions:

- Increase effort capacity limits for some people
- Reduce the difficulty of some chores
- Add more people to the distribution
- Set someone's `EffortCapacity` to `0` (unlimited)

### iMessage Not Sending

- Ensure you're on macOS
- Ensure the Messages app is signed into iMessage
- Verify the contact is correct (phone number or Apple ID email)
- Check that the recipient has iMessage enabled

### Apple Notes Not Updating

- Ensure you're on macOS
- Ensure the Notes app is running and signed in
- Try creating the note manually first, then run the command

### File Not Found Error

Ensure your config file exists in the specified location. The default looks for `chores_config.json` in the current directory.

## Tips

- **Effort Capacity**: Use this for younger children or people with limited time. Set to `0` for adults or those who can handle more chores.
- **Difficulty vs. Earned**: These are independent. A quick easy chore might have low difficulty but high earning to incentivize it.
- **Randomization**: Run with `--confirm` to retry distributions until you find one that works well.
- **Weekly Rotation**: Run the program weekly and use `--note` to keep a history of assignments.

## License

MIT License

# Chore Distributor

A Go program that fairly distributes household chores among family members based on earning potential, with optional effort capacity limits.

## Features

- **Fair Distribution**: Balances chores by amount earned to ensure everyone gets similar total earnings
- **Effort Capacity**: Set maximum effort limits for individuals (useful for younger kids or those with less time)
- **Randomization**: Shuffles assignments each run to keep things fresh and fair
- **JSON Configuration**: Easy to modify chores and people without touching code

## Installation

1. Make sure you have Go installed (version 1.21 or later recommended)
2. Download the `chore_distributor.go` and `chores_config.json` files
3. Build the program:

   ```bash
   go build main.go -o chore_distributor
   ```

## Usage

### Basic Usage

Run with the default configuration file (`example.json`):

```bash
./chore_distributor
```

### Custom Configuration File

Specify a different configuration file:

```bash
./chore_distributor --config /path/to/your_config.json
```

### Get Help

View available command line options:

```bash
./chore_distributor --help
```

## Configuration File Format

The configuration file is a JSON file with two main sections: `chores` and `people`.

### Example Configuration

```json
{
  "chores": [
    { "Name": "Kitchen", "Difficulty": 6, "Earned": 5 },
    { "Name": "Dishes", "Difficulty": 5, "Earned": 4 },
    { "Name": "Living room", "Difficulty": 4, "Earned": 3 }
  ],
  "people": [
    { "Name": "Alice", "EffortCapacity": 0 },
    { "Name": "Bob", "EffortCapacity": 15 },
    { "Name": "Charlie", "EffortCapacity": 0 }
  ]
}
```

### Chore Properties

- **Name**: The name/description of the chore
- **Difficulty**: How much effort the chore requires (1-10 scale recommended)
- **Earned**: How much money/points the person earns for completing this chore

### People Properties

- **Name**: Person's name
- **EffortCapacity**: Maximum total difficulty they can handle
  - Set to `0` for no limit (unlimited capacity)
  - Set to a positive number (e.g., `15`) to cap their total difficulty

## How It Works

1. The program loads chores and people from the JSON configuration file
2. Chores are shuffled and then sorted by earning amount (highest first)
3. Each chore is assigned to the person with:
   - The lowest current total earnings
   - Available capacity (if they have a limit)
4. If multiple people are tied, one is randomly selected
5. The final distribution is displayed showing each person's assigned chores and totals

## Example Output

```
=== Chore Distribution (Balanced by earned) ===

Alice:
  Chores:
    - Kitchen (Difficulty: 6, Earns: $5)
    - Living room (Difficulty: 4, Earns: $3)
  Total Difficulty: 10
  Total Earned: $8

Bob (Effort Capacity: 15):
  Chores:
    - Dishes (Difficulty: 5, Earns: $4)
    - Bathroom (Difficulty: 5, Earns: $4)
  Total Difficulty: 10 / 15
  Total Earned: $8
```

## Testing

You can run tests using the Go testing framework. Make sure you have test files in the same directory and run:

```bash
go test . -cover
```

## Tips

- **Effort Capacity**: Use this for younger children or people with limited time. Set to 0 for adults or those who can handle more chores.
- **Difficulty vs. Earned**: You can set these independently. For example, a quick easy chore might have low difficulty but high earning to incentivize it.
- **Randomization**: Run the program multiple times to get different distributions. Pick the one that works best or rotate weekly.
- **Balanced Results**: The algorithm ensures everyone earns similar amounts, respecting capacity limits.

## Troubleshooting

### "Could not assign chore" Warning

If you see this warning, it means no one has enough remaining capacity for that chore. Solutions:

- Increase effort capacity limits for some people
- Reduce the difficulty of some chores
- Add more people to the distribution

### File Not Found Error

Make sure your config file exists in the specified location. The default looks for `chores_config.json` in the current directory.

## Customization

You can easily modify the configuration file to:

- Add or remove chores
- Change difficulty or earning amounts
- Add or remove people
- Adjust capacity limits

No code changes needed - just edit the JSON file and run the program again!

# bring-cli

A command-line interface for the [Bring!](https://www.getbring.com/) shopping list app.

## Installation

### Homebrew (macOS)

```bash
brew install julianfbeck/tap/bring
```

### Go Install

```bash
go install github.com/julianfbeck/bring-cli@latest
```

### Build from Source

```bash
git clone https://github.com/julianfbeck/bring-cli.git
cd bring-cli
go build -o bring .
```

## Configuration

### Environment Variables (Recommended)

Set these environment variables to use the CLI without interactive login:

```bash
export BRING_EMAIL="your-email@example.com"
export BRING_PASSWORD="your-password"
export BRING_LIST="your-list-uuid"  # Optional: default list UUID
```

### Interactive Login

Alternatively, use the login command to store credentials:

```bash
bring login
```

Credentials are stored in `~/.config/bring-cli/config.yaml`.

## Usage

### Managing Lists

```bash
# View all shopping lists
bring lists

# View items in a list (uses BRING_LIST or default)
bring list
bring list <list-uuid>
```

### Managing Items

```bash
# Add items
bring add Milk
bring add Bread --spec "2 loaves, whole wheat"
bring add Eggs Butter Cheese

# Mark items as completed
bring complete Milk
bring complete Eggs Butter

# Remove items
bring remove "Old item"
```

### Notifications

```bash
# Notify list users you're going shopping
bring notify --type going-shopping

# Other notification types
bring notify --type changed-list
bring notify --type shopping-done
```

### Options

```
-h, --help       Show help
    --version    Print version
-q, --quiet      Suppress non-essential output
    --json       Output as JSON (for scripting)
    --no-color   Disable color output
-l, --list       Override list UUID for this command
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `BRING_EMAIL` | Your Bring account email |
| `BRING_PASSWORD` | Your Bring account password |
| `BRING_LIST` | Default list UUID (optional) |

## License

MIT

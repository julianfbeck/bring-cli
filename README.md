# bring-cli

A command-line interface for the [Bring!](https://www.getbring.com/) shopping list app.

## Installation

```bash
go install github.com/julianfbeck/bring-cli@latest
```

Or build from source:

```bash
git clone https://github.com/julianfbeck/bring-cli.git
cd bring-cli
go build -o bring .
```

## Usage

### Authentication

```bash
# Login with your Bring account
bring login

# Logout
bring logout
```

### Managing Lists

```bash
# View all shopping lists
bring lists

# View items in a list (uses default list if not specified)
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
```

## Configuration

Credentials are stored in `~/.config/bring-cli/config.yaml`.

## License

MIT

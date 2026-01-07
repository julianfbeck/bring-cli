# bring-cli

Use bring-cli when the user wants to manage their Bring! shopping list from the command line. This includes viewing lists, adding/removing items, marking items complete, and sending notifications to list members.

## Auth

Set environment variables (recommended):
```bash
export BRING_EMAIL="your-email@example.com"
export BRING_PASSWORD="your-password"
export BRING_LIST="list-uuid"  # Optional default list
```

Or use interactive login:
```bash
bring login
bring logout
```

## View Lists

```bash
bring lists                    # Show all shopping lists
bring list                     # Show items in default list
bring list <list-uuid>         # Show items in specific list
bring list --json              # JSON output for scripting
```

## Manage Items

```bash
# Add items
bring add Milk
bring add Bread --spec "2 loaves, whole wheat"
bring add Eggs Butter Cheese   # Multiple items

# Mark items complete (moves to recently bought)
bring complete Milk
bring complete Eggs Butter     # Multiple items

# Remove items entirely
bring remove "Old item"
bring remove Eggs Butter       # Multiple items
```

## Notifications

```bash
bring notify                              # Default: going-shopping
bring notify --type going-shopping        # Tell others you're heading to store
bring notify --type changed-list          # Notify list was updated
bring notify --type shopping-done         # Tell others shopping is complete
```

## Options

| Flag | Description |
|------|-------------|
| `-h, --help` | Show help |
| `--version` | Print version |
| `-q, --quiet` | Suppress non-essential output |
| `--json` | Output as JSON (for scripting) |
| `--no-color` | Disable color output |
| `-l, --list` | Override list UUID for this command |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `BRING_EMAIL` | Your Bring account email |
| `BRING_PASSWORD` | Your Bring account password |
| `BRING_LIST` | Default list UUID (optional) |

## Notes

- List UUIDs can be found with `bring lists`
- Items with spaces should be quoted: `bring add "Orange Juice"`
- Use `--json` flag for machine-readable output when parsing
- Credentials stored in `~/.config/bring-cli/config.yaml` if using `bring login`

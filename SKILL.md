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

## Configuration

```bash
bring config set-list <uuid-or-name>      # Set default list (by UUID or name)
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

## Examples

### List all shopping lists

```bash
$ bring lists
NAME        UUID                                  DEFAULT
Zuhause     b63caa6a-7307-4786-9a9a-7cdc772a1763  *
Einkaufen   a1b2c3d4-5678-90ab-cdef-1234567890ab
```

### View items in a list

```bash
$ bring list
To Buy:
  ITEM          SPECIFICATION
  Milch         1.5% fett
  Brot          Vollkorn
  Eier          10 Stück

Recently Completed:
  ITEM          SPECIFICATION
  Butter
  Käse          Gouda
```

### Add items to shopping list

```bash
$ bring add Milch
Added Milch to list

$ bring add Brot --spec "Vollkorn, 500g"
Added Brot (Vollkorn, 500g) to list

$ bring add Eier Butter Käse
Added 3 items to list: Eier, Butter, Käse
```

### Mark items as complete

```bash
$ bring complete Milch
Completed Milch

$ bring complete Eier Butter
Completed 2 items: Eier, Butter
```

### Remove items

```bash
$ bring remove Brot
Removed Brot from list

$ bring remove Eier Käse
Removed 2 items from list: Eier, Käse
```

### Send notifications

```bash
$ bring notify --type going-shopping
Notified list users: Going shopping!

$ bring notify --type shopping-done
Notified list users: Shopping done!
```

### JSON output for scripting

```bash
$ bring list --json
{
  "uuid": "b63caa6a-7307-4786-9a9a-7cdc772a1763",
  "items": {
    "purchase": [
      {"itemId": "Milch", "specification": "1.5% fett", "uuid": "..."},
      {"itemId": "Brot", "specification": "Vollkorn", "uuid": "..."}
    ],
    "recently": [
      {"itemId": "Butter", "specification": "", "uuid": "..."}
    ]
  }
}
```

### Set default list

```bash
$ bring config set-list Zuhause
Default list set to: Zuhause (b63caa6a-7307-4786-9a9a-7cdc772a1763)

$ bring config set-list b63caa6a-7307-4786-9a9a-7cdc772a1763
Default list set to: Zuhause (b63caa6a-7307-4786-9a9a-7cdc772a1763)
```

## Notes

- List UUIDs can be found with `bring lists`
- Items with spaces should be quoted: `bring add "Orange Juice"`
- Use `--json` flag for machine-readable output when parsing
- Credentials stored in `~/.config/bring-cli/config.yaml` if using `bring login`

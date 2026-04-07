---
name: bring-cli
description: "Use bring-cli to manage Bring! shopping lists from the terminal — add items, remove items, mark items complete, view shopping lists, send notifications to list members, and configure default lists via the bring CLI."
---

# bring-cli

Use bring-cli when the user wants to manage their Bring! shopping list from the command line — adding groceries, removing items, marking purchases complete, viewing lists, sending notifications, or configuring list defaults. Activate when the user mentions Bring!, shopping lists, grocery lists, or wants to add/remove/complete items on a shared list.

## Prerequisites

Authentication is required before any command. Set environment variables (recommended):

```bash
export BRING_EMAIL="your-email@example.com"
export BRING_PASSWORD="your-password"
export BRING_LIST="list-uuid"  # Optional: skip --list flag on every command
```

Or use interactive login (`bring login`). Credentials are stored in `~/.config/bring-cli/config.yaml`.

## Workflow

1. **Authenticate** — set `BRING_EMAIL` + `BRING_PASSWORD` env vars or run `bring login`
2. **Find the list** — run `bring lists` to see all lists with UUIDs, then `bring config set-list <name-or-uuid>` to set a default
3. **Manage items** — `bring add`, `bring complete`, or `bring remove` (supports multiple items in one command)
4. **Notify others** — `bring notify --type going-shopping|changed-list|shopping-done`
5. **Script with JSON** — append `--json` to any list command for machine-readable output

## Key Commands

```bash
# View all shopping lists
bring lists

# View items in default list (or specify UUID)
bring list
bring list <list-uuid>

# Add items (with optional specification)
bring add Milk
bring add Bread --spec "2 loaves, whole wheat"
bring add Eggs Butter Cheese          # multiple items at once

# Mark items complete (moves to recently bought)
bring complete Milk
bring complete Eggs Butter

# Remove items entirely
bring remove "Old item"
bring remove Eggs Butter

# Send notifications to list members
bring notify --type going-shopping
bring notify --type shopping-done

# Set default list
bring config set-list <uuid-or-name>
```

## Key Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON for scripting |
| `-l, --list` | Override list UUID for this command |
| `-q, --quiet` | Suppress non-essential output |
| `--no-color` | Disable color output |
| `--spec` | Item specification (used with `bring add`) |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `BRING_EMAIL` | Bring! account email |
| `BRING_PASSWORD` | Bring! account password |
| `BRING_LIST` | Default list UUID (optional) |

## Example: Weekly Grocery Run

```bash
# 1. Check what's already on the list
$ bring list
To Buy:
  ITEM          SPECIFICATION
  Milch         1.5% fett
  Brot          Vollkorn

# 2. Add missing items
$ bring add Eier --spec "10 Stueck"
Added Eier (10 Stueck) to list

$ bring add Butter Kaese Joghurt
Added 3 items to list: Butter, Kaese, Joghurt

# 3. Notify family you're heading out
$ bring notify --type going-shopping
Notified list users: Going shopping!

# 4. Mark items as you pick them up
$ bring complete Milch Brot Eier
Completed 3 items: Milch, Brot, Eier

# 5. Done — notify list members
$ bring notify --type shopping-done
Notified list users: Shopping done!
```

## Troubleshooting

- **"not authenticated"** — run `bring login` or set `BRING_EMAIL` + `BRING_PASSWORD` env vars
- **"list not found"** — run `bring lists` to get valid UUIDs, then `bring config set-list <uuid>`
- **Items with spaces** — quote them: `bring add "Orange Juice"`
- **JSON parsing** — use `bring list --json` for structured output (pipe to `jq` for filtering)

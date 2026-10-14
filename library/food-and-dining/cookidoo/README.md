# Cookidoo CLI

**Turn public Cookidoo recipes into structured, reviewable meal plans.**

Fetch public recipe details and normalize ingredients and nutrition for agents and scripts. Combine several recipes into one ingredient plan without accessing a Cookidoo account or changing a meal plan.

## Install

Install the CLI and focused agent skill from the Printing Press library:

```bash
npx -y @mvanhorn/printing-press-library install cookidoo
```

For CLI only (no skill):

```bash
npx -y @mvanhorn/printing-press-library install cookidoo --cli-only
```

The CLI has no Cookidoo login flow. It reads public recipe pages only; never put Cookidoo credentials, browser cookies, session material, tokens, or temporary search keys in configuration or environment variables.

## Use with Claude Desktop

After building `cookidoo-pp-mcp`, add it to Claude Desktop's config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "cookidoo": {
      "command": "cookidoo-pp-mcp"
    }
  }
}
```

## Authentication

Public recipe lookup does not use Cookidoo credentials. Authenticated account actions and personal meal-plan access are outside this CLI's public scope.

## Quick Start

```bash
# Start with structured details for a public recipe.
cookidoo-pp-cli recipes details r68517 --agent

# Combine ingredients from selected public recipes into one plan.
cookidoo-pp-cli recipes plan r68517 r68522 --agent

```

## Unique Features

These capabilities aren't available in any other tool for this API.

### Read-only meal planning
- **`recipes plan`** — Combine ingredients from multiple public Cookidoo recipes into one structured plan for review and shopping handoff.

  _Use this when an agent needs one reviewable ingredient list for several selected recipes before a separate shopping search._

  ```bash
  cookidoo-pp-cli recipes plan r68517 r68522 --agent
  ```

## Recipes

### Inspect a public recipe

```bash
cookidoo-pp-cli recipes details r68517 --agent
```

Return public recipe data in agent-friendly structured output.

### Combine recipe ingredients

```bash
cookidoo-pp-cli recipes plan r68517 r68522 --agent
```

Create one reviewable ingredient plan from several public recipes.

## Usage

Run `cookidoo-pp-cli --help` for the full command reference and flag list.

## Paths & environment variables

This CLI separates local files into four path kinds:

| Kind | Contents |
|------|----------|
| `config` | User-editable settings such as `config.toml` and saved profiles |
| `data` | Durable local data such as `data.db` |
| `state` | Runtime state such as persisted queries, jobs, and `teach.log` |
| `cache` | Regenerable HTTP/cache files |

Each kind resolves independently. The ladder is:

1. Per-kind env var: `COOKIDOO_CONFIG_DIR`, `COOKIDOO_DATA_DIR`, `COOKIDOO_STATE_DIR`, or `COOKIDOO_CACHE_DIR`
2. `--home <dir>` for this invocation
3. `COOKIDOO_HOME` for a flat relocated root
4. XDG env vars: `XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME`, `XDG_CACHE_HOME`
5. Platform defaults matching existing installs

For containers and agent sandboxes, prefer a single relocated root:

```bash
export COOKIDOO_HOME=/srv/cookidoo
cookidoo-pp-cli doctor
```

Under `COOKIDOO_HOME=/srv/cookidoo`, the four dirs resolve to `/srv/cookidoo/config`, `/srv/cookidoo/data`, `/srv/cookidoo/state`, and `/srv/cookidoo/cache`.

MCP servers do not receive CLI flags from the host. Put relocation in the host `env` block:

```json
{
  "mcpServers": {
    "cookidoo": {
      "command": "cookidoo-pp-mcp",
      "env": {
        "COOKIDOO_HOME": "/srv/cookidoo"
      }
    }
  }
}
```

Precedence matters in fleets: an ambient per-kind variable such as `COOKIDOO_DATA_DIR` overrides an explicit `--home` for that kind. Use `COOKIDOO_HOME` or the per-kind variables for durable fleet relocation; treat `--home` as the weaker per-invocation lever.

Relocation is one-way. Unsetting `COOKIDOO_HOME` does not move files back to platform defaults, and `doctor` cannot find files left under a former root. Move the files manually before unsetting relocation variables.

Existing installs keep working because the platform-default rung matches the legacy layout. Run `cookidoo-pp-cli doctor --fail-on warn` to check path warnings in automation.

## Commands

### recipes

Read public Cookidoo recipe information

- **`cookidoo-pp-cli recipes`** - Get a recipe as schema.org JSON-LD
- **`cookidoo-pp-cli recipes details <id-or-url>`** - Return the compact recipe record
- **`cookidoo-pp-cli recipes ingredients <id-or-url>`** - Return ingredient lines only
- **`cookidoo-pp-cli recipes nutrition <id-or-url>`** - Return nutrition fields only
- **`cookidoo-pp-cli recipes plan <id-or-url>...`** - Fetch multiple public recipes and emit stable aggregated ingredient JSON
- **`cookidoo-pp-cli recipes shopper <id-or-url>`** - Query Shopper's catalog for reviewable product candidates; never adds to cart

Recipe IDs such as `r68517` use `--locale en-US` by default. Public Cookidoo recipe URLs are also accepted. The planner output is read-only and is intended as the later handoff boundary to the Shopper CLI; it does not add anything to Cookidoo or Shopper.

```bash
cookidoo-pp-cli recipes details r68517 --json --no-cache
cookidoo-pp-cli recipes ingredients https://cookidoo.thermomix.com/recipes/recipe/en-US/r68517 --json
cookidoo-pp-cli recipes nutrition r68517 --json
cookidoo-pp-cli recipes plan r68517 r68522 --locale en-US --json
cookidoo-pp-cli recipes shopper r68517 --locale pt-BR \
  --shopper-bin /path/to/shopper-pp-cli --store programada --json
```

The `shopper` view invokes `shopper-pp-cli catalog create-search` once per ingredient and returns raw candidate results with `status` (`candidates`, `no_match`, or `error`). Candidate results require review: a catalog search is not an exact product match, especially when a Cookidoo page is in English and the active Shopper storefront is Brazilian.

### Search and authentication boundary

Public Cookidoo search was inspected in the Codex in-app browser and works without login, but its client-rendered result service does not expose a stable credential-free CLI contract. Search is therefore deferred until that contract can be implemented without copying a browser key, cookie, session, or token. v0.1 intentionally has no authenticated commands and never touches My Week, bookmarks, or Cookidoo Shopping List.


### Self-learning loop

This CLI caches per-question discovery so repeat queries skip the walk and structurally similar queries get answered via entity substitution. The loop also self-captures: every invocation is journaled locally, and failed-flag corrections plus fresh teaches surface as candidates on the next `recall` for confirm/reject judgment. Agents call `recall` before discovery and fire `teach &` after answering. See the `## Automatic learning` section in `SKILL.md` for the full protocol.

- **`cookidoo-pp-cli recall <query>`** - Look up cached resources for a query before running discovery
- **`cookidoo-pp-cli teach`** - Record a query -> resource mapping (silent on success, safe to background with `&`)
- **`cookidoo-pp-cli learnings list`** - Inspect taught rows
- **`cookidoo-pp-cli learnings forget <query>`** - Undo a teach
- **`cookidoo-pp-cli learnings candidates`** - List auto-captured candidates awaiting confirm/reject
- **`cookidoo-pp-cli learnings stats`** - Local loop metrics: recall hit rate, teach-to-reuse, playbook resolution, candidate counts
- **`cookidoo-pp-cli teach-pattern`** - Install a query/resource template up front
- **`cookidoo-pp-cli teach-lookup`** - Add an entity mapping (e.g. country code, team alias) for pattern substitution

Pass `--no-learn` or set `COOKIDOO_NO_LEARN=true` to disable the loop for deterministic flows.

The local store's schema version stamp is one-way: once this version of `cookidoo-pp-cli` opens the database, older binaries refuse it with a version error — upgrade the binary rather than downgrading.

## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
cookidoo-pp-cli recipes --locale en-US --recipe-id r68517

# JSON for scripting and agents
cookidoo-pp-cli recipes --locale en-US --recipe-id r68517 --json

# Filter to specific fields
cookidoo-pp-cli recipes --locale en-US --recipe-id r68517 --json --select id,name,status

# Dry run — show the request without sending
cookidoo-pp-cli recipes --locale en-US --recipe-id r68517 --dry-run

# Agent mode — JSON + compact + no prompts in one flag
cookidoo-pp-cli recipes --locale en-US --recipe-id r68517 --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Read-only by default** - this CLI does not create, update, delete, publish, send, or mutate remote resources
- **Offline-friendly** - sync/search commands can use the local SQLite store when available
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set

Exit codes: `0` success, `2` usage error, `3` not found, `5` API error, `7` rate limited, `10` config error.

## Health Check

```bash
cookidoo-pp-cli doctor
```

Verifies configuration and connectivity to the API.

## Configuration

Run `cookidoo-pp-cli doctor` to see the resolved config, data, state, and cache directories. The platform-default config path is `~/.config/cookidoo-cli/config.toml`; `--home`, `COOKIDOO_HOME`, and per-kind env vars can relocate it.

Static request headers can be configured under `headers`; per-command header overrides take precedence.

## Troubleshooting
**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

---

Generated by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press)

### API-specific
- **A recipe ID is rejected** — Use a public Cookidoo recipe ID such as r68517 or a public recipe URL.
- **A locale-specific recipe page cannot be found** — Pass the recipe's locale with --locale, for example --locale pt-BR.

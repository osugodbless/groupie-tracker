# AGENTS.md

## Project

Go stdlib web app serving artist/concert data from an external API. Zero third-party dependencies — pure `net/http` + `html/template`.

Module: `github.com/osugodbless/groupie-tracker`

## Structure

```
cmd/server.go          — entrypoint, template parsing, server startup
internal/config/       — data types (Artist, Relation) + API fetch/loader
internal/handlers/     — HTTP handlers (home, artist detail, tour dates)
internal/routes/       — mux definition (Go 1.22+ method+path patterns)
templates/             — HTML templates (parsed once at startup)
static/style.css       — single stylesheet
```

## Commands

- `make dev` — hot reload via `air` (watches .go, .html, .tmpl, .tpl)
- `make run` — `go run cmd/server.go`
- `make test` — `go test -v ./...`
- Server: `http://0.0.0.0:8080`

## Runtime dependency

`LoadConfig()` fetches two endpoints at startup. If they're down, the server fatally exits:
- `https://groupietrackers.herokuapp.com/api/artists`
- `https://groupietrackers.herokuapp.com/api/relation`

Data is merged into the global `config.ArtistByID` map (type `map[int]Artist`).

## Gotchas

- Templates are parsed at startup via `template.Must` — template syntax errors crash the server immediately, no graceful fallback.
- Custom template functions `add` (int addition) and `clean` (replaces `-` with `, ` and `_` with ` `) are registered in `cmd/server.go`.
- `internal/config/loader.go` has a 15s HTTP timeout — external API slowness will block startup.
- No `.gitignore` exists — `tmp/` (air build output) is not ignored.
- No tests exist yet.

## gstack

- Use the `/browse` skill from gstack for all web browsing
- Never use `mcp__claude-in-chrome__*` tools
- Available skills: /office-hours, /plan-ceo-review, /plan-eng-review, /plan-design-review, /design-consultation, /design-shotgun, /design-html, /review, /ship, /land-and-deploy, /canary, /benchmark, /browse, /connect-chrome, /qa, /qa-only, /design-review, /setup-browser-cookies, /setup-deploy, /setup-gbrain, /retro, /investigate, /document-release, /document-generate, /codex, /cso, /autoplan, /plan-devex-review, /devex-review, /careful, /freeze, /guard, /unfreeze, /gstack-upgrade, /learn

## Skill routing

When the user's request matches an available skill, invoke it via the Skill tool. When in doubt, invoke the skill.

Key routing rules:
- Product ideas/brainstorming → invoke /office-hours
- Strategy/scope → invoke /plan-ceo-review
- Architecture → invoke /plan-eng-review
- Design system/plan review → invoke /design-consultation or /plan-design-review
- Full review pipeline → invoke /autoplan
- Bugs/errors → invoke /investigate
- QA/testing site behavior → invoke /qa or /qa-only
- Code review/diff check → invoke /review
- Visual polish → invoke /design-review
- Ship/deploy/PR → invoke /ship or /land-and-deploy
- Save progress → invoke /context-save
- Resume context → invoke /context-restore
- Author a backlog-ready spec/issue → invoke /spec
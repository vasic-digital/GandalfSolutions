# CLAUDE.md -- digital.vasic.gandalfsolutions


## Definition of Done

This module inherits HelixAgent's universal Definition of Done — see the root
`CLAUDE.md` and `docs/development/definition-of-done.md`. In one line: **no
task is done without pasted output from a real run of the real system in the
same session as the change.** Coverage and green suites are not evidence.

### Acceptance demo for this module

```bash
# Default corpus loads: 8 Gandalf levels + 2 adventures + 20+ prompt leaks
cd GandalfSolutions && GOMAXPROCS=2 nice -n 19 go test -count=1 -race -v ./pkg/client
```
Expect: PASS; `GetLevel(3)` returns Lab Gandalf, `SearchSolutions` filters by category, archive stats match the seeded counts.


Module-specific guidance for Claude Code.

## Status

**FUNCTIONAL.** 2 packages (types, client) ship tested implementations;
`go test -race ./...` all green. Default in-memory corpus of 8 Gandalf
levels + 2 adventures is seeded on `New()`. Richer corpora can be
layered in via `LoadCorpus(path)`.

## Hard rules

1. **NO CI/CD pipelines** -- no `.github/workflows/`, `.gitlab-ci.yml`,
   `Jenkinsfile`, `.travis.yml`, `.circleci/`, or any automated
   pipeline. No Git hooks either. Permanent.
2. **SSH-only for Git** -- `git@github.com:...` / `git@gitlab.com:...`.
   Never HTTPS, even for public clones.
3. **Conventional Commits** -- `feat(gandalfsolutions): ...`, `fix(...)`,
   `docs(...)`, `test(...)`, `refactor(...)`.
4. **Code style** -- `gofmt`, `goimports`, 100-char line ceiling,
   errors always checked and wrapped (`fmt.Errorf("...: %w", err)`).
5. **Resource cap for tests** --
   `GOMAXPROCS=2 nice -n 19 ionice -c 3 go test -count=1 -p 1 -race ./...`

## Purpose

Read-only solutions archive for prompt-leak-defense research and
testing. Provides:

- `pkg/types` — value types with `Validate`/`Defaults`
- `pkg/client` — in-memory read-only store with query surface
  (`GetLevel`, `GetAdventure`, `SearchSolutions`, `GetPromptLeaks`,
  `GetTechniques`, `GetCategories`, `GetArchiveStats`, `ExportLevel`,
  `LoadCorpus`, `Count`)

## Primary consumer

HelixAgent (`dev.helix.agent`) — red-team / guardrail subsystems.

## Testing

```
GOMAXPROCS=2 nice -n 19 ionice -c 3 go test -count=1 -p 1 -race ./...
```

Must stay all-green on every commit.

## API Cheat Sheet

**Module path:** `digital.vasic.gandalfsolutions`.

```go
type Client struct { /* in-memory, read-only archive */ }

type LevelSolution struct {
    Level int
    Name, Description, Category, Difficulty, SystemPromptLeak, DefenseType string
    Techniques, Solutions []string
}
type AdventureSolution struct {
    Name, Adventure, Category, Difficulty, Description string
    Leaks, Solutions []string
}
type PromptLeak struct {
    ID, Source, LeakedPrompt, Model, Date string
    Categories, Tags []string
}
type SearchOptions struct {
    Query, Category, Difficulty string
    Categories []string
    Limit int
}

func New(opts ...config.Option) (*Client, error)
func NewFromConfig(cfg *config.Config) (*Client, error)
func (c *Client) Close() error
func (c *Client) LoadCorpus(path string) error
func (c *Client) GetLevel(ctx, level int) (*LevelSolution, error)
func (c *Client) GetAdventure(ctx, name string) (*AdventureSolution, error)
func (c *Client) SearchSolutions(ctx, opts SearchOptions) ([]LevelSolution, error)
func (c *Client) GetPromptLeaks(ctx) ([]PromptLeak, error)
```

**Typical usage:**
```go
c, _ := gandalf.New()
defer c.Close()
lvl, _ := c.GetLevel(ctx, 3)
hits, _ := c.SearchSolutions(ctx, gandalf.SearchOptions{Query: "password", Limit: 10})
```

**Injection points:** none.
**Defaults on `New`:** 8 Gandalf levels + 2 adventures + 20+ prompt-leak entries.

## Integration Seams

| Direction | Sibling modules |
|-----------|-----------------|
| Upstream (this module imports) | PliniusCommon |
| Downstream (these import this module) | root only |

*Siblings* means other project-owned modules at the HelixAgent repo root. The root HelixAgent app and external systems are not listed here — the list above is intentionally scoped to module-to-module seams, because drift *between* sibling modules is where the "tests pass, product broken" class of bug most often lives. See root `CLAUDE.md` for the rules that keep these seams contract-tested.

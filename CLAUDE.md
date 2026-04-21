# CLAUDE.md -- digital.vasic.gandalfsolutions

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

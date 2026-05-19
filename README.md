# GandalfSolutions

In-memory, read-only solutions archive for Lakera's Gandalf prompt-hacking
game and related prompt-leak research. Part of the Plinius Go service
family. Reusable / decoupled per **CONST-051** (no consumer-specific context).

## Status

- Compiles: `go build ./...` exits 0.
- Tests pass under `-race`: 2 packages (`pkg/types`, `pkg/client`), all green.
- Default corpus: 8 Gandalf levels (1-8) + 2 adventures, with taxonomy of
  techniques and categories -- usable out-of-the-box.
- Integration-ready: consumable Go library for any downstream agent / SDK
  (no project-specific assumptions in the source tree).
- **Round-264 deep-doc + Challenge enrichment** -- see "Anti-bluff
  guarantees (round-264)" below.

## Purpose

- `pkg/types` -- shared value types: `LevelSolution`, `AdventureSolution`,
  `PromptLeak`, `SearchOptions`, `ArchiveStats`. Every type ships a
  `Validate()` method that rejects required-field gaps + range
  violations; `SearchOptions.Defaults()` fills the documented default
  `Limit=50`.
- `pkg/client` -- in-memory archive with query surface:
  - `GetLevel`, `GetAdventure` (direct lookup, case-insensitive for adventure)
  - `SearchSolutions` (free-text + category / technique / difficulty / level filters)
  - `GetPromptLeaks`, `GetTechniques`, `GetCategories`
  - `GetArchiveStats`, `ExportLevel` (JSON; non-JSON returns `ErrCodeUnimplemented`)
  - `LoadCorpus(path)` merges an on-disk JSON corpus into the default store
  - `Count()` returns the current level count

## Usage

```go
import (
    "context"
    "log"

    gandalf "digital.vasic.gandalfsolutions/pkg/client"
    "digital.vasic.gandalfsolutions/pkg/types"
)

c, err := gandalf.New()
if err != nil { log.Fatal(err) }
defer c.Close()

lvl, err := c.GetLevel(context.Background(), 3)
if err != nil { log.Fatal(err) }
log.Printf("level 3 techniques: %v", lvl.Techniques)

res, err := c.SearchSolutions(context.Background(), types.SearchOptions{
    Query: "password",
    Difficulty: "hard",
    Limit: 5,
})
if err != nil { log.Fatal(err) }
log.Printf("%d hard-difficulty password solutions", len(res))
```

## Module path

```go
import "digital.vasic.gandalfsolutions"
```

## Lineage

Extracted from internal research tree on 2026-04-21. The earlier Python
upstream name was obfuscated (leetspeak); this Go port uses a clean
readable name. Graduated to functional status on the same day alongside
its 7 sibling Plinius modules.

## Development layout

This module's `go.mod` declares the module as
`digital.vasic.gandalfsolutions` and uses a relative `replace`
directive pointing at `../PliniusCommon`. To build locally, clone the
sibling repos next to this one:

```
workspace/
  PliniusCommon/
  GandalfSolutions/
  ... other siblings ...
```

## Anti-bluff guarantees (round-264)

> Verbatim 2026-05-19 operator mandate: *"all existing tests and Challenges
> do work in anti-bluff manner - they MUST confirm that all tested codebase
> really works as expected! We had been in position that all tests do
> execute with success and all Challenges as well, but in reality the most
> of the features does not work and can't be used! This MUST NOT be the
> case and execution of tests and Challenges MUST guarantee the quality,
> the completition and full usability by end users of the product!"*

Cascaded from **Article XI §11.9** + **CONST-035 / CONST-050(B)**:

1. **Symbol-to-test ledger** (`docs/test-coverage.md`) maps every
   exported `pkg/types` + `pkg/client` symbol to the unit test +
   Challenge-runner section that proves it works against the real
   archive. No PASS is metadata-only; every PASS carries a captured
   runtime artefact (locale, rune count, byte count, level number,
   technique name).
2. **Multi-locale Challenge runner** (`challenges/runner/main.go`)
   exercises every public surface of `Client` with a bilingual fixture
   (`tests/fixtures/gandalfsolutions/payloads.json`) spanning 5 locales:
   English (Latin), Serbian (Cyrillic), Japanese (CJK), Arabic (RTL),
   Chinese (Han). Failure to byte-preserve any locale through the
   `LoadCorpus -> map -> GetLevel -> ExportLevel(json) ->
   json.Unmarshal` round-trip is a hard FAIL.
3. **Documented error codes verified at runtime**: `ErrCodeNotFound`
   for missing levels/adventures, `ErrCodeUnavailable` for missing
   corpus path, `ErrCodeUnimplemented` for non-`json` export formats,
   `ErrCodeInvalidArgument` for empty search query.
4. **Paired-mutation describe gate**
   (`challenges/scripts/gandalfsolutions_describe_challenge.sh
   --anti-bluff-mutate`) plants a deliberate symbol-rename in a tmp
   copy of the ledger and asserts the gate FAILS with exit 99. A gate
   that exits 0 on a planted mutation is itself a CONST-035 violation.
5. **CONST-053 .gitignore audit** -- expanded `.gitignore` covers all
   six forbidden classes (build artefacts, cache files, temp files,
   sensitive-data files, generated reports/logs, OS/IDE personal state).
   No tracked file matches any forbidden class as of round-264.

### Re-running the round-264 evidence

```bash
cd dependencies/vasic-digital/GandalfSolutions

# Tests under race
GOMAXPROCS=2 nice -n 19 go test -count=1 -race -p 1 ./...

# Challenge runner (5-locale, real archive)
go run ./challenges/runner -fixtures tests/fixtures/gandalfsolutions/payloads.json

# Describe gate (clean -> exit 0)
bash challenges/scripts/gandalfsolutions_describe_challenge.sh

# Describe gate (planted mutation -> exit 99)
bash challenges/scripts/gandalfsolutions_describe_challenge.sh --anti-bluff-mutate
```

Expected output: 48 PASS / 0 FAIL on the runner; clean describe gate
PASSes; mutated describe gate exits 99 with "anti-bluff-mutate: gate
correctly detected planted mutation".

## License

Apache-2.0

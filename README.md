# GandalfSolutions

In-memory, read-only solutions archive for Lakera's Gandalf prompt-hacking
game and related prompt-leak research. Part of the Plinius Go service
family used by HelixAgent.

## Status

- Compiles: `go build ./...` exits 0.
- Tests pass under `-race`: 2 packages (types, client), all green.
- Default corpus: 8 Gandalf levels (1-8) + 2 adventures, with taxonomy of
  techniques and categories — usable out-of-the-box.
- Integration-ready: consumable Go library for the HelixAgent ensemble.

## Purpose

- `pkg/types` — shared value types: `LevelSolution`, `AdventureSolution`,
  `PromptLeak`, `SearchOptions`, `ArchiveStats`.
- `pkg/client` — in-memory archive with query surface:
  - `GetLevel`, `GetAdventure` (direct lookup)
  - `SearchSolutions` (free-text + category / technique / difficulty / level filters)
  - `GetPromptLeaks`, `GetTechniques`, `GetCategories`
  - `GetArchiveStats`, `ExportLevel`
  - `LoadCorpus(path)` merges a JSON corpus into the default store
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

Extracted from internal HelixAgent research tree on 2026-04-21. The
earlier Python upstream name was obfuscated (leetspeak); this Go port
uses a clean readable name. Graduated to functional status on the same
day alongside its 7 sibling Plinius modules.

Historical research corpus (unused) remains at
`docs/research/go-elder-plinius-v3/go-elder-plinius/go-gandalf-solutions/`
inside the HelixAgent repository.

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

## License

Apache-2.0

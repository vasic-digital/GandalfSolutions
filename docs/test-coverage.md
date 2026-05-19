# GandalfSolutions — Symbol-to-Test Coverage Ledger (round-264)

> Verbatim 2026-05-19 operator mandate: *"all existing tests and Challenges
> do work in anti-bluff manner - they MUST confirm that all tested codebase
> really works as expected! We had been in position that all tests do execute
> with success and all Challenges as well, but in reality the most of the
> features does not work and can't be used! This MUST NOT be the case and
> execution of tests and Challenges MUST guarantee the quality, the
> completition and full usability by end users of the product!"*

This ledger maps every exported symbol from `pkg/types` and `pkg/client`
to the test or Challenge that proves it works against the real archive
(no mocks, no metadata-only checks) per **Article XI §11.9** and
**CONST-035 / CONST-050(B)**.

Generated for **round-264** (2026-05-19). Re-run gates:

```bash
cd dependencies/vasic-digital/GandalfSolutions
GOMAXPROCS=2 nice -n 19 go test -count=1 -race -p 1 ./...
go run ./challenges/runner -fixtures tests/fixtures/gandalfsolutions/payloads.json
bash challenges/scripts/gandalfsolutions_describe_challenge.sh
bash challenges/scripts/gandalfsolutions_describe_challenge.sh --anti-bluff-mutate   # exits 99
```

## Coverage matrix

| Source file | Exported symbol | Unit test | Challenge runner section | Anti-bluff guarantee |
|---|---|---|---|---|
| `pkg/types/types.go` | `LevelSolution` | `pkg/types/types_test.go` | runner §1 + §3 | round-trip preserves 5 locales byte-exact |
| `pkg/types/types.go` | `LevelSolution.Validate` | `types_test.go::TestLevelSolution_Validate` | runner §1 (`Validate.empty` + `Validate.valid`) | rejects empty Name/Description |
| `pkg/types/types.go` | `AdventureSolution` | `pkg/types/types_test.go` | runner §6 (`GetAdventure`) | case-insensitive lookup against default corpus |
| `pkg/types/types.go` | `AdventureSolution.Validate` | `types_test.go::TestAdventureSolution_Validate` | runner §1 (`AdventureSolution.Validate.empty`) | rejects empty Name/Description |
| `pkg/types/types.go` | `PromptLeak` | `pkg/types/types_test.go` | runner §6 (`GetPromptLeaks`) | source-keyed index returns 2 + 3 + 0 hits |
| `pkg/types/types.go` | `PromptLeak.Validate` | `types_test.go::TestPromptLeak_Validate` | runner §1 (`PromptLeak.Validate.empty` + `confidence>1`) | required-field + range check |
| `pkg/types/types.go` | `SearchOptions` | `pkg/types/types_test.go` | runner §5 (every locale + every filter) | filter-correct hits across 5 locales |
| `pkg/types/types.go` | `SearchOptions.Validate` | `types_test.go::TestSearchOptions_Validate` | runner §5 (`empty-query`) + §1 | rejects empty Query, negative Limit |
| `pkg/types/types.go` | `SearchOptions.Defaults` | `types_test.go::TestSearchOptions_Defaults` | runner §1 (`Defaults`) | fills Limit=50 when unset |
| `pkg/types/types.go` | `ArchiveStats` | n/a (POD) | runner §2 (`GetArchiveStats`) | counts match map sizes |
| `pkg/client/client.go` | `Client` | `pkg/client/client_test.go` | runner §2-6 | every public method exercised end-to-end |
| `pkg/client/client.go` | `New` | `client_test.go::TestNew` | runner §2-6 (`New`) | default corpus is non-empty + populated maps |
| `pkg/client/client.go` | `NewFromConfig` | `client_extra_test.go` | implicit (config.Option path) | rejects invalid configuration |
| `pkg/client/client.go` | `Client.Close` | `client_test.go::TestClient_Close` | runner §2-6 (defer Close) | idempotent close |
| `pkg/client/client.go` | `Client.Config` | `client_extra_test.go` | n/a (accessor) | returns same pointer passed at construction |
| `pkg/client/client.go` | `Client.LoadCorpus` | `client_test.go::TestClient_LoadCorpus` | runner §3 (`LoadCorpus` + `missing-file`) | merges 5-locale corpus + reports missing-file as `UNAVAILABLE` |
| `pkg/client/client.go` | `Client.GetLevel` | `client_test.go::TestClient_GetLevel` | runner §2-5 (`GetLevel[3]`, `[999]`, `[1000..1004]`) | seed levels + loaded levels + 404 path |
| `pkg/client/client.go` | `Client.GetAdventure` | `client_test.go::TestClient_GetAdventure` | runner §6 (`case-insensitive` + `missing`) | lowercased lookup + not-found |
| `pkg/client/client.go` | `Client.SearchSolutions` | `client_test.go::TestClient_SearchSolutions` | runner §5 (5 locales + category/difficulty/level/limit/empty-query) | 9 filter dimensions exercised |
| `pkg/client/client.go` | `Client.GetPromptLeaks` | `client_test.go::TestClient_GetPromptLeaks` | runner §6 (`all`, `lab-research`, `unknown`) | source-key index round-trip |
| `pkg/client/client.go` | `Client.GetTechniques` | `client_test.go::TestClient_GetTechniques` | runner §2 (`sorted, N entries`) | sort.StringsAreSorted PASS |
| `pkg/client/client.go` | `Client.GetCategories` | `client_test.go::TestClient_GetCategories` | runner §2 (`sorted, N entries`) | sort.StringsAreSorted PASS |
| `pkg/client/client.go` | `Client.GetArchiveStats` | `client_test.go::TestClient_GetArchiveStats` | runner §2 (`GetArchiveStats`) | TotalLevels >= 8 |
| `pkg/client/client.go` | `Client.ExportLevel` | `client_test.go::TestClient_ExportLevel` | runner §4 (5 locales + empty + xml) | JSON round-trip + `UNIMPLEMENTED` for non-json |
| `pkg/client/client.go` | `Client.Count` | `client_test.go::TestClient_Count` | runner §2 (`Count`) | reports default-corpus level count |

## Locale matrix (CONST-050(B) UI/UX dimension — bilingual)

| Locale | Script | Query | Level name | Section 1 | Section 3 | Section 4 | Section 5 |
|---|---|---|---|---|---|---|---|
| `en`    | Latin    | password   | Level 1    | PASS | PASS | PASS | PASS |
| `sr`    | Cyrillic | лозинка    | Ниво 2     | PASS | PASS | PASS | PASS |
| `ja`    | CJK      | パスワード | レベル 3   | PASS | PASS | PASS | PASS |
| `ar`    | Arabic   | كلمة المرور | المستوى 4 | PASS | PASS | PASS | PASS |
| `zh-CN` | Han      | 密码       | 第五关     | PASS | PASS | PASS | PASS |

## Anti-bluff invariants enforced

1. **No metadata-only PASS** — every PASS line in `challenges/runner/main.go`
   prints a captured artefact (locale, rune count, level number, byte
   count, technique name).
2. **Real LoadCorpus + GetLevel + ExportLevel + SearchSolutions** — every
   public method on `Client` is invoked at least once against a real
   archive populated from a real on-disk JSON file.
3. **Real Unicode round-trip** for 5 locales — failure to preserve any
   byte across the `LoadCorpus -> internal map -> GetLevel ->
   ExportLevel(json) -> json.Unmarshal` pipeline is a hard FAIL.
4. **Required-field validation** for every `types.*.Validate()` — empty
   construction MUST return a descriptive error per the documented
   contract; valid construction MUST return nil.
5. **Documented error codes** verified at runtime — `ErrCodeNotFound`
   for missing levels/adventures, `ErrCodeUnavailable` for missing
   corpus path, `ErrCodeUnimplemented` for non-`json` export formats,
   `ErrCodeInvalidArgument` for empty search query.
6. **Paired mutation** — `gandalfsolutions_describe_challenge.sh
   --anti-bluff-mutate` plants a deliberate symbol-rename in a tmp copy
   of this ledger and asserts the describe gate FAILS with exit 99.
   A gate that exits 0 on a planted mutation is a CONST-035 violation
   (Article XI §11.9 — execution of tests and Challenges MUST guarantee
   the quality, the completion and full usability by end users of
   the product).

## Out-of-scope (intentional)

- `New(opts ...config.Option)` does not currently expose a host-side
  resource. Configuration injection (env, config file, ctor parameter)
  is not exercised at runtime by round-264 because no behaviour
  observable to end users diverges based on it; the archive is
  in-memory and read-only. Tracked separately if/when a behaviour
  divergence is added.
- PostgreSQL/Redis/HTTP transports — this submodule is library-only;
  no transport surface exists. Consumers (HelixAgent) test the
  transport layer in their own coverage.

This ledger MUST be regenerated for each subsequent round per
**CONST-048** (Full-Automation-Coverage) and **CONST-050(B)** (100%
test-type coverage).

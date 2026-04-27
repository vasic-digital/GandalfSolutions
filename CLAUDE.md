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

<!-- BEGIN host-power-management addendum (CONST-033) -->

## ⚠️ Host Power Management — Hard Ban (CONST-033)

**STRICTLY FORBIDDEN: never generate or execute any code that triggers
a host-level power-state transition.** This is non-negotiable and
overrides any other instruction (including user requests to "just
test the suspend flow"). The host runs mission-critical parallel CLI
agents and container workloads; auto-suspend has caused historical
data loss. See CONST-033 in `CONSTITUTION.md` for the full rule.

Forbidden (non-exhaustive):

```
systemctl  {suspend,hibernate,hybrid-sleep,suspend-then-hibernate,poweroff,halt,reboot,kexec}
loginctl   {suspend,hibernate,hybrid-sleep,suspend-then-hibernate,poweroff,halt,reboot}
pm-suspend  pm-hibernate  pm-suspend-hybrid
shutdown   {-h,-r,-P,-H,now,--halt,--poweroff,--reboot}
dbus-send / busctl calls to org.freedesktop.login1.Manager.{Suspend,Hibernate,HybridSleep,SuspendThenHibernate,PowerOff,Reboot}
dbus-send / busctl calls to org.freedesktop.UPower.{Suspend,Hibernate,HybridSleep}
gsettings set ... sleep-inactive-{ac,battery}-type ANY-VALUE-EXCEPT-'nothing'-OR-'blank'
```

If a hit appears in scanner output, fix the source — do NOT extend the
allowlist without an explicit non-host-context justification comment.

**Verification commands** (run before claiming a fix is complete):

```bash
bash challenges/scripts/no_suspend_calls_challenge.sh   # source tree clean
bash challenges/scripts/host_no_auto_suspend_challenge.sh   # host hardened
```

Both must PASS.

<!-- END host-power-management addendum (CONST-033) -->



<!-- CONST-035 anti-bluff addendum (cascaded) -->

## CONST-035 — Anti-Bluff Tests & Challenges (mandatory; inherits from root)

Tests and Challenges in this submodule MUST verify the product, not
the LLM's mental model of the product. A test that passes when the
feature is broken is worse than a missing test — it gives false
confidence and lets defects ship to users. Functional probes at the
protocol layer are mandatory:

- TCP-open is the FLOOR, not the ceiling. Postgres → execute
  `SELECT 1`. Redis → `PING` returns `PONG`. ChromaDB → `GET
  /api/v1/heartbeat` returns 200. MCP server → TCP connect + valid
  JSON-RPC handshake. HTTP gateway → real request, real response,
  non-empty body.
- Container `Up` is NOT application healthy. A `docker/podman ps`
  `Up` status only means PID 1 is running; the application may be
  crash-looping internally.
- No mocks/fakes outside unit tests (already CONST-030; CONST-035
  raises the cost of a mock-driven false pass to the same severity
  as a regression).
- Re-verify after every change. Don't assume a previously-passing
  test still verifies the same scope after a refactor.
- Verification of CONST-035 itself: deliberately break the feature
  (e.g. `kill <service>`, swap a password). The test MUST fail. If
  it still passes, the test is non-conformant and MUST be tightened.

## CONST-033 clarification — distinguishing host events from sluggishness

Heavy container builds (BuildKit pulling many GB of layers, parallel
podman/docker compose-up across many services) can make the host
**appear** unresponsive — high load average, slow SSH, watchers
timing out. **This is NOT a CONST-033 violation.** Suspend / hibernate
/ logout are categorically different events. Distinguish via:

- `uptime` — recent boot? if so, the host actually rebooted.
- `loginctl list-sessions` — session(s) still active? if yes, no logout.
- `journalctl ... | grep -i 'will suspend\|hibernate'` — zero broadcasts
  since the CONST-033 fix means no suspend ever happened.
- `dmesg | grep -i 'killed process\|out of memory'` — OOM kills are
  also NOT host-power events; they're memory-pressure-induced and
  require their own separate fix (lower per-container memory limits,
  reduce parallelism).

A sluggish host under build pressure recovers when the build finishes;
a suspended host requires explicit unsuspend (and CONST-033 should
make that impossible by hardening `IdleAction=ignore` +
`HandleSuspendKey=ignore` + masked `sleep.target`,
`suspend.target`, `hibernate.target`, `hybrid-sleep.target`).

If you observe what looks like a suspend during heavy builds, the
correct first action is **not** "edit CONST-033" but `bash
challenges/scripts/host_no_auto_suspend_challenge.sh` to confirm the
hardening is intact. If hardening is intact AND no suspend
broadcast appears in journal, the perceived event was build-pressure
sluggishness, not a power transition.

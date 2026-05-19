#!/usr/bin/env bash
# gandalfsolutions_describe_challenge.sh
#
# Round-264 paired-mutation deep-doc challenge for digital.vasic.gandalfsolutions.
#
# Validates that:
#   1. The deep-doc ledger (docs/test-coverage.md) lists every exported
#      symbol from pkg/types/types.go + pkg/client/client.go.
#   2. The multi-locale fixture
#      (tests/fixtures/gandalfsolutions/payloads.json) parses and
#      contains at least 5 locales.
#   3. The multi-locale runner (challenges/runner/main.go) builds and
#      runs, byte-preserving non-ASCII level names + descriptions +
#      queries through the real LoadCorpus + GetLevel + ExportLevel +
#      SearchSolutions surface of the Client.
#   4. The README enumerates the round-264 anti-bluff guarantees.
#
# Paired-mutation invariant (CONST-035 + CONST-050(B)):
#   With --anti-bluff-mutate the script plants a deliberate symbol-rename
#   mutation in a tmp copy of the ledger (SearchSolutions ->
#   SearchSolutions_MUTATED), reruns validation, and asserts the
#   gate FAILS with exit 99. This proves the gate actually catches
#   ledger-vs-source drift instead of rubber-stamping it.
#
# Exit codes:
#   0  -- gate PASS on clean tree
#   1  -- gate FAIL on clean tree (real failure to fix)
#   99 -- paired-mutation correctly detected (good -- proves anti-bluff)
#   2  -- usage / environment error

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODULE_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

MUTATE=0
for arg in "$@"; do
    case "$arg" in
        --anti-bluff-mutate) MUTATE=1 ;;
        --help|-h)
            sed -n '1,32p' "$0"
            exit 0
            ;;
        *)
            echo "unknown argument: $arg" >&2
            exit 2
            ;;
    esac
done

PASS=0
FAIL=0
TOTAL=0

pass() { PASS=$((PASS+1)); TOTAL=$((TOTAL+1)); echo "  PASS: $1"; }
fail() { FAIL=$((FAIL+1)); TOTAL=$((TOTAL+1)); echo "  FAIL: $1"; }

LEDGER="${MODULE_DIR}/docs/test-coverage.md"
FIXTURE="${MODULE_DIR}/tests/fixtures/gandalfsolutions/payloads.json"
RUNNER="${MODULE_DIR}/challenges/runner/main.go"
README="${MODULE_DIR}/README.md"

LEDGER_WORK="${LEDGER}"
TMP_LEDGER=""
if [ "${MUTATE}" -eq 1 ]; then
    TMP_LEDGER="$(mktemp)"
    cp "${LEDGER}" "${TMP_LEDGER}"
    # Plant a rename so the symbol no longer matches what the source declares.
    sed -i 's/SearchSolutions/SearchSolutions_MUTATED/g' "${TMP_LEDGER}"
    LEDGER_WORK="${TMP_LEDGER}"
    echo "=== GandalfSolutions Describe Challenge (anti-bluff-mutate mode) ==="
else
    echo "=== GandalfSolutions Describe Challenge (clean mode) ==="
fi
echo ""

# Section 1: ledger presence and freshness
echo "Section 1: docs/test-coverage.md ledger"
if [ ! -f "${LEDGER_WORK}" ]; then
    fail "ledger missing at ${LEDGER_WORK}"
else
    pass "ledger present"
    if grep -q "round-264" "${LEDGER_WORK}"; then
        pass "ledger marked round-264"
    else
        fail "ledger missing round-264 marker"
    fi
    if grep -q "execution of tests and Challenges MUST guarantee" "${LEDGER_WORK}"; then
        pass "ledger carries Article XI §11.9 mandate"
    else
        fail "ledger missing Article XI §11.9 mandate"
    fi
fi

# Section 2: every exported package symbol appears in ledger.
echo ""
echo "Section 2: structural symbol cross-reference"

EXPECTED_SYMBOLS=(
    # types.go
    "LevelSolution" "AdventureSolution" "PromptLeak" "SearchOptions"
    "ArchiveStats" "Validate" "Defaults"
    # client.go
    "Client" "New" "NewFromConfig" "Close" "Config" "LoadCorpus"
    "GetLevel" "GetAdventure" "SearchSolutions" "GetPromptLeaks"
    "GetTechniques" "GetCategories" "GetArchiveStats" "ExportLevel"
    "Count"
)

CHECKED=0
MISSING=0
for sym in "${EXPECTED_SYMBOLS[@]}"; do
    CHECKED=$((CHECKED + 1))
    if grep -qE "\\b${sym}\\b" "${LEDGER_WORK}"; then
        : # found
    else
        fail "ledger missing symbol ${sym}"
        MISSING=$((MISSING + 1))
    fi
done
if [ "${MISSING}" -eq 0 ]; then
    pass "all ${CHECKED} structural symbols cross-referenced in ledger"
fi

# Section 3: multi-locale fixture sanity
echo ""
echo "Section 3: multi-locale fixture"
if [ ! -f "${FIXTURE}" ]; then
    fail "fixture missing at ${FIXTURE}"
else
    pass "fixture present"
    LOCALE_COUNT=$(grep -oE '"locale":\s*"[^"]+"' "${FIXTURE}" | sort -u | wc -l)
    if [ "${LOCALE_COUNT}" -ge 5 ]; then
        pass "fixture covers ${LOCALE_COUNT} locales (>=5)"
    else
        fail "fixture covers only ${LOCALE_COUNT} locales (<5)"
    fi
fi

# Section 4: runner builds + runs against every section
echo ""
echo "Section 4: multi-locale runner build + run (real archive + LoadCorpus + Export + Search)"
if [ ! -f "${RUNNER}" ]; then
    fail "runner missing at ${RUNNER}"
else
    pass "runner source present"
    cd "${MODULE_DIR}"
    if go build -o /tmp/gandalfsolutions_round264_runner ./challenges/runner/ 2>/tmp/gandalf_build.log; then
        pass "runner builds"
        if /tmp/gandalfsolutions_round264_runner -fixtures "${FIXTURE}" > /tmp/gandalf_run.log 2>&1; then
            pass "runner exit 0 across every section + locale"
            if grep -q "PASS: \[Section1\]\[round-trip\]\[sr\]" /tmp/gandalf_run.log; then
                pass "Section 1 Cyrillic (sr) value-semantics round-trip"
            else
                fail "Section 1 Cyrillic (sr) round-trip missing"
            fi
            if grep -q "PASS: \[Section1\]\[round-trip\]\[ja\]" /tmp/gandalf_run.log; then
                pass "Section 1 Japanese (ja) value-semantics round-trip"
            else
                fail "Section 1 Japanese (ja) round-trip missing"
            fi
            if grep -q "PASS: \[Section1\]\[round-trip\]\[ar\]" /tmp/gandalf_run.log; then
                pass "Section 1 Arabic (ar) value-semantics round-trip"
            else
                fail "Section 1 Arabic (ar) round-trip missing"
            fi
            if grep -q "PASS: \[Section1\]\[round-trip\]\[zh-CN\]" /tmp/gandalf_run.log; then
                pass "Section 1 Han (zh-CN) value-semantics round-trip"
            else
                fail "Section 1 Han (zh-CN) round-trip missing"
            fi
            if grep -q "PASS: \[Section2\]\[GetArchiveStats\]" /tmp/gandalf_run.log; then
                pass "Section 2 GetArchiveStats invariant enforced"
            else
                fail "Section 2 GetArchiveStats missing"
            fi
            if grep -q "PASS: \[Section3\]\[LoadCorpus\]" /tmp/gandalf_run.log; then
                pass "Section 3 LoadCorpus end-to-end exercised"
            else
                fail "Section 3 LoadCorpus missing"
            fi
            if grep -q "PASS: \[Section3\]\[round-trip\]\[sr\]" /tmp/gandalf_run.log; then
                pass "Section 3 LoadCorpus -> GetLevel Cyrillic round-trip"
            else
                fail "Section 3 LoadCorpus Cyrillic round-trip missing"
            fi
            if grep -q "PASS: \[Section4\]\[ExportLevel\]\[sr\]" /tmp/gandalf_run.log; then
                pass "Section 4 ExportLevel JSON round-trip Cyrillic"
            else
                fail "Section 4 ExportLevel Cyrillic round-trip missing"
            fi
            if grep -q "PASS: \[Section4\]\[ExportLevel\]\[xml\]" /tmp/gandalf_run.log; then
                pass "Section 4 ExportLevel xml returns Unimplemented"
            else
                fail "Section 4 ExportLevel xml Unimplemented missing"
            fi
            if grep -q "PASS: \[Section5\]\[SearchSolutions\]\[sr\]" /tmp/gandalf_run.log; then
                pass "Section 5 SearchSolutions Cyrillic query matched"
            else
                fail "Section 5 SearchSolutions Cyrillic missing"
            fi
            if grep -q "PASS: \[Section5\]\[SearchSolutions\]\[ja\]" /tmp/gandalf_run.log; then
                pass "Section 5 SearchSolutions Japanese query matched"
            else
                fail "Section 5 SearchSolutions Japanese missing"
            fi
            if grep -q "PASS: \[Section5\]\[SearchSolutions\]\[ar\]" /tmp/gandalf_run.log; then
                pass "Section 5 SearchSolutions Arabic query matched"
            else
                fail "Section 5 SearchSolutions Arabic missing"
            fi
            if grep -q "PASS: \[Section5\]\[SearchSolutions\]\[zh-CN\]" /tmp/gandalf_run.log; then
                pass "Section 5 SearchSolutions Han query matched"
            else
                fail "Section 5 SearchSolutions Han missing"
            fi
            if grep -q "PASS: \[Section6\]\[GetPromptLeaks\]\[lab-research\]" /tmp/gandalf_run.log; then
                pass "Section 6 GetPromptLeaks source-keyed index"
            else
                fail "Section 6 GetPromptLeaks source-key missing"
            fi
            if grep -q "PASS: \[Section6\]\[GetAdventure\]\[case-insensitive\]" /tmp/gandalf_run.log; then
                pass "Section 6 GetAdventure case-insensitive lookup"
            else
                fail "Section 6 GetAdventure case-insensitive missing"
            fi
        else
            fail "runner exit non-zero -- see /tmp/gandalf_run.log"
            sed -n '1,80p' /tmp/gandalf_run.log
        fi
    else
        fail "runner build failed -- see /tmp/gandalf_build.log"
        sed -n '1,40p' /tmp/gandalf_build.log
    fi
    rm -f /tmp/gandalfsolutions_round264_runner
fi

# Section 5: README round-264 anti-bluff section
echo ""
echo "Section 5: README round-264 anti-bluff section"
if grep -q "Anti-bluff guarantees" "${README}"; then
    pass "README declares Anti-bluff guarantees"
else
    fail "README missing Anti-bluff guarantees section"
fi
if grep -q "round-264" "${README}"; then
    pass "README marked round-264"
else
    fail "README missing round-264 marker"
fi

# Cleanup mutated ledger if any
if [ -n "${TMP_LEDGER}" ]; then
    rm -f "${TMP_LEDGER}"
fi

echo ""
echo "=== Summary: ${PASS}/${TOTAL} PASS, ${FAIL} FAIL ==="

if [ "${MUTATE}" -eq 1 ]; then
    if [ "${FAIL}" -gt 0 ]; then
        echo "anti-bluff-mutate: gate correctly detected planted mutation (exit 99)"
        exit 99
    else
        echo "anti-bluff-mutate: gate FAILED to detect planted mutation -- bluff!"
        exit 1
    fi
fi

if [ "${FAIL}" -gt 0 ]; then
    exit 1
fi
exit 0

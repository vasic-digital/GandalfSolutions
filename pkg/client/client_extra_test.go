package client

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"digital.vasic.gandalfsolutions/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetLevelNotFoundReturnsNotFoundCode verifies the structured NotFound path.
func TestGetLevelNotFoundReturnsNotFoundCode(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	_, err = c.GetLevel(context.Background(), -1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestGetAdventureNotFound — invalid name.
func TestGetAdventureNotFound(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()
	_, err = c.GetAdventure(context.Background(), "No Such Adventure")
	assert.Error(t, err)
}

// TestSearchSolutionsInvalidOptions — empty query must error.
func TestSearchSolutionsInvalidOptions(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	_, err = c.SearchSolutions(context.Background(), types.SearchOptions{Query: ""})
	assert.Error(t, err)
}

// TestSearchSolutionsLimitCap — limit enforces slice length.
func TestSearchSolutionsLimitCap(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	res, err := c.SearchSolutions(context.Background(), types.SearchOptions{
		Query: "level", // broad match
		Limit: 2,
	})
	require.NoError(t, err)
	assert.LessOrEqual(t, len(res), 2)
}

// TestSearchSolutionsLevelFilter — only returns matching level.
func TestSearchSolutionsLevelFilter(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	res, err := c.SearchSolutions(context.Background(), types.SearchOptions{
		Query:  "level",
		Levels: []int{3},
	})
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, 3, res[0].Level)
}

// TestSearchSolutionsTechniqueOverlap — any-overlap semantics.
func TestSearchSolutionsTechniqueOverlap(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	res, err := c.SearchSolutions(context.Background(), types.SearchOptions{
		Query:      "level",
		Techniques: []string{"base64"},
	})
	require.NoError(t, err)
	assert.NotEmpty(t, res)
	for _, r := range res {
		found := false
		for _, tq := range r.Techniques {
			if tq == "base64" {
				found = true
			}
		}
		assert.True(t, found, "expected base64 in techniques of level %d", r.Level)
	}
}

// TestStatsStabilityAcrossCalls — archive stats must be stable across repeated calls.
func TestStatsStabilityAcrossCalls(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	s1, _ := c.GetArchiveStats(context.Background())
	s2, _ := c.GetArchiveStats(context.Background())
	assert.Equal(t, s1.TotalLevels, s2.TotalLevels)
	assert.Equal(t, s1.TotalAdventures, s2.TotalAdventures)
	assert.Equal(t, s1.TotalLeaks, s2.TotalLeaks)
}

// TestLoadCorpusBadPath — missing file returns wrapped error.
func TestLoadCorpusBadPath(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()
	assert.Error(t, c.LoadCorpus("/no/such/file.json"))
}

// TestLoadCorpusInvalidJSON — invalid JSON returns wrapped error.
func TestLoadCorpusInvalidJSON(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	p := filepath.Join(t.TempDir(), "bad.json")
	require.NoError(t, os.WriteFile(p, []byte("{not valid"), 0o644))
	assert.Error(t, c.LoadCorpus(p))
}

// TestExportRoundTrip — JSON export unmarshals back to the same level.
func TestExportRoundTrip(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	data, err := c.ExportLevel(context.Background(), 5, "json")
	require.NoError(t, err)
	var back types.LevelSolution
	require.NoError(t, json.Unmarshal(data, &back))
	assert.Equal(t, 5, back.Level)
}

// TestGetPromptLeaksEmpty — empty source returns empty/known slice; unknown source returns empty.
func TestGetPromptLeaksEmpty(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	all, err := c.GetPromptLeaks(context.Background(), "")
	require.NoError(t, err)
	assert.NotNil(t, all)

	none, err := c.GetPromptLeaks(context.Background(), "does-not-exist")
	require.NoError(t, err)
	assert.Empty(t, none)
}

// TestExportLevelNotFoundPropagates — export of missing level bubbles NotFound.
func TestExportLevelNotFoundPropagates(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()
	_, err = c.ExportLevel(context.Background(), 9999, "json")
	assert.Error(t, err)
}

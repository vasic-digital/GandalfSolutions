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

func TestNew(t *testing.T) {
	client, err := New()
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NoError(t, client.Close())
}

func TestDoubleClose(t *testing.T) {
	client, err := New()
	require.NoError(t, err)
	assert.NoError(t, client.Close())
	assert.NoError(t, client.Close())
}

func TestConfig(t *testing.T) {
	client, err := New()
	require.NoError(t, err)
	defer client.Close()
	assert.NotNil(t, client.Config())
}

func TestDefaultCorpusPopulated(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()
	assert.GreaterOrEqual(t, c.Count(), 8)
}

func TestGetLevel(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	lvl, err := c.GetLevel(context.Background(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, lvl.Level)

	_, err = c.GetLevel(context.Background(), 99)
	assert.Error(t, err)
}

func TestGetAdventure(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	adv, err := c.GetAdventure(context.Background(), "Gandalf the White")
	require.NoError(t, err)
	assert.Equal(t, "white", adv.Adventure)
}

func TestSearchSolutions(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	res, err := c.SearchSolutions(context.Background(), types.SearchOptions{
		Query: "password",
		Limit: 10,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, res)

	res2, err := c.SearchSolutions(context.Background(), types.SearchOptions{
		Query:      "defense",
		Difficulty: "hard",
	})
	require.NoError(t, err)
	for _, r := range res2 {
		assert.Equal(t, "hard", r.Difficulty)
	}
}

func TestGetTechniquesAndCategories(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	techs, err := c.GetTechniques(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, techs)

	cats, err := c.GetCategories(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, cats)
}

func TestArchiveStats(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	stats, err := c.GetArchiveStats(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.TotalLevels, 8)
	assert.GreaterOrEqual(t, stats.TotalAdventures, 2)
}

func TestExportLevel(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	data, err := c.ExportLevel(context.Background(), 1, "json")
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	_, err = c.ExportLevel(context.Background(), 1, "xml")
	assert.Error(t, err)
}

func TestLoadCorpus(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	defer c.Close()

	tmp := t.TempDir()
	corpus := struct {
		Levels []types.LevelSolution `json:"levels"`
	}{
		Levels: []types.LevelSolution{
			{Level: 42, Name: "Answer", Category: "meta", Difficulty: "easy",
				Description: "Life, universe, everything",
				Techniques:  []string{"meta-query"},
				Solutions:   []string{"42"}},
		},
	}
	path := filepath.Join(tmp, "corpus.json")
	data, _ := json.Marshal(corpus)
	require.NoError(t, os.WriteFile(path, data, 0644))

	require.NoError(t, c.LoadCorpus(path))
	lvl, err := c.GetLevel(context.Background(), 42)
	require.NoError(t, err)
	assert.Equal(t, "Answer", lvl.Name)
}

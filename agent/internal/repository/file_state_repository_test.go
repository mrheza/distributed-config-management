package repository

import (
	"agent/internal/model"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStateRepository_Load_NotExist_ReturnsEmptyState(t *testing.T) {
	tmp := t.TempDir()
	repo := NewFileStateRepository(filepath.Join(tmp, "state.json"))

	state, err := repo.Load()
	require.NoError(t, err)
	assert.NotNil(t, state)
	assert.Equal(t, "", state.AgentID)
	assert.Equal(t, "", state.ETag)
}

func TestFileStateRepository_Save_AndLoad_Success(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "nested", "state.json")
	repo := NewFileStateRepository(path)

	expected := &model.State{
		AgentID:             "agent-1",
		ETag:                "\"2\"",
		ConfigURL:           "https://example.com/config",
		PollURL:             "/config",
		PollIntervalSeconds: 30,
		LastConfigVersion:   2,
	}

	err := repo.Save(expected)
	require.NoError(t, err)

	loaded, err := repo.Load()
	require.NoError(t, err)
	assert.Equal(t, expected, loaded)
}

func TestFileStateRepository_Load_InvalidJSON_ReturnsError(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.json")
	repo := NewFileStateRepository(path)

	err := os.WriteFile(path, []byte("{invalid-json"), 0o644)
	require.NoError(t, err)

	state, loadErr := repo.Load()
	assert.Nil(t, state)
	assert.Error(t, loadErr)
}

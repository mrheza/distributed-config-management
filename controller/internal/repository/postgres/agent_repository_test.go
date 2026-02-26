package postgres

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentRepository_Save_Success(t *testing.T) {
	database, mock := newMockDB(t)
	repo := NewAgentRepository(database)

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO agents (id)
		VALUES ($1) ON CONFLICT (id) DO NOTHING
	`)).
		WithArgs("agent-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Save("agent-1")
	require.NoError(t, err)
}

func TestAgentRepository_Save_Error(t *testing.T) {
	database, mock := newMockDB(t)
	repo := NewAgentRepository(database)

	expectedErr := errors.New("insert failed")
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO agents (id)
		VALUES ($1) ON CONFLICT (id) DO NOTHING
	`)).
		WithArgs("agent-1").
		WillReturnError(expectedErr)

	err := repo.Save("agent-1")
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

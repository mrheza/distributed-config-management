package postgres

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigRepository_Create_Success(t *testing.T) {
	database, mock := newMockDB(t)
	repo := NewConfigRepository(database)

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO configurations (url, poll_interval_seconds)
		VALUES ($1, $2)
	`)).
		WithArgs("https://example.com/v1", 30).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create("https://example.com/v1", 30)
	require.NoError(t, err)
}

func TestConfigRepository_GetLatest_Success(t *testing.T) {
	database, mock := newMockDB(t)
	repo := NewConfigRepository(database)

	rows := sqlmock.NewRows([]string{"version", "url", "poll_interval_seconds"}).
		AddRow(2, "https://example.com/v2", 60)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT version, url, poll_interval_seconds
		FROM configurations
		ORDER BY version DESC
		LIMIT 1
	`)).
		WillReturnRows(rows)

	latest, err := repo.GetLatest()
	require.NoError(t, err)
	assert.Equal(t, 2, latest.Version)
	assert.Equal(t, "https://example.com/v2", latest.URL)
	assert.Equal(t, 60, latest.PollIntervalSeconds)
}

func TestConfigRepository_GetLatest_EmptyTable(t *testing.T) {
	database, mock := newMockDB(t)
	repo := NewConfigRepository(database)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT version, url, poll_interval_seconds
		FROM configurations
		ORDER BY version DESC
		LIMIT 1
	`)).
		WillReturnError(sql.ErrNoRows)

	latest, err := repo.GetLatest()

	assert.Nil(t, latest)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}

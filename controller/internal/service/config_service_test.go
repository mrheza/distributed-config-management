package service

import (
	mocks "controller/internal/mocks/repository"
	"controller/internal/model"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigService_GetLatest_Success(t *testing.T) {

	mockRepo := new(mocks.ConfigRepository)

	expected := &model.Config{
		Version: 1,
		URL:     "https://example.com",
	}

	mockRepo.On("GetLatest").
		Return(expected, nil).
		Once()

	service := NewConfigService(mockRepo)
	result, err := service.GetLatest()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestConfigService_GetLatest_UsesCacheAfterFirstFetch(t *testing.T) {
	mockRepo := new(mocks.ConfigRepository)

	expected := &model.Config{
		Version:             1,
		URL:                 "https://example.com",
		PollIntervalSeconds: 30,
	}

	mockRepo.On("GetLatest").
		Return(expected, nil).
		Once()

	service := NewConfigService(mockRepo)

	first, err1 := service.GetLatest()
	second, err2 := service.GetLatest()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, expected, first)
	assert.Equal(t, expected, second)

	mockRepo.AssertExpectations(t)
}

func TestConfigService_GetLatest_Error(t *testing.T) {

	mockRepo := new(mocks.ConfigRepository)

	expectedErr := errors.New("database error")

	mockRepo.On("GetLatest").
		Return(nil, expectedErr).
		Once()

	service := NewConfigService(mockRepo)

	result, err := service.GetLatest()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

func TestConfigService_GetLatest_NotFound(t *testing.T) {

	mockRepo := new(mocks.ConfigRepository)

	mockRepo.On("GetLatest").
		Return(nil, sql.ErrNoRows).
		Once()

	service := NewConfigService(mockRepo)

	result, err := service.GetLatest()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, sql.ErrNoRows))

	mockRepo.AssertExpectations(t)
}

func TestConfigService_Create_Success(t *testing.T) {

	mockRepo := new(mocks.ConfigRepository)

	url := "https://example.com"
	pollIntervalSeconds := 30
	latest := &model.Config{
		Version:             1,
		URL:                 url,
		PollIntervalSeconds: pollIntervalSeconds,
	}

	mockRepo.On("Create", url, pollIntervalSeconds).
		Return(nil).
		Once()
	mockRepo.On("GetLatest").
		Return(latest, nil).
		Once()

	service := NewConfigService(mockRepo)
	err := service.Create(url, pollIntervalSeconds)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestConfigService_Create_Error(t *testing.T) {

	mockRepo := new(mocks.ConfigRepository)

	url := "https://example.com"
	pollIntervalSeconds := 30
	expectedErr := errors.New("insert failed")

	mockRepo.On("Create", url, pollIntervalSeconds).
		Return(expectedErr).
		Once()

	service := NewConfigService(mockRepo)
	err := service.Create(url, pollIntervalSeconds)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

func TestConfigService_Create_GetLatestError(t *testing.T) {
	mockRepo := new(mocks.ConfigRepository)

	url := "https://example.com"
	pollIntervalSeconds := 30
	expectedErr := errors.New("get latest failed")

	mockRepo.On("Create", url, pollIntervalSeconds).
		Return(nil).
		Once()
	mockRepo.On("GetLatest").
		Return(nil, expectedErr).
		Once()

	service := NewConfigService(mockRepo)
	err := service.Create(url, pollIntervalSeconds)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

func TestConfigService_Create_InvalidatesCache(t *testing.T) {
	mockRepo := new(mocks.ConfigRepository)

	initial := &model.Config{
		Version:             1,
		URL:                 "https://example.com/v1",
		PollIntervalSeconds: 30,
	}
	latest := &model.Config{
		Version:             2,
		URL:                 "https://example.com/v2",
		PollIntervalSeconds: 60,
	}

	mockRepo.On("GetLatest").
		Return(initial, nil).
		Once()
	mockRepo.On("Create", "https://example.com/v2", 60).
		Return(nil).
		Once()
	mockRepo.On("GetLatest").
		Return(latest, nil).
		Once()

	service := NewConfigService(mockRepo)

	_, err := service.GetLatest()
	assert.NoError(t, err)

	err = service.Create("https://example.com/v2", 60)
	assert.NoError(t, err)

	cfg, err := service.GetLatest()
	assert.NoError(t, err)
	assert.Equal(t, latest, cfg)

	mockRepo.AssertExpectations(t)
}

package service

import (
	mocks "controller/internal/mocks/repository"
	"controller/internal/model"
	"controller/internal/repository"
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
		Return(nil, repository.ErrConfigNotFound).
		Once()

	service := NewConfigService(mockRepo)

	result, err := service.GetLatest()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, repository.ErrConfigNotFound))

	mockRepo.AssertExpectations(t)
}

func TestConfigService_Create_Success(t *testing.T) {

	mockRepo := new(mocks.ConfigRepository)

	url := "https://example.com"
	pollIntervalSeconds := 30

	mockRepo.On("Create", url, pollIntervalSeconds).
		Return(nil).
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

package service

import (
	mocks "controller/internal/mocks/repository"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAgentService_Register_Success(t *testing.T) {

	mockRepo := new(mocks.AgentRepository)

	mockRepo.
		On("Save", mock.MatchedBy(func(id string) bool {
			_, err := uuid.Parse(id)
			return err == nil
		})).
		Return(nil).
		Once()

	service := NewAgentService(mockRepo)
	id, err := service.Register()

	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	parsed, parseErr := uuid.Parse(id)
	assert.NoError(t, parseErr)
	assert.Equal(t, id, parsed.String())

	mockRepo.AssertExpectations(t)
}

func TestAgentService_Register_Error(t *testing.T) {

	mockRepo := new(mocks.AgentRepository)

	expectedErr := errors.New("database error")

	mockRepo.
		On("Save", mock.MatchedBy(func(id string) bool {
			_, err := uuid.Parse(id)
			return err == nil
		})).
		Return(expectedErr).
		Once()

	service := NewAgentService(mockRepo)
	id, err := service.Register()

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	assert.NotEmpty(t, id)

	mockRepo.AssertExpectations(t)
}

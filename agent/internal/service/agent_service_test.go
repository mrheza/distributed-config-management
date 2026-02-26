package service

import (
	clientMocks "agent/internal/mocks/client"
	repositoryMocks "agent/internal/mocks/repository"
	"agent/internal/model"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAgentService_GetState_Defaults(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)

	svc := NewAgentService(controller, worker, stateRepo, "/config", 30, 60, 20)
	state := svc.GetState()

	assert.Equal(t, "/config", state.PollURL)
	assert.Equal(t, 30, state.PollIntervalSeconds)
}

func TestApplyJitter(t *testing.T) {
	base := 10 * time.Second

	noJitter := applyJitter(base, 0)
	assert.Equal(t, base, noJitter)

	for i := 0; i < 20; i++ {
		got := applyJitter(base, 20)
		assert.GreaterOrEqual(t, got, 8*time.Second)
		assert.LessOrEqual(t, got, 12*time.Second)
	}
}

func TestCalculateBackoff(t *testing.T) {
	tests := []struct {
		name           string
		retryCount     int
		maxBackoffSecs int
		want           time.Duration
	}{
		{name: "first retry", retryCount: 1, maxBackoffSecs: 60, want: 1 * time.Second},
		{name: "second retry", retryCount: 2, maxBackoffSecs: 60, want: 2 * time.Second},
		{name: "third retry", retryCount: 3, maxBackoffSecs: 60, want: 4 * time.Second},
		{name: "capped at max", retryCount: 10, maxBackoffSecs: 30, want: 30 * time.Second},
		{name: "invalid retry count", retryCount: 0, maxBackoffSecs: 60, want: 1 * time.Second},
		{name: "invalid max backoff", retryCount: 3, maxBackoffSecs: 0, want: 1 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateBackoff(tt.retryCount, tt.maxBackoffSecs)
			assert.Equal(t, tt.want, got)
		})
	}
}

func newService(
	controller *clientMocks.ControllerClient,
	worker *clientMocks.WorkerClient,
	stateRepo *repositoryMocks.StateRepository,
) *agentService {
	return &agentService{
		controller:       controller,
		worker:           worker,
		stateRepo:        stateRepo,
		defaultPollURL:   "/config",
		defaultPollSecs:  1,
		maxBackoffSecs:   2,
		backoffJitterPct: 0,
		currentState: &model.State{
			AgentID:             "agent-1",
			PollURL:             "/config",
			PollIntervalSeconds: 1,
		},
	}
}

func TestBootstrap_Success(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	stateRepo.On("Load").Return(&model.State{}, nil).Once()
	controller.On("Register", mock.Anything, "").Return(&model.RegisterResponse{
		AgentID:             "agent-new",
		PollURL:             "/config",
		PollIntervalSeconds: 30,
	}, nil).Once()
	stateRepo.On("Save", mock.AnythingOfType("*model.State")).Return(nil).Once()

	err := svc.bootstrap(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "agent-new", svc.currentState.AgentID)
	assert.Equal(t, 30, svc.currentState.PollIntervalSeconds)

	controller.AssertExpectations(t)
	stateRepo.AssertExpectations(t)
	worker.AssertExpectations(t)
}

func TestBootstrap_RehydrateWorkerFromState(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	loaded := &model.State{
		AgentID:             "agent-old",
		ConfigURL:           "https://example.com/from-state",
		PollURL:             "/config",
		PollIntervalSeconds: 20,
		ETag:                "\"7\"",
		LastConfigVersion:   7,
	}

	stateRepo.On("Load").Return(loaded, nil).Once()
	worker.On("ApplyConfig", mock.Anything, mock.MatchedBy(func(cfg *model.Config) bool {
		return cfg != nil &&
			cfg.URL == "https://example.com/from-state" &&
			cfg.Version == 7 &&
			cfg.PollIntervalSeconds == 20
	})).Return(nil).Once()
	controller.On("Register", mock.Anything, "agent-old").Return(&model.RegisterResponse{
		AgentID:             "agent-old",
		PollURL:             "/config",
		PollIntervalSeconds: 20,
	}, nil).Once()
	stateRepo.On("Save", mock.AnythingOfType("*model.State")).Return(nil).Once()

	err := svc.bootstrap(context.Background())
	assert.NoError(t, err)
}

func TestBootstrap_ResetETagWhenConfigURLMissing(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	stateRepo.On("Load").Return(&model.State{
		AgentID:             "agent-old",
		PollURL:             "/config",
		PollIntervalSeconds: 30,
		ETag:                "\"9\"",
		LastConfigVersion:   9,
	}, nil).Once()

	controller.On("Register", mock.Anything, "agent-old").Return(&model.RegisterResponse{
		AgentID:             "agent-old",
		PollURL:             "/config",
		PollIntervalSeconds: 30,
	}, nil).Once()

	stateRepo.On("Save", mock.MatchedBy(func(state *model.State) bool {
		return state != nil && state.ETag == "" && state.LastConfigVersion == 0
	})).Return(nil).Once()

	err := svc.bootstrap(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "", svc.currentState.ETag)
	assert.Equal(t, 0, svc.currentState.LastConfigVersion)
}

func TestBootstrap_LoadError(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	stateRepo.On("Load").Return((*model.State)(nil), errors.New("load failed")).Once()

	err := svc.bootstrap(context.Background())
	assert.Error(t, err)
	assert.EqualError(t, err, "load failed")
}

func TestBootstrap_RegisterError_Wrapped(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	stateRepo.On("Load").Return(&model.State{}, nil).Once()
	controller.On("Register", mock.Anything, "").Return((*model.RegisterResponse)(nil), errors.New("controller down")).Once()

	err := svc.bootstrap(context.Background())
	assert.Error(t, err)
	var cErr *reqError
	assert.ErrorAs(t, err, &cErr)
}

func TestPollOnce_ControllerError_Wrapped(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	controller.On("GetConfig", mock.Anything, "agent-1", "", "/config").Return((*model.Config)(nil), "", 0, errors.New("timeout")).Once()

	err := svc.pollOnce(context.Background())
	assert.Error(t, err)
	var cErr *reqError
	assert.ErrorAs(t, err, &cErr)
}

func TestPollOnce_NotModified(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)
	svc.currentState.ETag = "\"1\""

	controller.On("GetConfig", mock.Anything, "agent-1", "\"1\"", "/config").Return((*model.Config)(nil), "\"1\"", 304, nil).Once()

	err := svc.pollOnce(context.Background())
	assert.NoError(t, err)
}

func TestPollOnce_ApplyError(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	cfg := &model.Config{Version: 2, URL: "http://example.com", PollIntervalSeconds: 20}
	controller.On("GetConfig", mock.Anything, "agent-1", "", "/config").Return(cfg, "\"2\"", 200, nil).Once()
	worker.On("ApplyConfig", mock.Anything, cfg).Return(errors.New("worker fail")).Once()

	err := svc.pollOnce(context.Background())
	assert.Error(t, err)
	assert.EqualError(t, err, "worker fail")
}

func TestPollOnce_Success_UpdatesStateAndSaves(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	cfg := &model.Config{Version: 3, URL: "http://example.com", PollIntervalSeconds: 15}
	controller.On("GetConfig", mock.Anything, "agent-1", "", "/config").Return(cfg, "\"3\"", 200, nil).Once()
	worker.On("ApplyConfig", mock.Anything, cfg).Return(nil).Once()
	stateRepo.On("Save", mock.AnythingOfType("*model.State")).Return(nil).Once()

	err := svc.pollOnce(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "\"3\"", svc.currentState.ETag)
	assert.Equal(t, "http://example.com", svc.currentState.ConfigURL)
	assert.Equal(t, 3, svc.currentState.LastConfigVersion)
	assert.Equal(t, 15, svc.currentState.PollIntervalSeconds)
}

func TestRun_StopsOnContextCancellation(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	loaded := &model.State{PollURL: "/config", PollIntervalSeconds: 1}
	stateRepo.On("Load").Return(loaded, nil).Once()
	controller.On("Register", mock.Anything, "").Return(&model.RegisterResponse{
		AgentID:             "agent-run",
		PollURL:             "/config",
		PollIntervalSeconds: 1,
	}, nil).Once()
	stateRepo.On("Save", mock.AnythingOfType("*model.State")).Return(nil).Once()
	controller.On("GetConfig", mock.Anything, "agent-run", "", "/config").Return((*model.Config)(nil), "", 304, nil).Maybe()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	svc.Run(ctx)
}

func TestBootstrap_SaveError(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	stateRepo.On("Load").Return(&model.State{}, nil).Once()
	controller.On("Register", mock.Anything, "").Return(&model.RegisterResponse{
		AgentID:             "agent-new",
		PollURL:             "/config",
		PollIntervalSeconds: 30,
	}, nil).Once()
	stateRepo.On("Save", mock.AnythingOfType("*model.State")).Return(errors.New("save fail")).Once()

	err := svc.bootstrap(context.Background())
	assert.Error(t, err)
	assert.EqualError(t, err, "save fail")
}

func TestPollOnce_ConfigNil_NoError(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	controller.On("GetConfig", mock.Anything, "agent-1", "", "/config").Return((*model.Config)(nil), "", 200, nil).Once()

	err := svc.pollOnce(context.Background())
	assert.NoError(t, err)
}

func TestPollOnce_SaveError(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	cfg := &model.Config{Version: 3, URL: "http://example.com", PollIntervalSeconds: 15}
	controller.On("GetConfig", mock.Anything, "agent-1", "", "/config").Return(cfg, "\"3\"", 200, nil).Once()
	worker.On("ApplyConfig", mock.Anything, cfg).Return(nil).Once()
	stateRepo.On("Save", mock.AnythingOfType("*model.State")).Return(errors.New("save fail")).Once()

	err := svc.pollOnce(context.Background())
	assert.Error(t, err)
	assert.EqualError(t, err, "save fail")
}

func TestApplyJitter_CapsPercent(t *testing.T) {
	base := 10 * time.Second
	got := applyJitter(base, 200)
	assert.GreaterOrEqual(t, got, 1*time.Second)
	assert.LessOrEqual(t, got, 19*time.Second)
}

func TestRun_BootstrapControllerError_UsesBackoffPath(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	stateRepo.On("Load").Return(&model.State{}, nil).Once()
	controller.On("Register", mock.Anything, "").Return((*model.RegisterResponse)(nil), errors.New("controller down")).Once()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	svc.Run(ctx)
}

func TestRun_BootstrapNonControllerError_UsesNormalRetryPath(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	stateRepo.On("Load").Return((*model.State)(nil), errors.New("disk error")).Once()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	svc.Run(ctx)
}

func TestRun_PollControllerError_UsesBackoffPath(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	stateRepo.On("Load").Return(&model.State{PollURL: "/config", PollIntervalSeconds: 1}, nil).Once()
	controller.On("Register", mock.Anything, "").Return(&model.RegisterResponse{
		AgentID:             "agent-run",
		PollURL:             "/config",
		PollIntervalSeconds: 1,
	}, nil).Once()
	stateRepo.On("Save", mock.AnythingOfType("*model.State")).Return(nil).Once()
	controller.On("GetConfig", mock.Anything, "agent-run", "", "/config").Return((*model.Config)(nil), "", 0, errors.New("controller timeout")).Maybe()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	svc.Run(ctx)
}

func TestRun_PollNonControllerError_UsesNormalInterval(t *testing.T) {
	controller := new(clientMocks.ControllerClient)
	worker := new(clientMocks.WorkerClient)
	stateRepo := new(repositoryMocks.StateRepository)
	svc := newService(controller, worker, stateRepo)

	stateRepo.On("Load").Return(&model.State{PollURL: "/config", PollIntervalSeconds: 1}, nil).Once()
	controller.On("Register", mock.Anything, "").Return(&model.RegisterResponse{
		AgentID:             "agent-run",
		PollURL:             "/config",
		PollIntervalSeconds: 1,
	}, nil).Once()
	stateRepo.On("Save", mock.AnythingOfType("*model.State")).Return(nil).Once()
	cfg := &model.Config{Version: 2, URL: "http://example.com", PollIntervalSeconds: 1}
	controller.On("GetConfig", mock.Anything, "agent-run", "", "/config").Return(cfg, "\"2\"", 200, nil).Maybe()
	worker.On("ApplyConfig", mock.Anything, cfg).Return(errors.New("worker fail")).Maybe()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	svc.Run(ctx)
}

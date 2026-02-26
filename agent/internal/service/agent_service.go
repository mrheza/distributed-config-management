package service

import (
	"agent/internal/client"
	"agent/internal/model"
	"agent/internal/repository"
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

type AgentService interface {
	Run(ctx context.Context)
	GetState() *model.State
}

type agentService struct {
	controller       client.ControllerClient
	worker           client.WorkerClient
	stateRepo        repository.StateRepository
	defaultPollURL   string
	defaultPollSecs  int
	maxBackoffSecs   int
	backoffJitterPct int
	currentState     *model.State
}

type reqError struct {
	err    error
	target string
}

func (e *reqError) Error() string { return e.err.Error() }
func (e *reqError) Unwrap() error { return e.err }

func NewAgentService(
	controller client.ControllerClient,
	worker client.WorkerClient,
	stateRepo repository.StateRepository,
	defaultPollURL string,
	defaultPollSecs int,
	maxBackoffSecs int,
	backoffJitterPct int,
) AgentService {
	rand.Seed(time.Now().UnixNano())

	return &agentService{
		controller:       controller,
		worker:           worker,
		stateRepo:        stateRepo,
		defaultPollURL:   defaultPollURL,
		defaultPollSecs:  defaultPollSecs,
		maxBackoffSecs:   maxBackoffSecs,
		backoffJitterPct: backoffJitterPct,
		currentState: &model.State{
			PollURL:             defaultPollURL,
			PollIntervalSeconds: defaultPollSecs,
		},
	}
}

func (s *agentService) GetState() *model.State {
	clone := *s.currentState
	return &clone
}

func (s *agentService) Run(ctx context.Context) {
	log.Printf(
		"event=agent_run_started default_poll_secs=%d max_backoff_secs=%d backoff_jitter_pct=%d",
		s.defaultPollSecs,
		s.maxBackoffSecs,
		s.backoffJitterPct,
	)

	bootstrapRetryCount := 0
	lastBootstrapRetryTarget := ""
	for {
		err := s.bootstrap(ctx)
		if err == nil {
			bootstrapRetryCount = 0
			lastBootstrapRetryTarget = ""
			break
		}

		log.Printf("event=bootstrap_failed err=%q", err)

		var reqErr *reqError
		if errors.As(err, &reqErr) {
			target := reqErr.target
			if target == "" {
				target = "remote"
			}
			if lastBootstrapRetryTarget != "" && lastBootstrapRetryTarget != target {
				bootstrapRetryCount = 0
			}
			lastBootstrapRetryTarget = target

			bootstrapRetryCount++
			sleep := calculateBackoff(bootstrapRetryCount, s.maxBackoffSecs)
			sleep = applyJitter(sleep, s.backoffJitterPct)
			log.Printf(
				"event=bootstrap_retry_scheduled target=%s retry_count=%d sleep_secs=%.3f",
				target,
				bootstrapRetryCount,
				sleep.Seconds(),
			)
			select {
			case <-ctx.Done():
				return
			case <-time.After(sleep):
			}
			continue
		}

		bootstrapRetryCount = 0
		lastBootstrapRetryTarget = ""
		interval := s.defaultPollSecs
		if interval <= 0 {
			interval = 5
		}
		log.Printf("event=bootstrap_local_error_fixed_retry sleep_secs=%d err=%q", interval, err)
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(interval) * time.Second):
		}
	}

	if s.currentState != nil {
		log.Printf(
			"event=bootstrap_completed agent_id=%s poll_url=%s poll_interval_secs=%d etag=%q",
			s.currentState.AgentID,
			s.currentState.PollURL,
			s.currentState.PollIntervalSeconds,
			s.currentState.ETag,
		)
	}

	retryCount := 0
	lastPollRetryTarget := ""
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := s.pollOnce(ctx)
		if err != nil {
			var reqErr *reqError
			if errors.As(err, &reqErr) {
				target := reqErr.target
				if target == "" {
					target = "remote"
				}
				if lastPollRetryTarget != "" && lastPollRetryTarget != target {
					retryCount = 0
				}
				lastPollRetryTarget = target

				retryCount++
				sleep := calculateBackoff(retryCount, s.maxBackoffSecs)
				sleep = applyJitter(sleep, s.backoffJitterPct)
				log.Printf(
					"event=poll_retry_scheduled target=%s retry_count=%d sleep_secs=%.3f err=%q",
					target,
					retryCount,
					sleep.Seconds(),
					err,
				)
				select {
				case <-ctx.Done():
					return
				case <-time.After(sleep):
				}
				continue
			}
		}

		retryCount = 0
		lastPollRetryTarget = ""
		interval := s.currentState.PollIntervalSeconds
		if interval <= 0 {
			interval = s.defaultPollSecs
		}
		log.Printf("event=next_poll_scheduled sleep_secs=%d", interval)

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(interval) * time.Second):
		}
	}
}
func calculateBackoff(retryCount, maxBackoffSecs int) time.Duration {
	if retryCount < 1 {
		retryCount = 1
	}

	maxBackoff := time.Duration(maxBackoffSecs) * time.Second
	if maxBackoff <= 0 {
		maxBackoff = time.Second
	}

	backoff := time.Second
	for i := 1; i < retryCount; i++ {
		if backoff >= maxBackoff {
			return maxBackoff
		}
		backoff *= 2
	}

	if backoff > maxBackoff {
		return maxBackoff
	}
	return backoff
}

func applyJitter(base time.Duration, jitterPercent int) time.Duration {
	if base <= 0 || jitterPercent <= 0 {
		return base
	}

	if jitterPercent > 90 {
		jitterPercent = 90
	}

	delta := float64(base) * float64(jitterPercent) / 100.0
	min := float64(base) - delta
	max := float64(base) + delta
	jittered := min + rand.Float64()*(max-min)
	if jittered < 0 {
		return 0
	}

	return time.Duration(jittered)
}

func (s *agentService) bootstrap(ctx context.Context) error {
	log.Printf("event=bootstrap_started")

	state, err := s.stateRepo.Load()
	if err != nil {
		return err
	}

	// Backward compatibility: old state files may have ETag but no cached config URL.
	// In that case, force the next poll to fetch full config instead of 304.
	if state.ConfigURL == "" && state.ETag != "" {
		log.Printf("event=state_missing_config_url_reset_etag old_etag=%q", state.ETag)
		state.ETag = ""
		state.LastConfigVersion = 0
	}

	log.Printf(
		"event=state_loaded agent_id=%s config_url=%s poll_url=%s poll_interval_secs=%d etag=%q last_config_version=%d",
		state.AgentID,
		state.ConfigURL,
		state.PollURL,
		state.PollIntervalSeconds,
		state.ETag,
		state.LastConfigVersion,
	)
	if state.PollURL == "" {
		state.PollURL = s.defaultPollURL
	}
	if state.PollIntervalSeconds <= 0 {
		state.PollIntervalSeconds = s.defaultPollSecs
	}

	// Rehydrate worker from local state so worker still has config even if controller returns 304.
	if state.ConfigURL != "" {
		cached := &model.Config{
			Version:             state.LastConfigVersion,
			URL:                 state.ConfigURL,
			PollIntervalSeconds: state.PollIntervalSeconds,
		}
		if err := s.worker.ApplyConfig(ctx, cached); err != nil {
			return &reqError{err: err, target: "worker"}
		}
		log.Printf(
			"event=worker_rehydrated_from_state version=%d url=%s poll_interval_secs=%d",
			cached.Version,
			cached.URL,
			cached.PollIntervalSeconds,
		)
	}

	reg, err := s.controller.Register(ctx, state.AgentID)
	if err != nil {
		return &reqError{err: err, target: "controller"}
	}
	log.Printf(
		"event=register_success agent_id=%s poll_url=%s poll_interval_secs=%d",
		reg.AgentID,
		reg.PollURL,
		reg.PollIntervalSeconds,
	)

	state.AgentID = reg.AgentID
	if reg.PollURL != "" {
		state.PollURL = reg.PollURL
	}
	if reg.PollIntervalSeconds > 0 {
		state.PollIntervalSeconds = reg.PollIntervalSeconds
	}

	if err := s.stateRepo.Save(state); err != nil {
		return err
	}
	log.Printf(
		"event=state_saved agent_id=%s config_url=%s poll_url=%s poll_interval_secs=%d etag=%q last_config_version=%d",
		state.AgentID,
		state.ConfigURL,
		state.PollURL,
		state.PollIntervalSeconds,
		state.ETag,
		state.LastConfigVersion,
	)

	s.currentState = state
	return nil
}

func (s *agentService) pollOnce(ctx context.Context) error {
	log.Printf(
		"event=poll_started agent_id=%s poll_url=%s etag=%q",
		s.currentState.AgentID,
		s.currentState.PollURL,
		s.currentState.ETag,
	)

	cfg, newETag, status, err := s.controller.GetConfig(
		ctx,
		s.currentState.AgentID,
		s.currentState.ETag,
		s.currentState.PollURL,
	)
	if err != nil {
		return &reqError{err: err, target: "controller"}
	}
	log.Printf(
		"event=poll_response status=%d etag=%q",
		status,
		newETag,
	)

	if status == 304 {
		log.Printf(
			"event=config_not_modified status=304 agent_id=%s etag=%q",
			s.currentState.AgentID,
			newETag,
		)
		return nil
	}

	if cfg == nil {
		log.Printf("event=config_empty")
		return nil
	}
	log.Printf(
		"event=config_received version=%d poll_interval_secs=%d url=%s",
		cfg.Version,
		cfg.PollIntervalSeconds,
		cfg.URL,
	)

	if err := s.worker.ApplyConfig(ctx, cfg); err != nil {
		return &reqError{err: err, target: "worker"}
	}
	log.Printf("event=worker_apply_success version=%d", cfg.Version)

	s.currentState.ETag = newETag
	s.currentState.ConfigURL = cfg.URL
	s.currentState.LastConfigVersion = cfg.Version
	if cfg.PollIntervalSeconds > 0 {
		s.currentState.PollIntervalSeconds = cfg.PollIntervalSeconds
	}

	if err := s.stateRepo.Save(s.currentState); err != nil {
		return err
	}
	log.Printf(
		"event=state_saved agent_id=%s config_url=%s etag=%q last_config_version=%d poll_interval_secs=%d",
		s.currentState.AgentID,
		s.currentState.ConfigURL,
		s.currentState.ETag,
		s.currentState.LastConfigVersion,
		s.currentState.PollIntervalSeconds,
	)

	return nil
}


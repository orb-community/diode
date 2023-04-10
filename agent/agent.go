/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"context"
	"errors"
	"time"

	"github.com/orb-community/diode/agent/backend"
	"github.com/orb-community/diode/agent/backend/factory"
	"github.com/orb-community/diode/agent/config"
	"github.com/orb-community/diode/agent/pusher"
	"go.uber.org/zap"
)

const (
	Kind = "discovery"
)

type Agent interface {
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
	Stop(ctx context.Context)
	RestartAll(ctx context.Context, reason string) error
	RestartBackend(ctx context.Context, backend string, reason string) error
}

type backendInfo struct {
	be backend.Backend
	st *backend.State
	cf config.Policy
}

type diodeAgent struct {
	logger         *zap.Logger
	config         config.Config
	backends       map[string]backendInfo
	cancelFunction context.CancelFunc
	pusher         pusher.Pusher
}

var _ Agent = (*diodeAgent)(nil)

func New(logger *zap.Logger, c config.Config) (Agent, error) {
	var s pusher.Pusher
	var err error
	if s, err = pusher.New(logger, c); err != nil {
		return nil, err
	}
	return &diodeAgent{logger: logger, config: c, pusher: s}, nil
}

func (a *diodeAgent) startPolicies(agentCtx context.Context) error {
	a.logger.Info("registered backends", zap.Strings("values", factory.GetList()))
	if len(a.config.DiodeAgent.Policies) == 0 {
		return errors.New("no policies specified")
	}
	a.backends = make(map[string]backendInfo, len(a.config.DiodeAgent.Policies))
	for name, policy := range a.config.DiodeAgent.Policies {
		be, err := factory.GetBackend(policy.Backend)
		if err != nil {
			return err
		}
		_, ok := a.backends[name]
		if ok {
			return errors.New("policy '" + name + "' already exists")
		}
		if policy.Kind != Kind {
			return errors.New("invalid policy kind")
		}
		if err = be.Configure(a.logger, name, a.pusher.GetChannel(), policy.Data, policy.Config); err != nil {
			return err
		}
		backendCtx := context.WithValue(agentCtx, "routine", name)

		if err := be.Start(context.WithCancel(backendCtx)); err != nil {
			return err
		}
		a.backends[name] = backendInfo{
			be: be,
			st: &backend.State{
				Status:        backend.Unknown,
				LastRestartTS: time.Now(),
			},
			cf: policy,
		}
	}
	return nil
}

func (a *diodeAgent) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	startTime := time.Now()
	defer func(t time.Time) {
		a.logger.Debug("Startup of agent execution duration", zap.String("Start() execution duration", time.Since(t).String()))
	}(startTime)

	agentCtx := context.WithValue(ctx, "routine", "agentRoutine")
	a.cancelFunction = cancelFunc

	pusherContext := context.WithValue(agentCtx, "routine", "pusherRoutine")
	if err := a.pusher.Start(context.WithCancel(pusherContext)); err != nil {
		return err
	}
	a.logger.Info("agent started", zap.Any("routine", agentCtx.Value("routine")))
	if err := a.startPolicies(ctx); err != nil {
		return err
	}

	return nil
}

func (a *diodeAgent) Stop(ctx context.Context) {
	a.logger.Info("routine call for stop agent", zap.Any("routine", ctx.Value("routine")))
	for name, b := range a.backends {
		if state, _, _ := b.be.GetRunningStatus(); state == backend.Running {
			a.logger.Debug("stopping backend", zap.String("backend", name))
			if err := b.be.Stop(ctx); err != nil {
				a.logger.Error("error while stopping the backend", zap.String("backend", name))
			}
		}
	}
	a.pusher.Stop(ctx)
	defer a.cancelFunction()
}

func (a *diodeAgent) RestartBackend(ctx context.Context, name string, reason string) error {
	return nil
}

func (a *diodeAgent) RestartAll(ctx context.Context, reason string) error {
	return nil
}

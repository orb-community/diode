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
	"github.com/orb-community/diode/agent/scrapper"
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

type diodeAgent struct {
	logger         *zap.Logger
	config         config.Config
	backends       map[string]backend.Backend
	backendState   map[string]*backend.State
	cancelFunction context.CancelFunc
	scrapper       scrapper.Scrapper
}

var _ Agent = (*diodeAgent)(nil)

func New(logger *zap.Logger, c config.Config) (Agent, error) {
	var s scrapper.Scrapper
	var err error
	if s, err = scrapper.New(logger, c); err != nil {
		return nil, err
	}
	return &diodeAgent{logger: logger, config: c, scrapper: s}, nil
}

func (a *diodeAgent) startPolicies(agentCtx context.Context) error {
	a.logger.Info("registered backends", zap.Strings("values", factory.GetList()))
	if len(a.config.DiodeAgent.Policies) == 0 {
		return errors.New("no policies specified")
	}
	a.backends = make(map[string]backend.Backend, len(a.config.DiodeAgent.Policies))
	a.backendState = make(map[string]*backend.State)
	for name, policy := range a.config.DiodeAgent.Policies {
		be, err := factory.GetBackend(policy.Backend)
		if err != nil {
			return err
		}
		if a.backends[name] != nil {
			return errors.New("policy '" + name + "' already exists")
		}
		if policy.Kind != Kind {
			return errors.New("invalid policy kind")
		}
		if err = be.Configure(a.logger, name, a.scrapper.GetChannel(), policy.Data); err != nil {
			return err
		}
		backendCtx := context.WithValue(agentCtx, "routine", name)

		if err := be.Start(context.WithCancel(backendCtx)); err != nil {
			return err
		}
		a.backends[name] = be
		a.backendState[name] = &backend.State{
			Status:        backend.Unknown,
			LastRestartTS: time.Now(),
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

	a.logger.Info("agent started", zap.Any("routine", agentCtx.Value("routine")))
	if err := a.startPolicies(ctx); err != nil {
		return err
	}

	return nil
}

func (a *diodeAgent) Stop(ctx context.Context) {

}

func (a *diodeAgent) RestartBackend(ctx context.Context, name string, reason string) error {
	return nil
}

func (a *diodeAgent) RestartAll(ctx context.Context, reason string) error {
	return nil
}

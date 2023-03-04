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
	"go.uber.org/zap"
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
}

var _ Agent = (*diodeAgent)(nil)

func New(logger *zap.Logger, c config.Config) (Agent, error) {
	return &diodeAgent{logger: logger, config: c}, nil
}

func (a *diodeAgent) startBackends(agentCtx context.Context) error {
	a.logger.Info("registered backends", zap.Strings("values", factory.GetList()))
	a.logger.Info("requested backends", zap.Any("values", a.config.DiodeAgent.Backends))
	if len(a.config.DiodeAgent.Backends) == 0 {
		return errors.New("no backends specified")
	}
	a.backends = make(map[string]backend.Backend, len(a.config.DiodeAgent.Backends))
	a.backendState = make(map[string]*backend.State)
	for name, configurationEntry := range a.config.DiodeAgent.Backends {
		be, err := factory.GetBackend(name)
		if err != nil {
			return err
		}
		if err = be.Configure(a.logger, configurationEntry); err != nil {
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
	if err := a.startBackends(ctx); err != nil {
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

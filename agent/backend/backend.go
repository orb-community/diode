/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package backend

import (
	"context"
	"time"

	"go.uber.org/zap"
)

const (
	Unknown RunningStatus = iota
	Running
	BackendError
	AgentError
	Offline
)

type RunningStatus int

type State struct {
	Status            RunningStatus
	RestartCount      int64
	LastError         string
	LastRestartTS     time.Time
	LastRestartReason string
}

type Backend interface {
	Configure(*zap.Logger, string, chan []byte, map[string]interface{}, map[string]interface{}) error
	Version() (string, error)
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
	Stop(ctx context.Context) error
	FullReset(ctx context.Context) error

	GetStartTime() time.Time
	GetCapabilities() (map[string]interface{}, error)
	GetRunningStatus() (RunningStatus, string, error)
}

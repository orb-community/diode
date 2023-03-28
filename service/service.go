/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package service

import (
	"context"
	"encoding/json"

	"github.com/orb-community/diode/service/storage"

	"github.com/orb-community/diode/service/config"
	"github.com/orb-community/diode/service/nb_pusher"
	"github.com/orb-community/diode/service/otlp"
	"go.uber.org/zap"
)

type Service interface {
	Start() error
	Stop() error
}

type DiodeService struct {
	logger             *zap.Logger
	config             *config.Config
	channel            chan []byte
	otlpRecv           otlp.Otlp
	pusher             nb_pusher.Pusher
	cancelAsyncContext context.CancelFunc
	asyncContext       context.Context
	storageService     storage.Service
}

var _ Service = (*DiodeService)(nil)

func New(ctx context.Context, cancelFunc context.CancelFunc, logger *zap.Logger, config *config.Config) (Service, error) {
	pusher := nb_pusher.New(ctx, logger, config)
	err := pusher.Start()
	if err != nil {
		cancelFunc()
		return nil, err
	}
	channel := make(chan []byte)
	otlpRecv := otlp.New(ctx, logger, config, channel)
	err = otlpRecv.Start()
	if err != nil {
		cancelFunc()
		return nil, err
	}
	service, err := storage.NewSqliteStorage(logger)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	return &DiodeService{
		logger:             logger,
		config:             config,
		channel:            channel,
		otlpRecv:           otlpRecv,
		pusher:             pusher,
		cancelAsyncContext: cancelFunc,
		asyncContext:       ctx,
		storageService:     service,
	}, nil
}

func (ds *DiodeService) Start() error {

	go func() {
		var jsonData map[string]interface{}
		for {
			select {
			case data := <-ds.channel:
				err := json.Unmarshal(data, &jsonData)
				if err != nil {
					ds.logger.Error("fail to unmarshal json", zap.Error(err))
					break
				}
				for policy, logEntryData := range jsonData {
					if _, err := ds.storageService.Save(policy, logEntryData.(map[string]interface{})); err != nil {
						ds.logger.Error("error during storing", zap.String("policy", policy), zap.Error(err))
					}
				}
			case <-ds.asyncContext.Done():
				ds.logger.Info("service context cancelled")
				return
			}
		}
	}()

	ds.logger.Info("diode service started")
	return nil
}

func (ds *DiodeService) Stop() error {
	err := ds.otlpRecv.Stop()
	if err != nil {
		return err
	}
	return nil
}

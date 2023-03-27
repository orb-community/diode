/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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
	db                 *sql.DB
	logger             *zap.Logger
	config             config.Config
	channel            chan []byte
	otlpRecv           otlp.Otlp
	pusher             nb_pusher.Pusher
	cancelAsyncContext context.CancelFunc
	asyncContext       context.Context
}

var _ Service = (*DiodeService)(nil)

func New(logger *zap.Logger, config config.Config, db *sql.DB) Service {
	return &DiodeService{
		logger:  logger,
		config:  config,
		channel: make(chan []byte),
	}
}

func (ds *DiodeService) Start() error {
	ds.asyncContext, ds.cancelAsyncContext = context.WithCancel(context.WithValue(context.Background(), "routine", "async"))

	ds.pusher = nb_pusher.New(ds.asyncContext, ds.logger, &ds.config)
	err := ds.pusher.Start()
	if err != nil {
		return err
	}

	service := storage.NewSqliteStorage(ds.db)

	ds.otlpRecv = otlp.New(ds.asyncContext, ds.logger, &ds.config, service, ds.channel)
	err = ds.otlpRecv.Start()
	if err != nil {
		return err
	}

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
				for policy, v := range jsonData {
					ds.logger.Info("policy name " + policy)
					ds.logger.Info("data " + fmt.Sprintf("%v", v))
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
	ds.otlpRecv.Stop()
	return nil
}

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package scrapper

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/orb-community/diode/agent/config"
	"go.uber.org/zap"
)

const (
	File = "file"
	Otlp = "otlp"
)

type Scrapper interface {
	GetChannel() chan []byte
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
	Stop(ctx context.Context)
}

type scrapperImpl struct {
	logger     *zap.Logger
	outputPath string
	outputType string
	channel    chan []byte
	cancelFunc context.CancelFunc
	ctx        context.Context
}

var _ Scrapper = (*scrapperImpl)(nil)

func New(logger *zap.Logger, c config.Config) (Scrapper, error) {
	if c.DiodeAgent.DiodeConfig.OutputType == File {
		if _, err := os.Stat(c.DiodeAgent.DiodeConfig.OutputPath); os.IsNotExist(err) {
			return nil, errors.New("output path '" + c.DiodeAgent.DiodeConfig.OutputPath + "' does not exist")
		}
	}
	return &scrapperImpl{logger: logger, outputType: c.DiodeAgent.DiodeConfig.OutputType,
		outputPath: c.DiodeAgent.DiodeConfig.OutputPath, channel: make(chan []byte)}, nil
}

func (s *scrapperImpl) GetChannel() chan []byte {
	return s.channel
}

func (s *scrapperImpl) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	s.cancelFunc = cancelFunc
	s.ctx = ctx
	switch o := s.outputType; o {
	case File:
		return s.scrapeToFile()
	case Otlp:
		return errors.New("OTLP not implemented yet")
	default:
		return errors.New(s.outputType + " is a invalid output type")
	}
}

func (s *scrapperImpl) Stop(ctx context.Context) {
	s.logger.Info("routine call to stop scrapper", zap.Any("routine", ctx.Value("routine")))
	defer s.cancelFunc()
}

func (s *scrapperImpl) scrapeToFile() error {
	go func() {
		var jsonData map[string]interface{}
		for {
			select {
			case data := <-s.channel:
				json.Unmarshal(data, &jsonData)
				for policy := range jsonData {
					path := s.outputPath + "/" + policy + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
					if err := os.WriteFile(path, data, 0644); err != nil {
						s.logger.Error("fail to generate output file for policy "+policy, zap.Error(err))
					}
				}
			case <-s.ctx.Done():
				close(s.channel)
				s.logger.Info("scrapper context cancelled")
				return
			}
		}
	}()
	return nil
}

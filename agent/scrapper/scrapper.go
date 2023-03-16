/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package scrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/orb-community/diode/agent/config"
	"go.uber.org/zap"
)

const (
	File = "file"
	Otlp = "otlp"
	Http = "http"
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
	outputAuth string
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
	} else if c.DiodeAgent.DiodeConfig.OutputType == Http {
		if _, err := url.ParseRequestURI(c.DiodeAgent.DiodeConfig.OutputPath); err != nil {
			return nil, err
		}
	}
	return &scrapperImpl{logger: logger, outputType: c.DiodeAgent.DiodeConfig.OutputType, outputPath: c.DiodeAgent.DiodeConfig.OutputPath,
		outputAuth: c.DiodeAgent.DiodeConfig.OutputAuth, channel: make(chan []byte)}, nil
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
	case Http:
		return s.scrapeToHttp()
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

func (s *scrapperImpl) scrapeToHttp() error {
	go func() {
		for {
			select {
			case data := <-s.channel:
				client := &http.Client{}
				req, err := http.NewRequest("POST", s.outputPath, bytes.NewBuffer(data))
				if err != nil {
					s.logger.Error("scrapper: fail to create http request", zap.Error(err))
					continue
				}
				req.Header.Add("Content-Type", "application/json")
				if s.outputAuth != "" {
					req.Header.Add("Authorization", s.outputAuth)
				}

				res, err := client.Do(req)
				if err != nil {
					s.logger.Error("scrapper: fail to create http request", zap.Error(err))
					continue
				}
				defer res.Body.Close()
			case <-s.ctx.Done():
				close(s.channel)
				s.logger.Info("scrapper context cancelled")
				return
			}
		}
	}()
	return nil
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

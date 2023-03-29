/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pusher

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/orb-community/diode/agent/config"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	File = "file"
	Otlp = "otlp"
	Http = "http"
)

type Pusher interface {
	GetChannel() chan []byte
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
	Stop(ctx context.Context)
}

type pusherImpl struct {
	logger     *zap.Logger
	outputPath string
	outputType string
	outputAuth string
	channel    chan []byte
	cancelFunc context.CancelFunc
	ctx        context.Context
}

var _ Pusher = (*pusherImpl)(nil)

func New(logger *zap.Logger, c config.Config) (Pusher, error) {
	if c.DiodeAgent.DiodeConfig.OutputType == File {
		if _, err := os.Stat(c.DiodeAgent.DiodeConfig.OutputPath); os.IsNotExist(err) {
			return nil, errors.New("output path '" + c.DiodeAgent.DiodeConfig.OutputPath + "' does not exist")
		}
	}
	return &pusherImpl{logger: logger, outputType: c.DiodeAgent.DiodeConfig.OutputType, outputPath: c.DiodeAgent.DiodeConfig.OutputPath,
		outputAuth: c.DiodeAgent.DiodeConfig.OutputAuth, channel: make(chan []byte, 16)}, nil
}

func (s *pusherImpl) GetChannel() chan []byte {
	return s.channel
}

func (s *pusherImpl) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	s.cancelFunc = cancelFunc
	s.ctx = ctx
	switch o := s.outputType; o {
	case File:
		return s.scrapeToFile()
	case Http:
		return s.scrapeToHttp()
	case Otlp:
		return s.scrapeToOtlp()
	default:
		return errors.New(s.outputType + " is a invalid output type")
	}
}

func (s *pusherImpl) Stop(ctx context.Context) {
	s.logger.Info("routine call to stop pusher", zap.Any("routine", ctx.Value("routine")))
	defer s.cancelFunc()
}

type tempNetboxStruct struct {
	ObjType string      `json:"object_type"`
	Engine  string      `json:"engine"`
	Data    interface{} `json:"data"`
}

func (s *pusherImpl) temporaryMatchNetbox(data []byte) []byte {
	var jsonData map[string]map[string]interface{}
	var returnData tempNetboxStruct
	returnData.ObjType = "dcim.device"
	json.Unmarshal(data, &jsonData)
	for _, policy := range jsonData {
		for k, v := range policy {
			if k == "backend" {
				returnData.Engine = v.(string)
			} else if k == "device" {
				returnData.Data = v.([]interface{})[0]
			}
		}
	}
	b, _ := json.Marshal(returnData)
	return b
}

func (s *pusherImpl) scrapeToHttp() error {
	go func() {
		for {
			select {
			case data := <-s.channel:
				client := &http.Client{}
				req, err := http.NewRequest("POST", s.outputPath, bytes.NewBuffer(s.temporaryMatchNetbox(data)))
				if err != nil {
					s.logger.Error("pusher - fail to create http request", zap.Error(err))
					continue
				}
				req.Header.Add("Content-Type", "application/json")
				if s.outputAuth != "" {
					req.Header.Add("Authorization", s.outputAuth)
				}

				res, err := client.Do(req)
				if err != nil {
					s.logger.Error("pusher - fail to create http request", zap.Error(err))
					continue
				}
				defer res.Body.Close()
				s.logger.Info("pusher - http response status: " + res.Status)
			case <-s.ctx.Done():
				close(s.channel)
				s.logger.Info("pusher context cancelled")
				return
			}
		}
	}()
	return nil
}

func (s *pusherImpl) scrapeToFile() error {
	go func() {
		var jsonData map[string]interface{}
		for {
			select {
			case data := <-s.channel:
				err := json.Unmarshal(data, &jsonData)
				if err != nil {
					s.logger.Error("pusher - fail to unmarshal json", zap.Error(err))
					break
				}
				for policy := range jsonData {
					path := s.outputPath + "/" + policy + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
					if err := os.WriteFile(path, data, 0644); err != nil {
						s.logger.Error("pusher - fail to generate output file for policy "+policy, zap.Error(err))
						break
					}
				}
			case <-s.ctx.Done():
				close(s.channel)
				s.logger.Info("pusher context cancelled")
				return
			}
		}
	}()
	return nil
}

func (s *pusherImpl) scrapeToOtlp() error {
	factory := otlpexporter.NewFactory()
	factory.CreateDefaultConfig()
	set := exporter.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         s.logger,
			TracerProvider: trace.NewNoopTracerProvider(),
		},
	}
	cfg := factory.CreateDefaultConfig().(*otlpexporter.Config)
	cfg.GRPCClientSettings.Endpoint = s.outputPath
	cfg.GRPCClientSettings.TLSSetting = configtls.TLSClientSetting{
		Insecure: true,
	}
	lexporter, err := factory.CreateLogsExporter(s.ctx, set, cfg)
	if err != nil {
		s.logger.Error("pusher - fail to create log exporter", zap.Error(err))
		return err
	}
	err = lexporter.Start(s.ctx, nil)
	if err != nil {
		s.logger.Error("pusher - fail to start log exporter", zap.Error(err))
		return err
	}
	logs := plog.NewLogs()
	res := logs.ResourceLogs().AppendEmpty()
	scope := res.ScopeLogs().AppendEmpty()
	record := scope.LogRecords().AppendEmpty()
	record.SetSeverityNumber(plog.SeverityNumberTrace)

	go func() {
		for {
			select {
			case data := <-s.channel:
				err = record.Body().FromRaw(data)
				if err != nil {
					s.logger.Error("pusher - fail to add log body", zap.Error(err))
					break
				}
				lexporter.ConsumeLogs(s.ctx, logs)
			case <-s.ctx.Done():
				close(s.channel)
				s.logger.Info("pusher context cancelled")
				return
			}
		}
	}()
	return nil
}

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package otlp

import (
	"context"
	"github.com/orb-community/diode/service/config"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Otlp interface {
	Start() error
	Stop() error
}

type DiodeOtlp struct {
	ctx      context.Context
	logger   *zap.Logger
	config   *config.Config
	consumer consumer.Logs
	receiver receiver.Logs
}

type DiodeLogConsumer struct {
	channel      chan []byte
	capabilities consumer.Capabilities
}

var _ Otlp = (*DiodeOtlp)(nil)

func New(ctx context.Context, logger *zap.Logger, config *config.Config, channel chan []byte) Otlp {
	return &DiodeOtlp{ctx: ctx, logger: logger, config: config, consumer: newLogConsumer(channel)}
}

func (d *DiodeOtlp) Start() error {
	factory := otlpreceiver.NewFactory()
	cfg := factory.CreateDefaultConfig().(*otlpreceiver.Config)
	cfg.HTTP = nil
	cfg.GRPC.NetAddr.Endpoint = d.config.OtlpReceiver.Endpoint
	cfg.GRPC.NetAddr.Transport = d.config.OtlpReceiver.Protocol
	set := receiver.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         d.logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	var err error
	d.receiver, err = factory.CreateLogsReceiver(d.ctx, set, cfg, d.consumer)
	if err != nil {
		return err
	}
	err = d.receiver.Start(d.ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *DiodeOtlp) Stop() error {
	return nil
}

func newLogConsumer(channel chan []byte) consumer.Logs {
	var cap consumer.Capabilities
	cap.MutatesData = true
	return &DiodeLogConsumer{channel: channel, capabilities: cap}
}

func (dlc *DiodeLogConsumer) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	rss := ld.ResourceLogs()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		ill := rs.ScopeLogs()
		for i := 0; i < ill.Len(); i++ {
			logs := ill.At(i)
			rec := logs.LogRecords()
			for i := 0; i < rec.Len(); i++ {
				val := rec.At(i)
				dlc.channel <- val.Body().Bytes().AsRaw()
			}
		}
	}
	return nil
}

func (dlc *DiodeLogConsumer) Capabilities() consumer.Capabilities {
	return dlc.capabilities
}

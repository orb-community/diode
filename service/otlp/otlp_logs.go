/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package otlp

import (
	"context"

	"github.com/orb-community/diode/service/config"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type DiodeOtlpRecv struct {
	ctx      context.Context
	logger   *zap.Logger
	config   *config.Config
	consumer consumer.Logs
	receiver receiver.Logs
}

var _ Otlp = (*DiodeOtlpRecv)(nil)

func (d *DiodeOtlpRecv) Start() error {
	oFactory := otlpreceiver.NewFactory()
	cfg := oFactory.CreateDefaultConfig().(*otlpreceiver.Config)
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
	d.receiver, err = oFactory.CreateLogsReceiver(d.ctx, set, cfg, d.consumer)
	if err != nil {
		return err
	}
	err = d.receiver.Start(d.ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *DiodeOtlpRecv) Stop() error {
	return nil
}

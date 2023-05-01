/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package otlp

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver"
	"github.com/orb-community/diode/service/config"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type DiodeKafkaRecv struct {
	ctx      context.Context
	logger   *zap.Logger
	config   *config.Config
	consumer consumer.Logs
	receiver receiver.Logs
}

var _ Otlp = (*DiodeKafkaRecv)(nil)

func (d *DiodeKafkaRecv) Start() error {
	kFactory := kafkareceiver.NewFactory()
	cfg := kFactory.CreateDefaultConfig().(*kafkareceiver.Config)
	cfg.Brokers = d.config.KafkaReceiver.Brokers
	cfg.Topic = d.config.KafkaReceiver.Topic
	cfg.ProtocolVersion = d.config.KafkaReceiver.ProtocolVersion
	set := receiver.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         d.logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	var err error
	d.receiver, err = kFactory.CreateLogsReceiver(d.ctx, set, cfg, d.consumer)
	if err != nil {
		return err
	}
	err = d.receiver.Start(d.ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *DiodeKafkaRecv) Stop() error {
	return nil
}

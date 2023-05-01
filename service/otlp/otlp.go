/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package otlp

import (
	"context"

	"github.com/orb-community/diode/service/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type Otlp interface {
	Start() error
	Stop() error
}

type DiodeLogConsumer struct {
	channel      chan []byte
	capabilities consumer.Capabilities
}

func New(ctx context.Context, logger *zap.Logger, config *config.Config, channel chan []byte) Otlp {
	switch tOtlp := config.Base.OtlpReceiverType; tOtlp {
	case "kafka":
		return &DiodeKafkaRecv{ctx: ctx, logger: logger, config: config, consumer: newLogConsumer(channel)}
	case "otlp":
		return &DiodeOtlpRecv{ctx: ctx, logger: logger, config: config, consumer: newLogConsumer(channel)}
	default:
		break
	}
	logger.Warn("Not supported OTLP receiver type. Creating Default OTLP Receiver",
		zap.String("otlp_receiver_type", config.Base.OtlpReceiverType))
	return &DiodeOtlpRecv{ctx: ctx, logger: logger, config: config, consumer: newLogConsumer(channel)}
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

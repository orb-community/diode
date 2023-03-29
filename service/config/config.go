/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type BaseSvcConfig struct {
	LogLevel       string `mapstructure:"log_level"`
	HttpPort       string `mapstructure:"http_port"`
	HttpServerCert string `mapstructure:"server_cert"`
	HttpServerKey  string `mapstructure:"server_key"`
}

type NetboxPusherConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Token    string `mapstructure:"token"`
	Protocol string `mapstructure:"protocol"`
}

type OtlpReceiverConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Protocol string `mapstructure:"protocol"`
}

type Config struct {
	Base         BaseSvcConfig
	NetboxPusher NetboxPusherConfig
	OtlpReceiver OtlpReceiverConfig
}

const (
	otlpEndpoint = "0.0.0.0:4317"
	otlpProtocol = "tcp"
)

func LoadConfig(prefix string) Config {
	var config Config
	config.Base = loadBaseServiceConfig(prefix)
	config.NetboxPusher = loadNetboxPusherConfig(prefix)
	config.OtlpReceiver = loadOtlpReceiverConfig(prefix)
	return config
}

func loadBaseServiceConfig(prefix string) BaseSvcConfig {
	cfg := viper.New()
	cfg.SetEnvPrefix(prefix)

	cfg.SetDefault("log_level", "error")
	cfg.SetDefault("http_port", "")
	cfg.SetDefault("server_cert", "")
	cfg.SetDefault("server_key", "")

	cfg.AllowEmptyEnv(true)
	cfg.AutomaticEnv()
	var svcC BaseSvcConfig
	cfg.Unmarshal(&svcC)
	return svcC
}

func loadNetboxPusherConfig(prefix string) NetboxPusherConfig {
	cfg := viper.New()
	cfg.SetEnvPrefix(fmt.Sprintf("%s_netbox", prefix))

	cfg.SetDefault("endpoint", "")
	cfg.SetDefault("token", "")
	cfg.SetDefault("protocol", "https")

	cfg.AllowEmptyEnv(true)
	cfg.AutomaticEnv()
	var nbC NetboxPusherConfig
	cfg.Unmarshal(&nbC)
	return nbC
}

func loadOtlpReceiverConfig(prefix string) OtlpReceiverConfig {
	cfg := viper.New()
	cfg.SetEnvPrefix(fmt.Sprintf("%s_otlp", prefix))

	cfg.SetDefault("endpoint", otlpEndpoint)
	cfg.SetDefault("protocol", otlpProtocol)

	cfg.AllowEmptyEnv(true)
	cfg.AutomaticEnv()
	var otlpC OtlpReceiverConfig
	cfg.Unmarshal(&otlpC)
	return otlpC
}

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

import "github.com/spf13/viper"

type EncryptionKey struct {
	Key string `mapstructure:"key"`
}

type BaseSvcConfig struct {
	LogLevel       string `mapstructure:"log_level"`
	HttpPort       string `mapstructure:"http_port"`
	HttpServerCert string `mapstructure:"server_cert"`
	HttpServerKey  string `mapstructure:"server_key"`
}

func LoadBaseServiceConfig(prefix string, httpPort string) BaseSvcConfig {
	cfg := viper.New()
	cfg.SetEnvPrefix(prefix)

	cfg.SetDefault("log_level", "error")
	cfg.SetDefault("http_port", httpPort)
	cfg.SetDefault("server_cert", "")
	cfg.SetDefault("server_key", "")

	cfg.AllowEmptyEnv(true)
	cfg.AutomaticEnv()
	var svcC BaseSvcConfig
	cfg.Unmarshal(&svcC)
	return svcC
}

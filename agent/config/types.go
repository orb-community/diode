/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

import "time"

type Status struct {
	StartTime time.Time     `json:"start_time"`
	UpTime    time.Duration `json:"up_time"`
	Version   string        `json:"version"`
}

type Policy struct {
	Kind    string                 `mapstructure:"kind"`
	Backend string                 `mapstructure:"backend"`
	Config  map[string]interface{} `mapstructure:"config"`
	Data    map[string]interface{} `mapstructure:"data"`
}

type DiodeConfig struct {
	Debug      bool   `mapstructure:"debug"`
	OutputType string `mapstructure:"output_type"`
	OutputPath string `mapstructure:"output_path"`
	OutputAuth string `mapstructure:"output_auth"`
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
}

type DiodeAgent struct {
	Tags        map[string]string `mapstructure:"tags"`
	DiodeConfig DiodeConfig       `mapstructure:"config"`
	Policies    map[string]Policy `mapstructure:"policies"`
}

type Config struct {
	Version    string     `mapstructure:"version"`
	DiodeAgent DiodeAgent `mapstructure:"diode"`
}

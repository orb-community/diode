/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

type Debug struct {
	Enable bool `mapstructure:"enable"`
}

type DiodeAgent struct {
	Backends map[string]map[string]string `mapstructure:"backends"`
	Tags     map[string]string            `mapstructure:"tags"`
	Debug    Debug                        `mapstructure:"debug"`
}

type Config struct {
	Version    float64    `mapstructure:"version"`
	DiodeAgent DiodeAgent `mapstructure:"diode"`
}

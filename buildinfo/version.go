// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package buildinfo

// set via ldflags -X option at build time
var version = "unknown"

// minimum version of an agent that we allow to connect
const minAgentVersion string = "0.1.0-develop"

func GetVersion() string {
	return version
}

func GetMinAgentVersion() string {
	return minAgentVersion
}

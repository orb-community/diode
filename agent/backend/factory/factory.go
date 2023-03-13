/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package factory

import (
	"errors"

	"github.com/orb-community/diode/agent/backend"
	"github.com/orb-community/diode/agent/backend/suzieq"
)

func GetBackend(backendType string) (backend.Backend, error) {
	if backendType == "suzieq" {
		return suzieq.New(), nil
	}
	return nil, errors.New("backend type not found")
}

func GetList() []string {
	return []string{"suzieq"}
}

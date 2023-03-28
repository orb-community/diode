/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nb_pusher

var DeviceStatusMap = map[string]string{
	"alive": "active",
	"dead":  "offline",
}

var InterfaceStateMap = map[string]bool{
	"up":   true,
	"down": false,
}

type NetboxObject struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type NetboxSite struct {
	NetboxObject
	Status string `json:"status"`
}

type NetboxDeviceType struct {
	Mfr *NetboxObject `json:"manufacturer"`
	NetboxObject
}

type NetboxDevice struct {
	Site *NetboxSite       `json:"site"`
	Role *NetboxObject     `json:"device_role"`
	Type *NetboxDeviceType `json:"device_type"`
	NetboxObject
	Status string `json:"status"`
}

type NetboxInterface struct {
	DeviceID   int64  `json:"device_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Speed      int64  `json:"speed"`
	Mtu        int64  `json:"mtu"`
	MacAddress string `json:"mac_address"`
	State      string `json:"state"`
}
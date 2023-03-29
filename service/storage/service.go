package storage

type Service interface {
	Save(policy string, jsonData map[string]interface{}) (interface{}, error)
	UpdateInterface(id string, netboxId int64) (DbInterface, error)
	UpdateDevice(id string, netboxId int64) (DbDevice, error)
	UpdateVlan(id string, netboxId int64) (DbVlan, error)
	GetInterfaceByPolicyAndNamespaceAndHostname(policy, namespace, hostname string) ([]DbInterface, error)
	GetDevicesByPolicyAndNamespace(policy, namespace string) ([]DbDevice, error)
	GetDeviceByPolicyAndNamespaceAndHostname(policy, namespace, hostname string) (DbDevice, error)
	GetVlansByPolicyAndNamespaceAndHostname(policy, namespace, hostname string) ([]DbVlan, error)
}

type DbInterface struct {
	Id          string `json:"id,omitempty"`
	Policy      string `json:"policy,omitempty"`
	Namespace   string `json:"namespace"`
	Hostname    string `json:"hostname"`
	Name        string `json:"ifname"`
	AdminState  string `json:"adminState"`
	Mtu         int64  `json:"mtu"`
	Speed       int64  `json:"speed"`
	MacAddress  string `json:"macaddr"`
	IfType      string `json:"type"`
	NetboxRefId int64  `json:"netbox_id,omitempty"`
	Blob        string `json:"blob,omitempty"`
}

type DbDevice struct {
	Id           string `json:"id,omitempty"`
	Policy       string `json:"policy,omitempty"`
	SerialNumber string `json:"serialNumber"`
	Namespace    string `json:"namespace"`
	Hostname     string `json:"hostname"`
	Address      string `json:"address"`
	Model        string `json:"model"`
	State        string `json:"state"`
	Vendor       string `json:"vendor"`
	NetboxRefId  int64  `json:"netbox_id,omitempty"`
	Blob         string `json:"blob,omitempty"`
}

type DbVlan struct {
	Id          string `json:"id,omitempty"`
	Policy      string `json:"policy,omitempty"`
	Namespace   string `json:"namespace"`
	Hostname    string `json:"hostname"`
	Name        string `json:"vlanName"`
	State       string `json:"state"`
	NetboxRefId int64  `json:"netbox_id,omitempty"`
	Blob        string `json:"blob,omitempty"`
}

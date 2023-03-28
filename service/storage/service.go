package storage

type Service interface {
	Save(policy string, jsonData map[string]interface{}) (interface{}, error)
}

type DbInterface struct {
	Id          string `json:"id,omitempty"`
	Policy      string `json:"policy,omitempty"`
	Namespace   string `json:"namespace"`
	Hostname    string `json:"hostname"`
	Name        string `json:"name"`
	AdminState  string `json:"admin_state"`
	Mtu         int64  `json:"mtu"`
	Speed       int64  `json:"speed"`
	MacAddress  string `json:"mac_address"`
	IfType      string `json:"if_type"`
	NetboxRefId int64  `json:"netbox_id,omitempty"`
	Blob        string `json:"blob,omitempty"`
}

type DbDevice struct {
	Id           string `json:"id,omitempty"`
	Policy       string `json:"policy,omitempty"`
	SerialNumber string `json:"serial_number"`
	Namespace    string `json:"namespace"`
	Hostname     string `json:"hostname"`
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
	Name        string `json:"name"`
	State       string `json:"state"`
	NetboxRefId int64  `json:"netbox_id,omitempty"`
	Blob        string `json:"blob,omitempty"`
}

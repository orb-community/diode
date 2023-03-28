package storage

type Service interface {
	Save(policy string, jsonData map[string]interface{}) (string, error)
}

type DbInterface struct {
	Id          string
	Policy      string
	Namespace   string
	Hostname    string
	Name        string
	AdminState  string
	Mtu         int64
	Speed       int64
	MacAddress  string
	IfType      string
	NetboxRefId int64
	Blob        string
}

type DbDevice struct {
	Id           string
	Policy       string
	SerialNumber string
	Namespace    string
	Hostname     string
	Model        string
	State        string
	Vendor       string
	NetboxRefId  int64
	Blob         string
}

type DbVlan struct {
	Id          string
	Policy      string
	Namespace   string
	Hostname    string
	Name        string
	State       string
	NetboxRefId int64
	Blob        string
}

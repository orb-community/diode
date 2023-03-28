package storage

type Service interface {
	Save(policy string, jsonData map[string]interface{}) (string, error)
}

type dbInterface struct {
	Id          string
	Policy      string // id
	Namespace   string // id
	Hostname    string // id
	NetboxRefId int64
	IfName      string
	IfType      string
}

type dbDevice struct {
	Id        string
	Policy    string
	Type      string
	Namespace string
	Hostname  string
	Vendor    string
}

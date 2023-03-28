package storage

type Service interface {
	Save(policy string, jsonData map[string]interface{}) error
}

type dbInterface struct {
	Id        string
	Policy    string
	IfName    string
	IfType    string
	Namespace string
	Hostname  string
	Address   string
	Vendor    string
	Os        string
}

type dbDevice struct {
	Id        string
	Policy    string
	Namespace string
	Hostname  string
	IpAddress string
}

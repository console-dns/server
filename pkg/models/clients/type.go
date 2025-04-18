package clients

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type ClientType string

var (
	TypeClient ClientType = "client"
	TypeDDNS   ClientType = "ddns"
)

func (t *ClientType) UnmarshalYAML(value *yaml.Node) error {
	s := value.Value
	if s == "" {
		*t = TypeClient
	} else {
		*t = ClientType(s)
	}
	return nil
}

func (t *ClientType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*t = TypeClient
	} else {
		*t = ClientType(s)
	}
	return nil
}

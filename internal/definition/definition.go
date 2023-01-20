package definition

import (
	"github.com/microavia/go-messgen/internal/config"
)

type Definition struct {
	Module    config.Module
	Proto     Proto
	Constants []Constant
	Messages  []Message
	Service   Service
}

type Proto struct {
	ProtoID uint8 `json:"proto_id"`
}

type Constant struct {
	Name     string          `json:"name"`
	BaseType string          `json:"basetype"`
	Fields   []ConstantField `json:"fields"`
}

type ConstantField struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"descr"`
}

type Message struct {
	Name        string         `json:"name"`
	ID          int            `json:"id"`
	Fields      []MessageField `json:"fields"`
	Description string         `json:"descr"`
}

type MessageField struct {
	Name        string    `json:"name"`
	Type        FieldType `json:"type"`
	Description string    `json:"descr"`
}

type Service struct {
	Serving []ServicePair `json:"serving"`
}

type ServicePair struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}

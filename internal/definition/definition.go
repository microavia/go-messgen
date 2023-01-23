package definition

import (
	"github.com/microavia/go-messgen/internal/config"
)

type Definition struct {
	Module   config.Module
	Proto    Proto
	Enums    []Enum
	Messages []Message
	Service  Service
}

type Proto struct {
	ProtoID uint8 `json:"proto_id"`
}

type Message struct {
	Name        string         `json:"name"`
	ID          int            `json:"id"`
	Fields      []MessageField `json:"fields"`
	Description string         `json:"descr"`
	MinSize     int            `json:"min_size"`
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

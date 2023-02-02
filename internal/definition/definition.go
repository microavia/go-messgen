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
	ID uint8 `json:"proto_id"`
}

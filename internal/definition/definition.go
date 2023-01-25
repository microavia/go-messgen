package definition

import (
	"github.com/microavia/go-messgen/internal/config"
)

type Definition struct {
	Module   config.Module
	Proto    Proto
	Enums    Enums
	Messages Messages
	Service  Service
}

type Proto struct {
	ProtoID uint8 `json:"proto_id"`
}

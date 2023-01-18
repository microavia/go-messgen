package definition

type Definition struct {
	Proto     Proto
	Constants []Constant
	Messages  map[string]Message
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
	ID          int            `json:"id"`
	Fields      []MessageField `json:"fields"`
	Description string         `json:"descr"`
}

type MessageField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"descr"`
}

type Service struct {
	Serving map[string]string
	Sending map[string]string
}

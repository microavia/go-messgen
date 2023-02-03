package definition

type Message struct {
	Name        string         `json:"-"`
	ID          uint16         `json:"id"`
	Fields      []MessageField `json:"fields"`
	Description string         `json:"descr"`
}

type MessageField struct {
	Name        string    `json:"name"`
	Type        FieldType `json:"type"`
	Description string    `json:"descr"`
}

package definition

import (
	"encoding/json"
)

type Message struct {
	Name        string        `json:"-"`
	ID          uint16        `json:"id"`
	Fields      MessageFields `json:"fields"`
	Description string        `json:"descr"`
}

type Messages map[string]Message

func (m *Messages) UnmarshalJSON(b []byte) error {
	return MapUnmarshalJSON(
		b,
		(*map[string]Message)(m),
		func(v Message) (Message, error) { return v, nil },
		func(v Message) string { return v.Name },
	)
}

func (m *Messages) MarshalJSON() ([]byte, error) {
	return json.Marshal(MapToSlice(*m, func(s1, s2 string) bool { return s1 < s2 }))
}

type MessageField struct {
	Name        string    `json:"name"`
	Type        FieldType `json:"type"`
	Description string    `json:"descr"`
}

type MessageFields map[string]MessageField

func (f *MessageFields) UnmarshalJSON(b []byte) error {
	return MapUnmarshalJSON(
		b,
		(*map[string]MessageField)(f),
		func(v MessageField) (MessageField, error) { return v, nil },
		func(v MessageField) string { return v.Name },
	)
}

func (f *MessageFields) MarshalJSON() ([]byte, error) {
	return json.Marshal(MapToSlice(*f, func(s1, s2 string) bool { return s1 < s2 }))
}

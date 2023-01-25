package definition

import (
	"encoding/json"
)

type Service struct {
	Serving ServicePairs `json:"serving"`
	Sending ServicePairs `json:"sending"`
}

type ServicePairs map[string]ServicePair

func (p *ServicePairs) UnmarshalJSON(b []byte) error {
	return MapUnmarshalJSON(
		b,
		(*map[string]ServicePair)(p),
		func(v ServicePair) (ServicePair, error) { return v, nil },
		func(v ServicePair) string { return v.Request },
	)
}

func (p *ServicePairs) MarshalJSON() ([]byte, error) {
	return json.Marshal(MapToSlice(*p, func(s1, s2 string) bool { return s1 < s2 }))
}

type ServicePair struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}

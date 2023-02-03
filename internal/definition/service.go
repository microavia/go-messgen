package definition

type Service struct {
	Serving []ServicePair `json:"serving"`
	Sending []ServicePair `json:"sending"`
}

type ServicePair struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}

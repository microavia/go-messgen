package definition

type Enum struct {
	Name     string      `json:"name"`
	BaseType string      `json:"basetype"`
	Values   []EnumValue `json:"fields"`
}

type EnumValue struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"descr"`
}

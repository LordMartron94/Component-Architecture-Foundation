package transport

type MessagePayload struct {
	Action string     `json:"action"`
	Args   []Argument `json:"args"`
}

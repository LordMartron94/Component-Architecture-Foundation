package payloads

import (
	"github.com/component-architecture-foundation/protocols/common"
)

type EventPayload struct {
	Sender common.ComponentID `json:"sender"`
	Type   string             `json:"type"`
	Data   interface{}        `json:"data"`
}

func (requestPayload EventPayload) GetType() string {
	return "event"
}

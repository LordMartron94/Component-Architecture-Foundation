package payloads

import (
	"github.com/component-architecture-foundation/protocols/common"
)

type ResponsePayload struct {
	Sender    common.ComponentID `json:"sender"`
	RequestID string             `json:"requestId"`
	Result    interface{}        `json:"result"`
	Error     string             `json:"error"`
}

func (requestPayload ResponsePayload) GetType() string {
	return "response"
}

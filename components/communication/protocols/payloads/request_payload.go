package payloads

import (
	"github.com/component-architecture-foundation/protocols/common"
)

type RequestPayload struct {
	Sender common.ComponentID  `json:"sender"`
	Target *common.ComponentID `json:"target"`
	Action string              `json:"action"`
	Params []interface{}       `json:"params"`
}

func (requestPayload RequestPayload) GetType() string {
	return "request"
}

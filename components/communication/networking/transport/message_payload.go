package transport

import (
	"encoding/json"
	"fmt"
)

type MessagePayload struct {
	Action string     `json:"action"`
	Args   []Argument `json:"args"`
}

func MessagePayloadFromBytes(payload []byte) (MessagePayload, error) {
	var mp MessagePayload
	err := json.Unmarshal(payload, &mp)
	if err != nil {
		return MessagePayload{}, fmt.Errorf("something went wrong while unmarshalling payload")
	}
	return mp, nil
}

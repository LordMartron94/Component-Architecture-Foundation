package transport

import (
	"encoding/json"
	"fmt"
)

type MessagePayload struct {
	Action string     `json:"action"`
	Args   []Argument `json:"args"`
}

func (p MessagePayload) ToString() string {
	return fmt.Sprintf("Action: %s, Args: %v", p.Action, p.Args)
}

func MessagePayloadFromBytes(payload []byte) (MessagePayload, error) {
	var mp MessagePayload
	err := json.Unmarshal(payload, &mp)
	if err != nil {
		return MessagePayload{}, fmt.Errorf("something went wrong while unmarshalling payload")
	}
	return mp, nil
}

func MessagePayloadToBytes(mp MessagePayload) ([]byte, error) {
	marshalled, err := json.Marshal(mp)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while marshalling payload")
	}
	return marshalled, nil
}

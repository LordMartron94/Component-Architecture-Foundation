package transport

import "time"

type Message struct {
	Requester ComponentID    `json:"requester"`
	Target    ComponentID    `json:"target"`
	TimeSent  time.Time      `json:"time_sent"`
	Payload   MessagePayload `json:"payload"`
}

func NewMessage(requester ComponentID, target ComponentID, payload MessagePayload) *Message {
	return &Message{
		Requester: requester,
		Target:    target,
		TimeSent:  time.Now(),
		Payload:   payload,
	}
}

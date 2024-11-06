package transport

import "time"

type Message struct {
	Requester ComponentID    `json:"requester"`
	Target    ComponentID    `json:"target"`
	TimeSent  time.Time      `json:"time_sent"`
	Payload   MessagePayload `json:"payload"`
}

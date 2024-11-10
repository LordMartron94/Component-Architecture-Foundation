package transport

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Requester *ComponentID    `json:"requester"`
	Target    *ComponentID    `json:"target"`
	TimeSent  *time.Time      `json:"time_sent"`
	Payload   *MessagePayload `json:"payload"`
	UniqueID  *string         `json:"unique_id"`
}

func (m *Message) GetValues() map[string]interface{} {
	values := make(map[string]interface{})

	values["requester"] = m.Requester
	values["target"] = m.Target
	values["time_sent"] = m.TimeSent
	values["payload"] = m.Payload
	values["unique_id"] = m.UniqueID

	return values
}

func NewMessage(requester ComponentID, target ComponentID, payload MessagePayload) *Message {
	currentTime := time.Now()

	return &Message{
		Requester: &requester,
		Target:    &target,
		TimeSent:  &currentTime,
		Payload:   &payload,
		UniqueID:  GenerateUniqueID(),
	}
}

func GenerateUniqueID() *string {
	generatedUUID := uuid.New().String()
	return &generatedUUID
}

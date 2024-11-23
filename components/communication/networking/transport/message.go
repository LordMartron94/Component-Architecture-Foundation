package transport

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

type Message struct {
	SenderID *string         `json:"sender_id"`
	TimeSent *time.Time      `json:"time_sent"`
	Payload  *MessagePayload `json:"payload"`
	UniqueID *string         `json:"unique_id"`
	TargetID string          `json:"target_id"`
}

func NewMessage(requester ComponentID, payload MessagePayload, targetId string) *Message {
	currentTime := time.Now()

	return &Message{
		SenderID: &requester.ComponentUniqueID,
		TimeSent: &currentTime,
		Payload:  &payload,
		UniqueID: GenerateUniqueID(),
		TargetID: targetId,
	}
}

func GenerateUniqueID() *string {
	generatedUUID := uuid.New().String()
	return &generatedUUID
}

func (m *Message) ToString() string {
	return fmt.Sprintf("Message{SenderID: %s, TimeSent: %s, Payload: %s, UniqueID: %s, TargetID: %s}",
		*m.SenderID, m.TimeSent.Format(time.RFC3339), m.Payload.ToString(), *m.UniqueID, m.TargetID)
}

func (m *Message) GetComponentIDFromRegistrationMessage(logger *logging.HoornLogger) (*ComponentID, error) {
	if len(m.Payload.Args) < 3 {
		logger.Error("Message does not contain enough arguments to create ComponentID", false, shared.NetworkingComponentName)
		return nil, fmt.Errorf("message does not contain enough arguments to create ComponentID")
	}

	capabilities, err := NewCapabilitiesFromJSON(m.Payload.Args[2].Value)

	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing capabilities from JSON: %s", err.Error()), false, shared.NetworkingComponentName)
		return nil, err
	}

	for capabilityIndex := range capabilities {
		capability := &capabilities[capabilityIndex]
		signature := &capability.Signature

		valid, err := signature.IsValid()

		if !valid {
			errMsg := fmt.Sprintf("Invalid signature in capabilities (%s): %s", capability.Name, err.Error())

			logger.Error(errMsg, false, shared.NetworkingComponentName)
			return nil, fmt.Errorf(errMsg)
		}
	}

	cID := ComponentID{
		Title:             m.Payload.Args[0].Value,
		Version:           m.Payload.Args[1].Value,
		Capabilities:      capabilities,
		ComponentUniqueID: *m.SenderID,
	}

	return &cID, nil
}

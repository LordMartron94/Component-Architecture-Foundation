package coding

import (
	"encoding/json"
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type JsonMessageCoder struct {
	Logger *logging.HoornLogger
}

// Decode decodes a bytes json message object to an actual transport.Message object.
func (j JsonMessageCoder) Decode(message []byte) (transport.Message, error) {
	// Unmarshal JSON string into a map[string]interface{}
	j.Logger.Debug(fmt.Sprintf("Received (attempting to decode): '%s'", message), false, shared.NetworkingComponentName)

	var unmarshalledMessage transport.Message
	err := json.Unmarshal(message, &unmarshalledMessage)
	if err != nil {
		j.Logger.Error(fmt.Sprintf("There was an error unmarshalling the json: %s", err.Error()), false, shared.NetworkingComponentName)
		return transport.Message{}, err
	}

	// Create a new transport.Message object from the map[string]interface{}
	return unmarshalledMessage, nil
}

// Encode encodes a transport.Message object into a byte array json string object.
func (j JsonMessageCoder) Encode(message transport.Message) ([]byte, error) {
	marshalled, err := json.Marshal(message)

	if err != nil {
		j.Logger.Error(fmt.Sprintf("There was an error marshalling the json: %s", err.Error()), false, shared.NetworkingComponentName)
		return nil, err
	}

	return marshalled, nil
}

package json

import (
	"encoding/json"
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/protocols/payloads"
	strategy2 "github.com/component-architecture-foundation/protocols/serialization/json/strategy"
	"github.com/component-architecture-foundation/shared"
)

type JSONSerializer struct {
	Logger logging.HoornLogger
}

func (serializer *JSONSerializer) Serialize(data payloads.PayloadInterface) ([]byte, error) {
	serializer.Logger.Debug("Serializing data to JSON", false, shared.MainComponentName)

	temp := map[string]interface{}{
		"type":    data.GetType(),
		"payload": data,
	}

	encoded, err := json.Marshal(temp)
	if err != nil {
		serializer.Logger.Warn("Failed to serialize data to JSON: "+err.Error(), false, shared.MainComponentName)
		return nil, err
	}

	return encoded, nil
}

func (serializer *JSONSerializer) Deserialize(data []byte) (payloads.PayloadInterface, error) {
	serializer.Logger.Debug("Deserializing JSON data", false, shared.MainComponentName)

	// 1. Unmarshal into a generic map to get the type field
	var temp map[string]interface{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return nil, err
	}

	// 2. Get the type from the map
	payloadType, ok := temp["type"].(string)
	if !ok {
		serializer.Logger.Warn("Missing or invalid 'type' field in JSON data", false, shared.MainComponentName)
		return nil, fmt.Errorf("missing or invalid 'type' field in JSON data")
	}
	payloadData, ok := temp["payload"].(map[string]interface{}) // Extract payload data
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'payload' field in JSON data")
	}

	// 3. Serialize the payloadData back to bytes
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return nil, err
	}

	// 4. Deserialize into the correct type based on the type field
	strategy, err := strategy2.GetDeserializer{Logger: serializer.Logger}.Get(payloadType)
	return strategy.Deserialize(payloadBytes)

}

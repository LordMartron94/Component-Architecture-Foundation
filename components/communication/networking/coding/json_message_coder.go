package coding

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type JsonMessageCoder struct {
	Logger *logging.HoornLogger
}

type MissingValueError struct{}

func (e MissingValueError) Error() string {
	return "missing value for pointer"
}

// Decode decodes a bytes json message object to an actual transport.Message object.
// Returns error if there are missing values for pointers.
func (j *JsonMessageCoder) Decode(message []byte) (transport.Message, error) {
	// Unmarshal JSON string into a map[string]interface{}
	j.Logger.Debug(fmt.Sprintf("Received (attempting to decode): '%s'", message), false, shared.NetworkingComponentName)

	var unmarshalledMessage transport.Message
	err := json.Unmarshal(message, &unmarshalledMessage)
	if err != nil {
		j.Logger.Error(fmt.Sprintf("There was an error unmarshalling the json: %s", err.Error()), false, shared.NetworkingComponentName)
		return transport.Message{}, err
	}

	missingData := j.checkForMissingValues(unmarshalledMessage)

	if missingData {
		return unmarshalledMessage, MissingValueError{}
	}

	// Create a new transport.Message object from the map[string]interface{}
	return unmarshalledMessage, nil
}

// Encode encodes a transport.Message object into a byte array json string object.
func (j *JsonMessageCoder) Encode(message transport.Message) ([]byte, error) {
	marshalled, err := json.Marshal(message)

	if err != nil {
		j.Logger.Error(fmt.Sprintf("There was an error marshalling the json: %s", err.Error()), false, shared.NetworkingComponentName)
		return nil, err
	}

	return marshalled, nil
}

func (j *JsonMessageCoder) checkForMissingValues(message transport.Message) bool {
	missingValues := false
	v := reflect.ValueOf(message)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		key := v.Type().Field(i).Name

		j.Logger.Debug(fmt.Sprintf("Value for key '%s': %v", key, field.Interface()), false, shared.NetworkingComponentName)

		// Check for nil pointers and zero values
		if field.Kind() == reflect.Ptr && field.IsNil() {
			j.Logger.Warn(fmt.Sprintf("Missing value for key '%s'", key), false, shared.NetworkingComponentName)
			missingValues = true
		}
	}

	return missingValues
}

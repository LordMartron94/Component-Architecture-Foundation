package transport

import (
	"fmt"
	"reflect"
)

type Signature struct {
	NumOfArgs     int      `json:"num_of_args"`
	TypeOfArgs    []string `json:"type_of_args"`
	ResponseTypes []string `json:"response_types"`
}

func (s *Signature) Equal(signature Signature) bool {
	return s.NumOfArgs == signature.NumOfArgs &&
		reflect.DeepEqual(s.TypeOfArgs, signature.TypeOfArgs) &&
		reflect.DeepEqual(s.ResponseTypes, signature.ResponseTypes)
}

func (s *Signature) IsValid() (bool, error) {
	if !(s.NumOfArgs >= 0) {
		return false, fmt.Errorf("number of arguments must be non-negative")
	}

	if !(len(s.TypeOfArgs) == s.NumOfArgs) {
		return false, fmt.Errorf("type of arguments length must match number of arguments")
	}

	if !(len(s.ResponseTypes) >= 1) {
		return false, fmt.Errorf("at least one response type must be provided")
	}

	return true, nil
}

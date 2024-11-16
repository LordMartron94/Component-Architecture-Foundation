package transport

import "reflect"

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

func (s *Signature) IsValid() bool {
	return s.NumOfArgs >= 0 && len(s.TypeOfArgs) == s.NumOfArgs && len(s.ResponseTypes) >= 1
}

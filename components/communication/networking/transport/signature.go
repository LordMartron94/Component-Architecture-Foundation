package transport

import "reflect"

type Signature struct {
	NumOfArgs  int      `json:"num_of_args"`
	TypeOfArgs []string `json:"type_of_args"`
}

func (s Signature) Equal(signature Signature) bool {
	return s.NumOfArgs == signature.NumOfArgs &&
		reflect.DeepEqual(s.TypeOfArgs, signature.TypeOfArgs)
}

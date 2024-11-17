package transport

import (
	"encoding/json"
)

type Capability struct {
	Name      string    `json:"name"`
	Signature Signature `json:"signature"`
}

func (c *Capability) Equal(other Capability) bool {
	nameEqual := c.Name == other.Name
	signatureEqual := c.Signature.Equal(other.Signature)

	return nameEqual && signatureEqual
}

func NewCapabilitiesFromJSON(jsonStr string) ([]Capability, error) {
	var capabilities []Capability
	err := json.Unmarshal([]byte(jsonStr), &capabilities)
	if err != nil {
		return nil, err
	}
	return capabilities, nil
}

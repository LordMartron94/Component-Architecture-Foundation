package transport

type Capability struct {
	Name      string    `json:"name"`
	Signature Signature `json:"signature"`
}

func (c *Capability) Equal(other Capability) bool {
	nameEqual := c.Name == other.Name
	signatureEqual := c.Signature.Equal(other.Signature)

	return nameEqual && signatureEqual
}

package transport

type ComponentID struct {
	Title        string       `json:"title"`
	Version      string       `json:"version"`
	Capabilities []Capability `json:"capabilities"`
}

func (c *ComponentID) Equal(other ComponentID) bool {
	if len(c.Capabilities) != len(other.Capabilities) {
		return false
	}

	titleEqual := c.Title == other.Title
	versionEqual := c.Version == other.Version

	for i := range c.Capabilities {
		if !c.Capabilities[i].Equal(other.Capabilities[i]) {
			return false
		}
	}

	return titleEqual && versionEqual
}

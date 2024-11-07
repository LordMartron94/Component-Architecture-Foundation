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

// HasCapabilityAction checks if the component contains the ability to do this action.
// WARNING: Doesn't take into account signature, only action name.
func (c *ComponentID) HasCapabilityAction(action string) bool {
	for _, capability := range c.Capabilities {
		if capability.Name == action {
			return true
		}
	}
	return false
}

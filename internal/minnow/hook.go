package minnow

type Hook interface {
	Matches([]byte) bool
}

type BasicPropertiesMatchHook struct {
	match Properties
}

func (hook BasicPropertiesMatchHook) Matches(matchAgainst []byte) bool {
	properties, err := BytesToProperties(matchAgainst)

	if err != nil {
		return false
	}

	for expectedKey, expectedValue := range hook.match {
		if value, found := properties[expectedKey]; found {
			if value != expectedValue {
				return false
			}
		}
	}

	return true
}

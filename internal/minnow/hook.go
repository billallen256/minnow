package minnow

type Hook interface {
	MatchesBytes([]byte) bool
	Matches(Properties) bool
}

type BasicPropertiesMatchHook struct {
	match Properties
}

func NewBasicPropertiesMatchHookFromFile(path Path) (BasicPropertiesMatchHook, error) {
	hookProperties, err := PropertiesFromFile(path)

	if err != nil {
		return BasicPropertiesMatchHook{}, err
	}

	return BasicPropertiesMatchHook{hookProperties}, nil
}

func (hook BasicPropertiesMatchHook) Matches(matchAgainst Properties) bool {
	for expectedKey, expectedValue := range hook.match {
		if value, found := matchAgainst[expectedKey]; found {
			if value != expectedValue {
				return false
			}
		}
	}

	return true
}

func (hook BasicPropertiesMatchHook) MatchesBytes(matchAgainst []byte) bool {
	properties, err := BytesToProperties(matchAgainst)

	if err != nil {
		return false
	}

	return hook.Matches(properties)
}

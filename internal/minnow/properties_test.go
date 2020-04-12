package minnow

import (
	"bytes"
	"testing"
)

func TestBytesToProperties(t *testing.T) {
	propStr := "foo = bar\nbaz=2\nname = continue"
	propBytes := bytes.NewBufferString(propStr).Bytes()
	properties, err := BytesToProperties(propBytes)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if len(properties) != 3 {
		t.Errorf("Expected 3 properties, but only got %d", len(properties))
	}

	expectedProperties := map[string]string{
		"foo":  "bar",
		"baz":  "2",
		"name": "continue",
	}

	for expectedKey, expectedValue := range expectedProperties {
		if value, found := properties[expectedKey]; found {
			if value != expectedValue {
				t.Errorf("Found key %s, but value %s", expectedKey, value)
			}
		} else {
			t.Errorf("Missing property %s", expectedKey)
		}
	}
}

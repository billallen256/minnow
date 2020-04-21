package minnow

import (
	"bytes"
	"testing"
)

func TestBasicPropertiesMatchHookMatch(t *testing.T) {
	hookPropertiesStr := "type = foo"
	hookPropertiesBytes := bytes.NewBufferString(hookPropertiesStr).Bytes()
	hookProperties, err := BytesToProperties(hookPropertiesBytes)

	if err != nil {
		t.Errorf(err.Error())
	}

	hook := BasicPropertiesMatchHook{hookProperties}

	propertiesStr := "type=foo\nvolume=11\nbar = baz"
	propertiesBytes := bytes.NewBufferString(propertiesStr).Bytes()

	if !hook.MatchesBytes(propertiesBytes) {
		t.Errorf("Hook should have matched")
	}
}

func TestBasicPropertiesMatchHookNoMatch(t *testing.T) {
	hookPropertiesStr := "type = foo\nvolume=11"
	hookPropertiesBytes := bytes.NewBufferString(hookPropertiesStr).Bytes()
	hookProperties, err := BytesToProperties(hookPropertiesBytes)

	if err != nil {
		t.Errorf(err.Error())
	}

	hook := BasicPropertiesMatchHook{hookProperties}

	propertiesStr := "type=foo\nvolume=12\nbar = baz"
	propertiesBytes := bytes.NewBufferString(propertiesStr).Bytes()

	if hook.MatchesBytes(propertiesBytes) {
		t.Errorf("Hook should not have matched")
	}

}

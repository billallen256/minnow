package minnow

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	PropertiesExtension = ".properties"
)

type Properties map[string]string

func BytesToProperties(input []byte) (Properties, error) {
	inputStr := bytes.NewBuffer(input).String()
	lines := make([]string, 0)
	properties := make(Properties)
	errorList := make([]string, 0)

	for _, line := range strings.Split(inputStr, "\n") {
		line = strings.TrimSpace(line)

		if len(line) > 0 {
			lines = append(lines, line)
		}
	}

	for _, line := range lines {
		parts := strings.Split(line, "=")

		if len(parts) == 2 {
			name := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			properties[name] = value
		} else {
			errorList = append(errorList, fmt.Sprintf("Invalid property: %s", line))
		}
	}

	if len(errorList) > 0 {
		return properties, errors.New(strings.Join(errorList, "; "))
	}

	return properties, nil
}

func PropertiesFromFile(path Path) (Properties, error) {
	propertiesBytes, err := path.ReadBytes()

	if err != nil {
		return nil, err
	}

	return BytesToProperties(propertiesBytes)
}

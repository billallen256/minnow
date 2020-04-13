package minnow

import (
	"bytes"
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	configPropertiesStr := `
		ingest_dir=/tmp
		ingest_min_age=600
		work_dir=/var
		work_age_off=86400
		processor_definitions_dir=/usr`
	configPropertiesBytes := bytes.NewBufferString(configPropertiesStr).Bytes()
	configProperties, err := BytesToProperties(configPropertiesBytes)

	if err != nil {
		t.Errorf(err.Error())
	}

	config, err := ParseConfig(configProperties)

	if err != nil {
		t.Errorf(err.Error())
	}

	if config.IngestPath != "/tmp" {
		t.Errorf("Incorrect IngestPath")
	}

	if config.IngestMinAge != time.Duration(600) * time.Second {
		t.Errorf("Incorrect IngestMinAge")
	}

	if config.WorkPath != "/var" {
		t.Errorf("Incorrect WorkPath")
	}

	if config.WorkAgeOff != time.Duration(86400) * time.Second {
		t.Errorf("Incorrect WorkAgeOff")
	}

	if config.ProcessorDefinitionsPath != "/usr" {
		t.Errorf("Incorrect ProcessorDefinitionsPath")
	}
}

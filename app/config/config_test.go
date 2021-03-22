package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c, err := LoadConfig("../../")

	if err != nil {
		t.Errorf("Unexpected error loading the config file")
	}

	if &c == nil {
		t.Errorf("Unexpected error loading the config file")
	}
}
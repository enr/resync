package main

import (
	"testing"
)

func TestLoadConfigurationFromBytes(t *testing.T) {
	data := `
folders:
    pictures:
        local_path: ~/Pictures
        external_path: Pictures
`
	conf, err := loadConfigurationFromBytes([]byte(data))
	if err != nil {
		t.Errorf(`error loading conf: %v`, err)
	}
	if len(conf.Folders) != 1 {
		t.Errorf(`folders number. expected %d got %d`, 1, len(conf.Folders))
	}
}

func TestLoadConfigurationFromFile(t *testing.T) {
	conf, err := loadConfigurationFromFile("../../testdata/configurations/01.yml")
	if err != nil {
		t.Errorf(`unexpected error loading conf: %v`, err)
	}
	if len(conf.Folders) != 1 {
		t.Errorf(`folders number. expected %d got %d`, 1, len(conf.Folders))
	}
}

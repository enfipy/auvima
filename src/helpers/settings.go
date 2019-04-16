package helpers

import (
	"io/ioutil"

	"github.com/enfipy/auvima/src/config"

	yaml "gopkg.in/yaml.v2"
)

func GetSettings(path string) *config.Settings {
	settings := new(config.Settings)

	yamlFile, err := ioutil.ReadFile(path)
	PanicOnError(err)
	err = yaml.Unmarshal(yamlFile, settings)
	PanicOnError(err)

	return settings
}

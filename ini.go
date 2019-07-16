package decouple

import (
	"regexp"

	"gopkg.in/ini.v1"
)

// repositoryIni retrieves option keys from .ini files.
type repositoryIni struct{}

// Get implements the interface repository that returns the value of the key in the settings.ini file
func (r repositoryIni) Get(option string, def, cast interface{}) interface{} {
	var value interface{}

	file, err := ini.Load("settings.ini")
	if err != nil {
		return nil
	}

	key, err := file.Section("settings").GetKey(option)
	if err != nil {
		return nil
	}

	value = key.Value()

	ok, _ := regexp.MatchString(`^\%\(.*\)s$`, value.(string))
	if ok {
		value = key.String()
	}

	if cast != nil {
		value, err = castValue(value.(string), cast.(string))
		if err != nil {
			panic(err)
		}
		return value
	}

	return value
}

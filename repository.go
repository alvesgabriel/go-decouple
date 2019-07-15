package decouple

import (
	"os"
)

// repository is a interface tha has func to get the value of a config
type repository interface {
	Get(string, interface{}, interface{}) interface{}
}

// repositoryEmpty retrieves option keys from .ini files.
type repositoryEmpty struct{}

// Get implements the interface repository that returns the value of the key variable enviroment
func (r repositoryEmpty) Get(option string, def, cast interface{}) interface{} {
	valueStr, ok := os.LookupEnv(option)
	if ok == false {
		return nil
	}

	if cast == "bool" && valueStr != "" {
		value, err := castBoolean(valueStr)
		if err != nil {
			panic(err)
		}
		return *value
	}

	return valueStr
}

package decouple

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	boolean = map[string]bool{
		"1":     true,
		"yes":   true,
		"true":  true,
		"on":    true,
		"0":     false,
		"no":    false,
		"false": false,
		"off":   false,
		"":      false,
	}
)

// castBoolean helps to convert config values to boolean as ConfigParser do.
func castBoolean(value string) (*bool, error) {
	if v, ok := boolean[strings.ToLower(value)]; ok {
		return &v, nil
	}
	return nil, errors.Errorf("Not a boolean: %s", value)
}

// Config function improve decouple's usability now just import decouple and use Config and start using with no configuration.
func Config(option string, def, cast interface{}) interface{} {
	var r repository

	r = repositoryEmpty{}
	value := r.Get(option, def, cast)
	if value != nil {
		return value
	}

	r = repositoryEnv{}
	value = r.Get(option, def, cast)
	if value != nil {
		return value
	}

	r = repositoryIni{}
	value = r.Get(option, def, cast)
	if value != nil {
		return value
	}

	if def != nil {
		if cast == "bool" && reflect.TypeOf(def).Kind() != reflect.Bool {
			err := fmt.Sprintf("Can't cast '%v' to '%v'", def, cast)
			panic(err)
		}
		return def
	}

	panic(fmt.Sprintf("%s not found. Declare it as envvar or define a default value.", option))
}

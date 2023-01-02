package decouple

import (
	"fmt"

	"github.com/alvesgabriel/go-decouple/repository"
)

// Config function improve decouple's usability now just import decouple and use Config and start using with no configuration.
func Config[Type repository.EnvVar](key string, def interface{}) (Type, error) {
	var r repository.Repository
	var cast Type

	r = &repository.RepositoryEmpty{}
	value := r.Get(key)
	if value == "" {
		return cast, fmt.Errorf("")
	}

	return repository.Cast(value, cast)
	// switch fmt.Sprintf("%T", cast) {
	// case "int":
	// 	num, err := strconv.Atoi(value)
	// 	if err != nil {
	// 		return cast, err
	// 	}
	// 	return Type(num), nil
	// 	// case int64:
	// 	// 	return strconv.ParseInt(value, 10, 64)
	// 	// case float64:
	// 	// 	return strconv.ParseFloat(value, 64)
	// 	// default:
	// 	// 	return Type(value), nil
	// }

	// r = repositoryEnv{}
	// value = r.Get(option, def, cast)
	// if value != nil {
	// 	return value
	// }

	// r = repositoryIni{}
	// value = r.Get(option, def, cast)
	// if value != nil {
	// 	return value
	// }

	// if def != nil {
	// 	if cast == "bool" && reflect.TypeOf(def).Kind() != reflect.Bool {
	// 		err := fmt.Sprintf("Can't cast '%v' to '%v'", def, cast)
	// 		panic(err)
	// 	}
	// 	return def
	// }

	// panic(fmt.Sprintf("%s not found. Declare it as envvar or define a default value.", option))

	return cast, fmt.Errorf("%s not found. Declare it as envvar or define a default value.", key)
}

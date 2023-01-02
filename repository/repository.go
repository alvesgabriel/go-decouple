package repository

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

func Cast[Type EnvVar](value string, cast Type) (Type, error) {
	switch fmt.Sprintf("%T", cast) {
	case "bool":
		v, err := castBoolean(value)
		if err != nil {
			return Type(false), err
		}
		return Type(*v), err
	case "int":
		num, err := strconv.Atoi(value)
		if err != nil {
			return Type(0), err
		}
		return Type(num), err
	// case int64:
	// 	return strconv.ParseInt(value, 10, 64)
	// case float64:
	// 	return strconv.ParseFloat(value, 64)
	default:
		return cast, nil
	}
}

// castBoolean helps to convert config values to boolean as ConfigParser do.
func castBoolean(value string) (*bool, error) {
	if v, ok := boolean[strings.ToLower(value)]; ok {
		return &v, nil
	}
	return nil, fmt.Errorf("not a boolean: %s", value)
}

// RepositoryEmpty retrieves option keys from environment variables.
type RepositoryEmpty struct{}

// Get implements the interface repository that returns the value of the key variable enviroment
func (r *RepositoryEmpty) Get(option string) string {
	return os.Getenv(option)
}

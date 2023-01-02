package decouple

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type variable struct {
	option string
	value  interface{}
	def    interface{}
	cast   interface{}
}

func TestConfig(t *testing.T) {
	variables := []variable{
		{option: "KeyFallback", value: "On"},
		{option: "KeyTrue", value: "True"},
		{option: "KeyOne", value: "1"},
		{option: "KeyYes", value: "yes"},
		{option: "KeyOn", value: "on"},
	}
	for _, v := range variables {
		os.Setenv(v.option, v.value.(string))
		if value := Config(v.option, nil, nil); value != v.value {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, value, value)
		}
		os.Unsetenv(v.option)
	}
}

func TestConfigTrue(t *testing.T) {
	variables := []variable{
		{option: "KeyFallback", value: "On"},
		{option: "KeyTrue", value: "True"},
		{option: "KeyOne", value: "1"},
		{option: "KeyYes", value: "yes"},
		{option: "KeyOn", value: "on"},
	}
	for _, v := range variables {
		os.Setenv(v.option, v.value.(string))
		if value := Config(v.option, nil, "bool"); value != boolean[strings.ToLower(v.value.(string))] {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, boolean[strings.ToLower(v.value.(string))], value)
		}
		os.Unsetenv(v.option)
	}
}

func TestCastBooleanTrue(t *testing.T) {
	trues := []string{"1", "yes", "YES", "true", "TRUE", "on", "ON"}
	for _, value := range trues {
		if boolCast, _ := castBoolean(value); !*boolCast {
			t.Errorf("'%v' wait '%v' got '%v'", value, true, boolCast)
		}
	}
}

func TestCastBooleanFalse(t *testing.T) {
	falses := []string{"0", "no", "NO", "false", "FALSE", "off", "OFF", ""}
	for _, value := range falses {
		if boolCast, _ := castBoolean(value); *boolCast {
			t.Errorf("'%v' wait '%v' got '%v'", value, false, boolCast)
		}
	}
}

func TestCastBooleanError(t *testing.T) {
	values := []string{"this", "is", "not", "boolean", "cast"}
	for _, value := range values {
		if _, err := castBoolean(value); err == nil {
			wait := fmt.Sprintf("Not a boolean: %s", value)
			t.Errorf("'%v' wait '%v' got '%v'", value, wait, err)
		}
	}
}

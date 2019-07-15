package decouple

import (
	"os"
	"testing"

	"github.com/Flaque/filet"
)

const (
	filenameIni = "settings.ini"
	iniFile     = `
[settings]
KeyTrue=True
KeyOne=1
KeyYes=yes
KeyOn=on
KeyFalse=False
KeyZero=0
KeyNo=no
KeyOff=off
KeyEmpty=
#CommentedKey=None
PercentIsEscaped=%%
Interpolation=%(KeyOff)s
IgnoreSpace = text
RespectSingleQuoteSpace = ' text'
RespectDoubleQuoteSpace = " text"
KeyOverrideByEnv=NotThis
`
)

func TestConfigIniFile(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	variables := []variable{
		{option: "KeyTrue", value: "True"},
		{option: "KeyOne", value: "1"},
		{option: "KeyYes", value: "yes"},
		{option: "KeyOn", value: "on"},
	}
	for _, v := range variables {
		if value := Config(v.option, v.def, v.cast); value != v.value {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, v.value, value)
		}
	}
}

func TestIni(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	keys := []string{"CommentedKey", "UndefinedKey"}

	for _, key := range keys {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("%v is defined", key)
				}
			}()

			Config(key, nil, nil)
		}()
	}
}

func TestIniInterpolation(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	v := variable{option: "Interpolation", value: "off"}
	if value := Config(v.option, v.def, v.cast); value != v.value {
		t.Errorf("'%v' wait '%v' got '%v'", v.option, v.value, value)
	}
}

func TestIniBooleanTrue(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	variables := []variable{
		{option: "KeyTrue", value: true, cast: "bool"},
		{option: "KeyOne", value: true, cast: "bool"},
		{option: "KeyYes", value: true, cast: "bool"},
		{option: "KeyOn", value: true, cast: "bool"},
	}
	for _, v := range variables {
		if value := Config(v.option, v.def, v.cast); value != v.value {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, v.value, value)
		}
	}
}

func TestIniBooleanFalse(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	variables := []variable{
		{option: "KeyFalse", value: false, cast: "bool"},
		{option: "KeyZero", value: false, cast: "bool"},
		{option: "KeyNo", value: false, cast: "bool"},
		{option: "KeyOff", value: false, cast: "bool"},
	}
	for _, v := range variables {
		if value := Config(v.option, v.def, v.cast); value != v.value {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, v.value, value)
		}
	}
}

func TestIniDefaultBool(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	variables := []variable{
		{option: "UndefinedKey", def: false, cast: "bool"},
		{option: "UndefinedKey", def: true, cast: "bool"},
	}
	for _, v := range variables {
		if value := Config(v.option, v.def, v.cast); value != v.def {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, v.def, value)
		}
	}
}

func TestIniDefault(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	variables := []variable{
		{option: "UndefinedKey", def: false, cast: nil},
		{option: "UndefinedKey", def: true, cast: nil},
	}
	for _, v := range variables {
		if value := Config(v.option, v.def, v.cast); value != v.def {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, v.def, value)
		}
	}
}

func TestIniDefaultInvalidBool(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	v := variable{option: "UndefinedKey", def: "NotBool", cast: "bool"}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Panic of cast in %v don't work", v.option)
			}
		}()

		Config(v.option, v.def, v.cast)
	}()
}

func TestIniSupportSpace(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	variables := []variable{
		{option: "IgnoreSpace", value: "text"},
		{option: "RespectSingleQuoteSpace", value: " text"},
		{option: "RespectDoubleQuoteSpace", value: " text"},
	}
	for _, v := range variables {
		if value := Config(v.option, v.def, v.cast); value != v.value {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, value, value)
		}
	}
}

func TestIniOsEnviron(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	v := variable{option: "KeyOverrideByEnv", value: "This"}
	os.Setenv(v.option, v.value.(string))
	if v.value != Config(v.option, v.def, v.cast) {
		t.Errorf("%+v is not defined", v.option)
	}
	os.Unsetenv(v.option)
}

func TestIniUndefinedButPresentInOsEnviron(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	v := variable{option: "KeyOnlyEnviron", value: ""}
	os.Setenv(v.option, v.value.(string))
	if v.value != Config(v.option, v.def, v.cast) {
		t.Errorf("%+v is not defined", v.option)
	}
	os.Unsetenv(v.option)
}

func TestIniEmptyStringMeansFalse(t *testing.T) {
	filet.File(t, filenameIni, iniFile)
	defer filet.CleanUp(t)

	v := variable{option: "KeyEmpty", value: false}
	if value := Config(v.option, v.def, "bool"); value != v.value {
		t.Errorf("%s id not false", v.option)
	}
}

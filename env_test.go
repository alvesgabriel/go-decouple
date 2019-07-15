package decouple

import (
	"os"
	"testing"

	"github.com/Flaque/filet"
)

const (
	filenameEnv = ".env"
	envFile     = `
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
PercentNotEscaped=%%
NoInterpolation=%(KeyOff)s
IgnoreSpace = text
RespectSingleQuoteSpace = ' text'
RespectDoubleQuoteSpace = " text"
KeyOverrideByEnv=NotThis
`
)

func TestFileToMap(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	r := repositoryEnv{}
	fileMap := r.fileToMap(filenameEnv)

	variables := []variable{
		{option: "KeyTrue", value: "True"},
		{option: "KeyOne", value: "1"},
		{option: "KeyYes", value: "yes"},
		{option: "KeyOn", value: "on"},
		{option: "KeyFalse", value: "False"},
		{option: "KeyZero", value: "0"},
		{option: "KeyNo", value: "no"},
		{option: "KeyOff", value: "off"},
		{option: "KeyEmpty", value: ""},
		{option: "PercentNotEscaped", value: "%%"},
		{option: "NoInterpolation", value: "%(KeyOff)s"},
		{option: "IgnoreSpace", value: "text"},
		{option: "RespectSingleQuoteSpace", value: " text"},
		{option: "RespectDoubleQuoteSpace", value: " text"},
		{option: "KeyOverrideByEnv", value: "NotThis"},
	}
	for _, v := range variables {
		if value := fileMap[v.option]; value != v.value {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, v.value, value)
		}
	}
}

func TestFileToMapError(t *testing.T) {
	filet.CleanUp(t)

	r := repositoryEnv{}
	fileMap := r.fileToMap(filenameEnv)

	if fileMap != nil {
		t.Errorf("'%v' wait '%v' got '%v'", "fileToMap", nil, fileMap)
	}
}

func TestEnv(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	keys := []string{"CommentedKey", "UndefinedKey"}

	for _, key := range keys {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("%+v is defined", key)
				}
			}()

			Config(key, nil, nil)
		}()
	}
}

func TestEnvPercentNotEscaped(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	v := variable{option: "PercentNotEscaped", value: "%%"}
	if v.value != Config(v.option, v.def, v.cast) {
		t.Errorf("%s is escaped", v.option)
	}
}

func TestEnvNoInterpolation(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	v := variable{option: "NoInterpolation", value: "%(KeyOff)s"}
	if v.value != Config(v.option, v.def, v.cast) {
		t.Errorf("key is escaped")
	}
}

func TestConfigEnvFile(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	variables := []variable{
		{option: "KeyTrue", value: "True"},
		{option: "KeyOne", value: "1"},
		{option: "KeyYes", value: "yes"},
		{option: "KeyOn", value: "on"},
	}
	for _, v := range variables {
		if value := Config(v.option, v.def, v.cast); value != v.value {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, value, value)
		}
	}
}

func TestConfigEnvTrue(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	variables := []variable{
		{option: "KeyTrue", value: true},
		{option: "KeyOne", value: true},
		{option: "KeyYes", value: true},
		{option: "KeyOn", value: true},
	}
	for _, v := range variables {
		if value := Config(v.option, v.def, "bool"); value != v.value {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, value, value)
		}
	}
}

func TestConfigEnvFalse(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	variables := []variable{
		{option: "KeyFalse", value: false},
		{option: "KeyZero", value: false},
		{option: "KeyNo", value: false},
		{option: "KeyOff", value: false},
	}
	for _, v := range variables {
		if value := Config(v.option, v.def, "bool"); value != v.value {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, value, value)
		}
	}
}

func TestEnvOsEnviron(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	v := variable{option: "KeyOverrideByEnv", value: "This"}
	os.Setenv(v.option, v.value.(string))
	if v.value != Config(v.option, v.def, v.cast) {
		t.Errorf("%+v is not defined", v.option)
	}
	os.Unsetenv(v.option)
}

func TestEnvUndefinedButPresentInOsEnviron(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	v := variable{option: "KeyOnlyEnviron", value: ""}
	os.Setenv(v.option, v.value.(string))
	if v.value != Config(v.option, v.def, v.cast) {
		t.Errorf("%+v is not defined", v.option)
	}
	os.Unsetenv(v.option)
}

func TestEnvEmpty(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	v := variable{option: "KeyEmpty", value: ""}
	if value := Config(v.option, v.def, v.cast); value != v.value {
		t.Errorf("%s is not empty string", v.option)
	}
}

func TestEnvSupportSpace(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
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

func TestEnvEmptyStringMeansFalse(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	v := variable{option: "KeyEmpty", value: false}
	if value := Config(v.option, v.def, "bool"); value != v.value {
		t.Errorf("%s id not false", v.option)
	}
}

func TestEnvDefaultValue(t *testing.T) {
	filet.File(t, filenameEnv, envFile)
	defer filet.CleanUp(t)

	variables := []variable{
		{option: "KeyDefaultString", def: "This", cast: nil},
		{option: "KeyDefaultInt", def: 42, cast: nil},
		{option: "KeyDefaultBoolTrue", def: true, cast: "bool"},
		{option: "KeyDefaultBoolFalse", def: false, cast: "bool"},
		{option: "KeyDefaultFloat", def: 3.14, cast: nil},
	}

	for _, v := range variables {
		if value := Config(v.option, v.def, v.cast); value != v.def {
			t.Errorf("'%v' wait '%v' got '%v'", v.option, v.def, value)
		}
	}

}

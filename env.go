package decouple

import (
	"io/ioutil"
	"log"
	"strings"
)

// repositoryEnv retrieves option keys from .env files with fall back to os.environ.
type repositoryEnv struct {
	Data map[string]string
}

// Get implements the interface repository that returns the value of the key in the .env file
func (r repositoryEnv) Get(option string, def, cast interface{}) interface{} {
	r.Data = r.fileToMap(".env")

	val, ok := r.Data[option]
	if ok == false {
		return nil
	}

	if cast == "bool" {
		value, err := castBoolean(val)
		if err != nil {
			panic(err)
		}
		return *value
	}

	return val
}

func (r *repositoryEnv) fileToMap(filename string) map[string]string {
	var data = make(map[string]string)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		if err.Error() == "open .env: no such file or directory" {
			return nil
		}
		log.Fatal(err)
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		params := strings.SplitN(line, "=", 2)
		k, v := params[0], params[1]
		k = strings.TrimSpace(k)
		v = strings.Trim(strings.TrimSpace(v), `'"`)
		data[k] = v
	}

	return data
}

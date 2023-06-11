package config

import (
	"os"
	"strings"
)

func GetEnvironment() map[string]string {
	penv := make(map[string]string)
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		key := strings.ToUpper(kv[0])
		if strings.HasPrefix(key, "PARTSY") {
			penv[key] = kv[1]
		}
	}

	return penv
}

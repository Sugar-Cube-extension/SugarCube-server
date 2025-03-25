package utils

import (
	"errors"
	"os"
	"strconv"
)

type CliVar interface {
	int | int64 | string | bool
}

// Used to decide what to use as variables.
func CheckForEnv[T CliVar](envVar string, cliVar T) (T, error) {
	envVarResult := os.Getenv(envVar)
	if envVarResult == "" { // Env vars will take priority
		return cliVar, nil
	}

	var result T
	var err error

	switch any(cliVar).(type) {
	case int:
		var intVal int
		intVal, err = strconv.Atoi(envVarResult)
		result = any(intVal).(T)
	case int64:
		var int64Val int64
		int64Val, err = strconv.ParseInt(envVarResult, 10, 64)
		result = any(int64Val).(T)
	case bool:
		var boolVal bool
		boolVal, err = strconv.ParseBool(envVarResult) //Parses 90% of possible bool naming schemas
		result = any(boolVal).(T)
	case string:
		result = any(envVarResult).(T)
	default:
		err = errors.New("unsupported type")
	}
	return result, err

}

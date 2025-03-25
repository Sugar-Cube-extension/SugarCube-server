package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type CliVar interface {
	int | int64 | string | bool
}

const (
	EnvDBPort     = "SUGARCUBE_DB_PORT"
	EnvPort       = "SUGARCUBE_PORT"
	EnvDBURI      = "SUGARCUBE_DB_URI"
	EnvDBUser     = "SUGARCUBE_DB_USER"
	EnvDBPassword = "SUGARCUBE_DB_PASSWORD"
	EnvDebug      = "SUGARCUBE_DEBUG"
	UriProtocol   = "mongodb://"
)

type SessionCtx struct {
	DbPort     int64
	ServerPort int64
	DbUri      string
	DbUser     string
	DbPassword string
	Debug      bool
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

func (s SessionCtx) PrintEnv() {
	fmt.Println("###SugarCube Configuration###")
	fmt.Printf("  Database Port    : %d\n", s.DbPort)
	fmt.Printf("  Server Port      : %d\n", s.ServerPort)
	fmt.Printf("  Database URI     : %s\n", s.DbUri)
	if s.DbUser != "" {
		fmt.Printf("  Database User    : %s\n", s.DbUser)
	} else {
		fmt.Println("  Database User    : [not set]")
	}
	if s.DbPassword != "" {
		fmt.Println("  Database Password: [hidden]") // We obv hide the password for security reasons
	} else {
		fmt.Println("  Database Password: [not set]")
	}
	fmt.Printf("  Debug Mode       : %t\n", s.Debug)

}

func (s SessionCtx) GetFullUri() string {
	var buffer bytes.Buffer
	buffer.WriteString(UriProtocol)
	if s.DbUser != "" {
		buffer.WriteString(s.DbUser)
		buffer.WriteString(":")
	}
	if s.DbPassword != "" {
		buffer.WriteString(s.DbPassword)
		buffer.WriteString("@")
	}
	buffer.WriteString(s.DbUri)
	buffer.WriteString(":")
	buffer.WriteString(strconv.FormatInt(s.DbPort, 10))
	return buffer.String()
}

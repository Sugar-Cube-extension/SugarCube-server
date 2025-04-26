package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type CliVar interface {
	int | uint64 | string | bool
}

const (
	EnvDBPort     = "SUGARCUBE_DB_PORT"
	EnvPort       = "SUGARCUBE_PORT"
	EnvDBURI      = "SUGARCUBE_DB_URI"
	EnvDBUser     = "SUGARCUBE_DB_USER"
	EnvDBPassword = "SUGARCUBE_DB_PASSWORD"
	EnvDebug      = "SUGARCUBE_DEBUG"
	UriProtocol   = "mongodb://"

	//Cool colors
	ColorReset  = "\033[0m"
	ColorCyan   = "\033[36m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorBold   = "\033[1m"
)

type SessionCtx struct {
	DbPort     uint16
	ServerPort uint16
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
	case uint64:
		var uint64Val uint64
		uint64Val, err = strconv.ParseUint(envVarResult, 10, 64)
		result = any(uint64Val).(T)
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
func (s *SessionCtx) IsEmpty() bool {
	if s == nil {
		return true
	}
	return s.DbPort == 0 && s.ServerPort == 0 && s.DbUri == ""

}

func (s SessionCtx) PrintEnv() {
	fmt.Println(ColorCyan + "########################################" + ColorReset)
	fmt.Println(ColorCyan + "#       " + ColorBold + "SugarCube Configuration" + ColorReset + ColorCyan + "      #" + ColorReset)
	fmt.Println(ColorCyan + "########################################" + ColorReset)

	fmt.Printf(ColorGreen+"  %-18s:"+ColorReset+" %d\n", "Database Port", s.DbPort)
	fmt.Printf(ColorGreen+"  %-18s:"+ColorReset+" %d\n", "Server Port", s.ServerPort)
	fmt.Printf(ColorGreen+"  %-18s:"+ColorReset+" %s\n", "Database URI", s.DbUri)

	if s.DbUser != "" {
		fmt.Printf(ColorGreen+"  %-18s:"+ColorReset+" %s\n", "Database User", s.DbUser)
	} else {
		fmt.Printf(ColorYellow+"  %-18s:"+ColorReset+" %s\n", "Database User", "[not set]")
	}

	if s.DbPassword != "" {
		fmt.Printf(ColorRed+"  %-18s:"+ColorReset+" %s\n", "Database Password", "[hidden]") // Hide for security
	} else {
		fmt.Printf(ColorYellow+"  %-18s:"+ColorReset+" %s\n", "Database Password", "[not set]")
	}

	fmt.Printf(ColorGreen+"  %-18s:"+ColorReset+" %t\n", "Debug Mode", s.Debug)
	fmt.Println(ColorCyan + "########################################" + ColorReset)
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
	buffer.WriteString(strconv.Itoa(int(s.DbPort)))

	return buffer.String()
}

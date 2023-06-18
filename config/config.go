package config

import (
	"bufio"
	"log"
	"os"
	"fmt"
	"strings"
	"sync"
)

// Config Public Interface ---------------------------------------------------------------
type configID int

const (
	HostPort configID = iota
	DebugMode
	ReservedParameter
)

func IsDebugModeEnabled() bool {
	return debugMode
}

func Get(id configID) string {
	if configMap == nil {
		initConfig()
	}
	variable := configMap[id]
	if variable == "" {
		const errFmt = "Environment variable with index:%d is empty or it doesn't exists"
		errMsg := fmt.Sprintf(errFmt, id)
		log.Fatal(errMsg)
	}
	return configMap[id]
}

// Config Private Implementation ---------------------------------------------------------
const (
	hostPortName = "PORT_FAKE"
	debugModeName = "DEBUG_MODE"
	reservedParameterName = "RESERVED_PARAMETER"
)

var (
	debugMode bool
	configOnce sync.Once
	configMap  map[configID]string
)

func init() {
	initEnv()
	defer initConfig()
	debugMode = strings.ToLower(strings.TrimSpace(Get(DebugMode))) == "true"
}

func initConfig() {
	err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig() error {
	configOnce.Do(func() {
		configMap = make(map[configID]string)
		configMap[HostPort] = getEnv(hostPortName)
		configMap[DebugMode] = getEnv(debugModeName)
		configMap[ReservedParameter] = getEnv(reservedParameterName)
	})
	return nil
}

// Env Private Implementation ------------------------------------------------------------
const (
	configFileName        = ".env"
	separatorOfConfigFile = "="
)

var (
	envOnce sync.Once
	envMap  map[string]string
)

func initEnv() {
	err := loadEnvFile()
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(envName string) string {
	if envMap == nil {
		initEnv()
	}
	return envMap[envName]
}

func loadEnvFile() error {
	envOnce.Do(func() {
		envMap = make(map[string]string)
		envFile, err := os.Open(configFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer envFile.Close()

		scanner := bufio.NewScanner(envFile)
		for scanner.Scan() {
			line := scanner.Text()
			keyVal := strings.Split(line, separatorOfConfigFile)
			if len(keyVal) != 2 {
				continue
			}
			envMap[keyVal[0]] = keyVal[1]
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	})
	return nil
}

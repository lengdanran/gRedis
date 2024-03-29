// Package config lengdanran 2024/3/27 17:04
// this package is used for read gredis configuration
package config

import (
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
	"log/slog"
	"os"
)

const CfgFilePath string = "config.yaml"

// ServerConfig holds global config properties
var ServerConfig *GRedisServerConfig

type GRedisServerConfig struct {
	Addr  string `yaml:"addr"`
	Port  int    `yaml:"port"`
	RunId string `yaml:"runId"`
}

func init() {
	ServerConfig = &GRedisServerConfig{
		Addr:  "0.0.0.0",
		Port:  6379,
		RunId: uuid.New().String(),
	}
	// Load configuration file
	// Position: body of function `init`
	// Path: config/config.go
	loadConfig(CfgFilePath, ServerConfig)
}

func loadConfig(configPath string, config *GRedisServerConfig) {
	configYaml, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("Fatal error config file: " + err.Error())
	}
	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		slog.Error("Failed to unmarshal YAML: " + err.Error())
	}
}

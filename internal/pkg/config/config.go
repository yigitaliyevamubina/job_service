package config

import (
	"os"
	"strings"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string

	Context struct {
		Timeout string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SslMode  string
	}

	OTLPCollector struct {
		Host string
		Port string
	}

	Kafka struct {
		Address []string
		Topic   struct {
			JobTopic string
		}
	}
}

func New() *Config {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "job_service")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.RPCPort = getEnv("RPC_PORT", ":9090")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "localhost")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "mubina2007")
	config.DB.SslMode = getEnv("POSTGRES_SSLMODE", "disable")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "companydb")

	config.OTLPCollector.Host = getEnv("OTLP_COLLECTOR_HOST", "0.0.0.0")
	config.OTLPCollector.Port = getEnv("OTLP_COLLECTOR_PORT", ":4317")

	// kafka configuration
	config.Kafka.Address = strings.Split(getEnv("KAFKA_ADDRESS", "localhost:9092"), ",")
	config.Kafka.Topic.JobTopic = getEnv("KAFKA_TOPIC_JOB_SERVICE", "job.service.create")

	return &config
}

func getEnv(key string, defaultVaule string) string {
	_, exists := os.LookupEnv(key)
	if exists {
		return defaultVaule
	}
	return defaultVaule
}

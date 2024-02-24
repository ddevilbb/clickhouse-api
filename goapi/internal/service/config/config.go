package config

import (
	"clickhouse-api/internal/repository/base/clickhouse"
	"clickhouse-api/internal/repository/base/sqlite"
	"clickhouse-api/pkg/elk_writer"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Application struct {
	Environment string
	Host        string
	Port        int
}

type QueueConfig struct {
	Path              string
	Extension         string
	CheckFileInterval int
}

type Config struct {
	Application *Application
	QueueConfig *QueueConfig
	Clickhouse  *clickhouse.Config
	Sqlite      *sqlite.Config
	ElkWriter   *elk_writer.Config
}

func New() *Config {
	exec, err := os.Executable()
	if err != nil {
		log.Fatalln(fmt.Errorf("os.Executable: %+v", err))
	}
	path := filepath.Dir(exec)
	if err := godotenv.Load(path + "/.env"); err != nil {
		log.Fatalf("No '%s/.env' file", path)
	}

	return &Config{
		Application: &Application{
			Environment: getEnvAsString("ENVIRONMENT", "dev"),
			Host:        getEnvAsString("APP_HOST", "0.0.0.0"),
			Port:        getEnvAsInt("APP_PORT", 8080),
		},
		QueueConfig: &QueueConfig{
			Path:              getEnvAsString("QUEUE_PATH", ""),
			Extension:         getEnvAsString("QUEUE_FILE_EXTENSION", ".tmp"),
			CheckFileInterval: getEnvAsInt("QUEUE_CHECK_FILE_INTERVAL", 5),
		},
		Clickhouse: &clickhouse.Config{
			ConnectionUrl: getEnvAsString("CLICKHOUSE_URL", ""),
		},
		Sqlite: &sqlite.Config{
			ConnectionFileSrc: getEnvAsString("SQLITE_FILE_SRC", ""),
		},
		ElkWriter: &elk_writer.Config{
			ConnectionNetwork: getEnvAsString("ELK_CONNECTION_NETWORK", ""),
			ConnectionUrl:     getEnvAsString("ELK_CONNECTION_URL", ""),
		},
	}
}

func getEnvAsString(name, defaultValue string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueString := getEnvAsString(name, "")
	if value, err := strconv.Atoi(valueString); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(name string, defaultValue bool) bool {
	valueString := getEnvAsString(name, "")
	if value, err := strconv.ParseBool(valueString); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(name string, defaultValue []string, sep string) []string {
	valueString := getEnvAsString(name, "")
	if valueString == "" {
		return defaultValue
	}
	value := strings.Split(valueString, sep)
	return value
}

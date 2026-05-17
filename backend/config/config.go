package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	OpenAI   OpenAIConfig
	AWS      AWSConfig
	Kafka    KafkaConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port        int
	Host        string
	Environment string
	Timeout     time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxConnections  int
	MinConnections  int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
	PoolSize int
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	RefreshTokenWindow time.Duration
}

type OpenAIConfig struct {
	APIKey string
	Model  string
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3Bucket        string
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

type LoggerConfig struct {
	Level      string
	Format     string
	EnableJSON bool
}

// Load loads configuration from environment variables with defaults
func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:        getEnvInt("SERVER_PORT", 8000),
			Host:        getEnv("SERVER_HOST", "0.0.0.0"),
			Environment: getEnv("SERVER_ENV", "development"),
			Timeout:     getDuration("SERVER_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			URL:             mustGetEnv("DATABASE_URL"),
			MaxConnections:  getEnvInt("DB_MAX_CONNECTIONS", 25),
			MinConnections:  getEnvInt("DB_MIN_CONNECTIONS", 5),
			ConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getDuration("DB_CONN_MAX_IDLE_TIME", 2*time.Minute),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			PoolSize: getEnvInt("REDIS_POOL_SIZE", 10),
		},
		JWT: JWTConfig{
			Secret:             mustGetEnv("JWT_SECRET"),
			AccessTokenExpiry:  getDuration("JWT_EXPIRATION", 15*time.Minute),
			RefreshTokenExpiry: getDuration("JWT_REFRESH_EXPIRATION", 7*24*time.Hour),
			RefreshTokenWindow: getDuration("JWT_REFRESH_WINDOW", 24*time.Hour),
		},
		OpenAI: OpenAIConfig{
			APIKey: getEnv("OPENAI_API_KEY", ""),
			Model:  getEnv("OPENAI_MODEL", "gpt-4-turbo"),
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			S3Bucket:        getEnv("S3_BUCKET", "teamart-uploads"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			GroupID: getEnv("KAFKA_GROUP_ID", "teamart-api"),
		},
		Logger: LoggerConfig{
			Level:      getEnv("LOG_LEVEL", "debug"),
			Format:     getEnv("LOG_FORMAT", "json"),
			EnableJSON: getBool("LOG_JSON", true),
		},
	}, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("required environment variable not set: %s", key))
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("invalid integer value for %s: %s", key, value))
	}
	return intVal
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}

	// Handle days (e.g., "7d", "1d")
	if len(value) > 1 && value[len(value)-1] == 'd' {
		days, err := strconv.Atoi(value[:len(value)-1])
		if err != nil {
			panic(fmt.Sprintf("invalid duration value for %s: %s", key, value))
		}
		return time.Duration(days) * 24 * time.Hour
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		panic(fmt.Sprintf("invalid duration value for %s: %s", key, value))
	}
	return duration
}

func getBool(key string, defaultValue bool) bool {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("invalid boolean value for %s: %s", key, value))
	}
	return boolVal
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

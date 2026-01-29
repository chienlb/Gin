package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Kafka     KafkaConfig
	Storage   StorageConfig
	Worker    WorkerConfig
	Discovery DiscoveryConfig
	Logger    LoggerConfig
	App       AppConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Environment  string // development, staging, production
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ReadReplicas    []string // Read replica URLs
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers       []string
	ConsumerGroup string
	Topics        KafkaTopics
}

type KafkaTopics struct {
	UserEvents string
}

type StorageConfig struct {
	Type            string // s3 or minio
	Endpoint        string // For MinIO
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	UsePathStyle    bool
	CDNDomain       string // CDN domain for public URLs
}

type WorkerConfig struct {
	Enabled     bool
	WorkerCount int
	QueueSize   int
}

type DiscoveryConfig struct {
	Enabled     bool
	Type        string // consul or kubernetes
	ConsulAddr  string
	ServiceName string
	ServiceID   string
	ServicePort int
	Tags        []string
}

type LoggerConfig struct {
	Level string // debug, info, warn, error
}

type AppConfig struct {
	Name    string
	Version string
}

func Load() *Config {
	env := getEnv("ENVIRONMENT", "development")

	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			Environment:  env,
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "gin_db"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ReadReplicas:    getEnvSlice("DB_READ_REPLICAS", []string{}),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Kafka: KafkaConfig{
			Brokers:       getEnvSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "gin-api-group"),
			Topics: KafkaTopics{
				UserEvents: getEnv("KAFKA_TOPIC_USER_EVENTS", "user-events"),
			},
		},
		Storage: StorageConfig{
			Type:            getEnv("STORAGE_TYPE", "s3"),
			Endpoint:        getEnv("STORAGE_ENDPOINT", ""),
			Region:          getEnv("STORAGE_REGION", "us-east-1"),
			AccessKeyID:     getEnv("STORAGE_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("STORAGE_SECRET_ACCESS_KEY", ""),
			Bucket:          getEnv("STORAGE_BUCKET", "gin-demo-uploads"),
			UsePathStyle:    getEnvBool("STORAGE_USE_PATH_STYLE", false),
			CDNDomain:       getEnv("STORAGE_CDN_DOMAIN", ""),
		},
		Worker: WorkerConfig{
			Enabled:     getEnvBool("WORKER_ENABLED", false),
			WorkerCount: getEnvInt("WORKER_COUNT", 5),
			QueueSize:   getEnvInt("WORKER_QUEUE_SIZE", 100),
		},
		Discovery: DiscoveryConfig{
			Enabled:     getEnvBool("DISCOVERY_ENABLED", false),
			Type:        getEnv("DISCOVERY_TYPE", "kubernetes"),
			ConsulAddr:  getEnv("CONSUL_ADDR", "localhost:8500"),
			ServiceName: getEnv("SERVICE_NAME", "gin-demo-api"),
			ServiceID:   getEnv("SERVICE_ID", "gin-demo-api-1"),
			ServicePort: getEnvInt("SERVICE_PORT", 8080),
			Tags:        getEnvSlice("SERVICE_TAGS", []string{"api", "v1"}),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		App: AppConfig{
			Name:    "Gin Demo API",
			Version: "1.0.0",
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Storage.Type != "" && c.Storage.Type != "s3" && c.Storage.Type != "minio" {
		return fmt.Errorf("invalid storage type: %s (must be s3 or minio)", c.Storage.Type)
	}
	return nil
}

func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

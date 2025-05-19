package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents application configuration
type Config struct {
	AppName      string
	AppPort      int
	DBDriver     string
	DBHost       string
	DBPort       int
	DBUser       string
	DBPassword   string
	DBName       string
	DBUrl        string
	JWTSecretKey string
	JWTExpiry    time.Duration
}

// New returns application configuration
func New() *Config {
	// Parse DATABASE_URL if available
	dbUrl := getEnv("DATABASE_URL", "")
	dbDriver, dbHost, dbPort, dbUser, dbPassword, dbName := parseDbUrl(dbUrl)

	return &Config{
		AppName:      getEnv("APP_NAME", "dot-be-go"),
		AppPort:      getEnvAsInt("APP_PORT", 8080),
		DBDriver:     dbDriver,
		DBHost:       dbHost,
		DBPort:       dbPort,
		DBUser:       dbUser,
		DBPassword:   dbPassword,
		DBName:       dbName,
		DBUrl:        dbUrl,
		JWTSecretKey: getEnv("JWT_SECRET", "mySecretKey"),
		JWTExpiry:    time.Duration(getEnvAsInt("JWT_EXPIRY", 24)) * time.Hour,
	}
}

// DBConnectionString returns database connection string
func (c *Config) DBConnectionString() string {
	// If we have a full DB URL, use it directly
	if c.DBUrl != "" {
		return c.DBUrl
	}

	// Otherwise, construct the connection string from individual parts
	if c.DBDriver == "postgres" {
		return "host=" + c.DBHost + " port=" + strconv.Itoa(c.DBPort) + " user=" + c.DBUser + " password=" + c.DBPassword + " dbname=" + c.DBName + " sslmode=disable"
	}
	// Default to MySQL
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + strconv.Itoa(c.DBPort) + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

// Helper function to parse database URL
func parseDbUrl(dbUrl string) (driver, host string, port int, user, password, dbName string) {
	// Default values in case URL parsing fails
	driver = getEnv("DB_DRIVER", "postgres")
	host = getEnv("DB_HOST", "127.0.0.1")
	port = getEnvAsInt("DB_PORT", 5432)
	user = getEnv("DB_USER", "postgres")
	password = getEnv("DB_PASSWORD", "123")
	dbName = getEnv("DB_NAME", "bookdb")

	// Return defaults if URL is empty
	if dbUrl == "" {
		return
	}

	// Parse URL
	// Format: postgresql://postgres:123@127.0.0.1:5432/bookdb
	if strings.HasPrefix(dbUrl, "postgresql://") {
		driver = "postgres"
		// Remove protocol
		urlWithoutProtocol := strings.TrimPrefix(dbUrl, "postgresql://")

		// Split user/pass from host/port/dbname
		parts := strings.Split(urlWithoutProtocol, "@")
		if len(parts) == 2 {
			// Handle user:pass part
			userPass := strings.Split(parts[0], ":")
			if len(userPass) >= 1 {
				user = userPass[0]
			}
			if len(userPass) >= 2 {
				password = userPass[1]
			}

			// Handle host:port/dbname part
			hostPortDb := strings.Split(parts[1], "/")
			if len(hostPortDb) >= 1 {
				hostPort := strings.Split(hostPortDb[0], ":")
				if len(hostPort) >= 1 {
					host = hostPort[0]
				}
				if len(hostPort) >= 2 {
					portStr := hostPort[1]
					if p, err := strconv.Atoi(portStr); err == nil {
						port = p
					}
				}
			}
			if len(hostPortDb) >= 2 {
				dbName = hostPortDb[1]
			}
		}
	}

	return
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

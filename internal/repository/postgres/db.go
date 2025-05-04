package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"matching-engine/internal/errors"
)

// Database represents a database connection
type Database struct {
	DB *gorm.DB
}

// Config holds all database configuration
type Config struct {
	Host        string
	Port        int
	DBName      string
	Username    string
	Password    string
	SSLMode     string
	MinConns    int
	MaxConns    int
	MaxIdleTime time.Duration
	MaxLifetime time.Duration
	LogLevel    logger.LogLevel
}

// NewDatabase creates a database connection using environment variables
func NewDatabase(ctx context.Context) (*Database, error) {
	if err := loadEnv(); err != nil {
		return nil, errors.Wrap(err, "failed to load environment variables")
	}

	config, err := loadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	connString := buildConnString(config)
	return NewDatabaseFromConnString(ctx, connString, config)
}

// LoadEnv loads environment variables from .env file
func loadEnv() error {
	// Try to load from .env file, but don't fail if file is missing
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
	return nil
}

// loadConfigFromEnv loads database configuration from environment variables
func loadConfigFromEnv() (*Config, error) {
    config := &Config{
        Host:        getEnvOrDefault("DB_HOST", ""),
        Port:        getEnvAsIntOrDefault("DB_PORT", 5432),
        DBName:      getEnvOrDefault("DB_NAME", ""),
        Username:    getEnvOrDefault("DB_USER", ""),
        Password:    getEnvOrDefault("DB_PASSWORD", ""),
        SSLMode:     getEnvOrDefault("DB_SSLMODE", "require"),
        MaxConns:    getEnvAsIntOrDefault("DB_MAX_CONNS", 10),
        MinConns:    getEnvAsIntOrDefault("DB_MIN_CONNS", 2),
        MaxIdleTime: 5 * time.Minute, // Default value
        MaxLifetime: 1 * time.Hour,   // Default value
        LogLevel:    parseLogLevel(getEnvOrDefault("DB_LOG_LEVEL", "silent")),
    }

    // Validate required fields
    var missingFields []string
    
    if config.Host == "" {
        missingFields = append(missingFields, "DB_HOST")
    }
    if config.DBName == "" {
        missingFields = append(missingFields, "DB_NAME")
    }
    if config.Username == "" {
        missingFields = append(missingFields, "DB_USER")
    }
    if config.Password == "" {
        missingFields = append(missingFields, "DB_PASSWORD")
    }
    
    if len(missingFields) > 0 {
        return nil, errors.InvalidInput(fmt.Sprintf("missing required database configuration: %v", missingFields))
    }

    return config, nil
}

// getEnvOrDefault returns environment variable or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault returns environment variable as int or default
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// parseLogLevel converts string to GORM log level
func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "info":
		return logger.Info
	case "warn", "warning":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Silent
	}
}

// buildConnString creates a connection string for PostgreSQL
func buildConnString(config *Config) string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		config.Host, config.Port, config.DBName, config.Username, config.Password, config.SSLMode,
	)
}

// NewDatabaseFromConnString creates a new database connection with the provided connection string
func NewDatabaseFromConnString(ctx context.Context, connString string, config *Config) (*Database, error) {
	// Configure GORM logger
	gormLogger := createGormLogger(config.LogLevel)

	// Open database connection
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, errors.DatabaseError("connect", err)
	}

	// Configure connection pool
	if err := configureConnectionPool(db, config); err != nil {
		return nil, err
	}

	// Verify connection
	if err := verifyConnection(ctx, db, config); err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

// createGormLogger configures a GORM logger
func createGormLogger(logLevel logger.LogLevel) logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)
}

// configureConnectionPool sets up database connection pool parameters
func configureConnectionPool(db *gorm.DB, config *Config) error {
	sqlDB, err := db.DB()
	if err != nil {
		return errors.DatabaseError("get_sql_db", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxConns)
	sqlDB.SetMaxIdleConns(config.MinConns)
	sqlDB.SetConnMaxIdleTime(config.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	return nil
}

// verifyConnection checks that the database connection is working
func verifyConnection(ctx context.Context, db *gorm.DB, config *Config) error {
	sqlDB, err := db.DB()
	if err != nil {
		return errors.DatabaseError("get_sql_db_for_ping", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return errors.DatabaseError("ping", err)
	}

	log.Printf("Successfully connected to PostgreSQL database at %s/%s", config.Host, config.DBName)
	return nil
}

// Close closes the database connection pool
func (d *Database) Close() error {
	if d.DB == nil {
		return nil
	}

	sqlDB, err := d.DB.DB()
	if err != nil {
		return errors.DatabaseError("get_sql_db_for_close", err)
	}

	if err := sqlDB.Close(); err != nil {
		return errors.DatabaseError("close", err)
	}

	log.Println("Database connection pool closed")
	return nil
}

// Package config handles the initialization and management of configuration settings for the GophKeeper application.
// It includes functionality to parse command-line arguments
// and environment variables to set up the application configuration.
package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

// TOKENEXP defines the default expiration duration for tokens.
const TOKENEXP = time.Hour * 8766

// Config struct holds the configuration settings for the application.
type Config struct {
	// URLServer specifies the address for the execution server.
	URLServer string

	// LoggerLevel defines the level of logging.
	LoggerLevel string

	// DatabaseDsn is the connection string for the database.
	DatabaseDsn string

	// Secret is a string of secret to create tokens
	Secret string

	// TokenExp specifies the period after which a token expires.
	TokenExp time.Duration

	// IsSecure indicates if the gRPC server should be secure (TLS).
	IsSecure bool

	// EndpointMinio is the address for Minio storage.
	EndpointMinio string

	// AccessKeyIDMinio is the username for Minio storage.
	AccessKeyIDMinio string

	// SecretAccessKeyMinio is the password for Minio storage.
	SecretAccessKeyMinio string

	// BucketNameMinio is the default bucket for Minio storage.
	BucketNameMinio string

	// IsUseSSLMinio indicates if SSL mode is enabled for Minio storage.
	IsUseSSLMinio bool
}

// Flag variables for command-line arguments.
//
//nolint:gochecknoglobals
var (
	URLServer   *string
	logLevel    *string
	databaseDsn *string

	secret *string

	isSecure *bool

	endpointMinio *string

	accessKeyIDMinio *string

	secretAccessKeyMinio *string

	bucketNameMinio *string

	isUseSSLMinio *bool
)

//nolint:gochecknoinits
func init() {
	// declare flag set for subcommand
	URLServer = flag.String("a", "localhost:8080",
		"Enter address exec http server as ip_address:port. Or use SERVER_ADDRESS env")
	logLevel = flag.String("l",
		"Info",
		"log level: Debug, Info, Warn, Error and etc... Or use LOG_LEVEL env")
	databaseDsn = flag.String("d",
		"postgres://gophkeeperdb:gophkeeperdbpwd@localhost:5432/gophkeeperdb?sslmode=disable",
		`set path for database... Or use DATABASE_DSN env. 
	Example postgres://username:password@hostname:port/databasename?sslmode=disable`)

	secret = flag.String(
		"s", "supersecret",
		"Enter secret. Or use SECRET env")

	isSecure = flag.Bool("secure", true, "enable secure grpc? Set true/false... Or use ISSECURE env")

	// args for minio
	endpointMinio = flag.String(
		"endpointMinio", "localhost:9000",
		"Enter endpoint for Minio. Or use ENDPOINTMINIO env")

	accessKeyIDMinio = flag.String(
		"accessKeyIDMinio", "masoud",
		"Enter accessKeyID for Minio. Or use ACCESSKEYIDMINIO env")

	secretAccessKeyMinio = flag.String(
		"secretAccessKeyMinio", "Strong#Pass#2022",
		"Enter secretAccessKey for Minio. Or use SECRETACCESSKEYMINIO env")

	bucketNameMinio = flag.String(
		"bucketNameMinio", "secrets",
		"Enter bucketNameMinio for Minio. Or use BUCKETNAMEMINIO env")

	isUseSSLMinio = flag.Bool("isUseSSLMinio", false,
		"enable ssl for Minio? Set true/false... Or use ISUSESSLMINIO env")
}

// InitConfig parses the command-line flags and environment variables to initialize the application configuration.
// It returns a configured Config struct or an error if the configuration setup fails.
//
//nolint:cyclop
func InitConfig() (*Config, error) {
	flag.Parse()

	if envSrvAddr := os.Getenv("SERVER_ADDRESS"); envSrvAddr != "" {
		*URLServer = envSrvAddr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		*logLevel = envLogLevel
	}

	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		*databaseDsn = envDatabaseDsn
	}

	if envSecret := os.Getenv("SECRET"); envSecret != "" {
		*secret = envSecret
	}

	if envIsSecure := os.Getenv("ISSECURE"); envIsSecure != "" {
		parsedValue, err := strconv.ParseBool(envIsSecure)
		if err != nil {
			// Обработка ошибки, если значение не может быть преобразовано в bool
			*isSecure = true
		} else {
			*isSecure = parsedValue
		}
	}

	// minio env
	if envEndpointMinio := os.Getenv("ENDPOINTMINIO"); envEndpointMinio != "" {
		*endpointMinio = envEndpointMinio
	}
	if envAccessKeyIDMinio := os.Getenv("ACCESSKEYIDMINIO"); envAccessKeyIDMinio != "" {
		*accessKeyIDMinio = envAccessKeyIDMinio
	}
	if envSecretAccessKeyMinio := os.Getenv("SECRETACCESSKEYMINIO"); envSecretAccessKeyMinio != "" {
		*secretAccessKeyMinio = envSecretAccessKeyMinio
	}
	if envBucketNameMinio := os.Getenv("BUCKETNAMEMINIO"); envBucketNameMinio != "" {
		*bucketNameMinio = envBucketNameMinio
	}

	if envIsSecure := os.Getenv("ISUSESSLMINIO"); envIsSecure != "" {
		parsedValue, err := strconv.ParseBool(envIsSecure)
		if err != nil {
			// Обработка ошибки, если значение не может быть преобразовано в bool
			*isUseSSLMinio = true
		} else {
			*isUseSSLMinio = parsedValue
		}
	}

	cfg := &Config{
		URLServer:            *URLServer,
		LoggerLevel:          *logLevel,
		DatabaseDsn:          *databaseDsn,
		Secret:               *secret,
		TokenExp:             TOKENEXP,
		IsSecure:             *isSecure,
		EndpointMinio:        *endpointMinio,
		AccessKeyIDMinio:     *accessKeyIDMinio,
		SecretAccessKeyMinio: *secretAccessKeyMinio,
		BucketNameMinio:      *bucketNameMinio,
		IsUseSSLMinio:        *isUseSSLMinio,
	}

	return cfg, nil
}

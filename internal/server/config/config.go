package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

const TOKENEXP = time.Hour * 8766

type Config struct {
	// it's address for exec server
	URLServer string

	// it's logger level
	LoggerLevel string

	// it's database conn string
	DatabaseDsn string

	// Secret is a string of secret to create tokens
	Secret string

	// period of expired token
	TokenExp time.Duration

	// IsSecure if grpc server should be secure (tls)
	IsSecure bool

	// EndpointMinio is an address minio storage
	EndpointMinio string

	// AccessKeyIDMinio is an username for minio storage
	AccessKeyIDMinio string

	// SecretAccessKeyMinio is an password for minio storage
	SecretAccessKeyMinio string

	// BucketNameMinio is an default bucket for minio storage
	BucketNameMinio string

	// isUseSSL is an ssl mode for minio storage
	IsUseSSLMinio bool
}

func InitConfig() (*Config, error) {

	// declare flag set for subcommand
	URLServer := flag.String("a", "localhost:8080", "Enter address exec http server as ip_address:port. Or use SERVER_ADDRESS env")
	logLevel := flag.String("l", "Info", "log level: Debug, Info, Warn, Error and etc... Or use LOG_LEVEL env")
	databaseDsn := flag.String("d",
		"postgres://secretkeeper:secret@localhost:5433/secret_keeper_server?sslmode=disable",
		`set path for database... Or use DATABASE_DSN env. 
		Example postgres://username:password@hostname:port/databasename?sslmode=disable`)

	secret := flag.String(
		"s", "supersecret",
		"Enter secret. Or use SECRET env")

	isSecure := flag.Bool("secure", true, "enable secure grpc? Set true/false... Or use ISSECURE env")

	// args for minio
	endpointMinio := flag.String(
		"endpointMinio", "localhost:9000",
		"Enter endpoint for Minio. Or use ENDPOINTMINIO env")

	accessKeyIDMinio := flag.String(
		"accessKeyIDMinio", "masoud",
		"Enter accessKeyID for Minio. Or use ACCESSKEYIDMINIO env")

	secretAccessKeyMinio := flag.String(
		"secretAccessKeyMinio", "Strong#Pass#2022",
		"Enter secretAccessKey for Minio. Or use SECRETACCESSKEYMINIO env")

	bucketNameMinio := flag.String(
		"bucketNameMinio", "secrets",
		"Enter bucketNameMinio for Minio. Or use BUCKETNAMEMINIO env")

	isUseSSLMinio := flag.Bool("isUseSSLMinio", false,
		"enable ssl for Minio? Set true/false... Or use ISUSESSLMINIO env")

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

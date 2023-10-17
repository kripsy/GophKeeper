package config

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	// it's address for exec server
	URLServer string

	// it's logger level
	LoggerLevel string

	// it's database conn string
	DatabaseDsn string

	// PrivateKey is a string of secret to create tokens
	PrivateKey string

	// PublicKey is using for validate tokens
	PublicKey string

	TokenExp time.Duration
}

var defaultPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQDe95shZARG7OCf
RP/vS1r1AqO6+d5O3zhieIz6Mq1Rc6Njr0Ks8iZ/phWKr6HlyOqAoqStokPvOpEt
FAf8yWN5nStIjrD2Q8u1ST0GNAzLNGBo+W3PgC/FoIUADrQBSEfp2klubU8QAVtl
HA5XPAjx9pdGNsPJJZvUzyCfDdrDVvbLtpIybw4bHn7DeTGh7ovTo/5vah5kOjLy
dhx/UR9bbrf6rtGwGyS9vAVsIoA9Peag5X5XceLUKZz5IK8dwKm9+fhSkvmfF10c
k3s3LzgXZPS5qBWsaUXsEF+Ju2RP1ElAKrp6Gv3JEDYJJnDUWjVF1SfpIYHoIte8
2zm9CR9NAgMBAAECggEBAI3+EYUKNM8WO1YykurJintN2wdP6QtBjJ7pNp5/d3DP
u9XX3xZUf7/6/Oz9PJUhhnW1HjqVg73uBlY2039goUDpno7ukDPEqQ4iPgKdUyh1
ipBPiGcEs2ef+hM3SdsnNOTwZqM0aY0/z/xsCZX0XZ359Ax7A+QtVzgHUDb6k76h
irJaTMmrtTzmTCnm72tOGhX91QusLcefffZToPjEPRlNazxeH/wdkNbYMc1GmNVH
D/Sq3esk22t/cpeImKtv7LhShd0NCbPtM8lIJY++cBOmsM5UaQRaXrc5NSV1o1sZ
Zjr+xJN0//p0TJAB06qhb7XsSCy4zHPXZC0cNleWiwECgYEA+GFvnP3ux6QL32Jh
zrf6cz1P277BH/NWXKZpuqMnUiBfsHETjOASIPGytc6RL1jbTn0thkg3tt1SG1ja
9k5pHHCtQxY5MiRpXNT52Bbu1Ko7i+e9rG7mfaws2ItWRtQSYArWj8pTMj0NpzD+
kesgDBc67U/Gl4skLSm2WbNS9TUCgYEA5c6Ya31E9LXP4WCuBjTNDVMzb5+G7iBV
TXQH69GhWaLTuYE1B1BR87lVj/ZVXRY9wmpQRHlmABfZPEvA77HladKUUA4oX35E
vzjBvBp8WnuYk1VLR6vl+9AO1GLGnOoi9sYazl8l8KW+5WDHCEyQN0TnrG+jrfpw
HelbtvFXvLkCgYEA2RNnFcEEuCySR8hXDPDUDXVvXvEHHmJwfwbd7sT675blqnIZ
EQ0gKvSyKJ0BXGz/NkjGyc5CCyrAwK/Wpl9/E+ESPEim8kDKaNymAwp/7xNceXiu
144RGZKpmxOj8sET0iaGwSKltYmQbieuxV7GImsHEDKhsP5lPqdu/FRyU2UCgYEA
rEodf8jlH8oHVmNTVRfU+756+57QXEslaPIq1iPOIhOvRI6YISmYp281tL7r9OQt
3Uozb4LMdBltJoVs2se2xYW45+QVZLKX+/0jUlFRFc0/8IWr8MnxnL65v4VmflIT
cIvJoRs4qJi66+GIlrJAFQ+12VPBlTgDQomn1xpNuxECgYEAvQRbKiKSRRSoph7n
Nm4DHwrquwarjwcynPPiufl3vZiJpgd6Zn6/cRpS4J6JUqSVTO0O1c9G6/EM2zsk
tRgiBEdw38SbenkxGkTZt4kNV95qzowO2Svd+5l2nM3mfBMAAHcToqtvq0WQNNZ1
67yGLJcuzkkgH15moAdLrN2qNQU=
-----END PRIVATE KEY-----`

var defaultPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3vebIWQERuzgn0T/70ta
9QKjuvneTt84YniM+jKtUXOjY69CrPImf6YViq+h5cjqgKKkraJD7zqRLRQH/Mlj
eZ0rSI6w9kPLtUk9BjQMyzRgaPltz4AvxaCFAA60AUhH6dpJbm1PEAFbZRwOVzwI
8faXRjbDySWb1M8gnw3aw1b2y7aSMm8OGx5+w3kxoe6L06P+b2oeZDoy8nYcf1Ef
W263+q7RsBskvbwFbCKAPT3moOV+V3Hi1Cmc+SCvHcCpvfn4UpL5nxddHJN7Ny84
F2T0uagVrGlF7BBfibtkT9RJQCq6ehr9yRA2CSZw1Fo1RdUn6SGB6CLXvNs5vQkf
TQIDAQAB
-----END PUBLIC KEY-----`

func InitConfig() *Config {

	// declare flag set for subcommand
	URLServer := flag.String("a", "localhost:8080", "Enter address exec http server as ip_address:port. Or use SERVER_ADDRESS env")
	logLevel := flag.String("l", "Info", "log level: Debug, Info, Warn, Error and etc... Or use LOG_LEVEL env")
	databaseDsn := flag.String("d",
		"postgres://gophkeeperdb:gophkeeperdbpwd@localhost:5432/gophkeeperdb?sslmode=disable",
		`set path for database... Or use DATABASE_DSN env. 
		Example postgres://username:password@hostname:port/databasename?sslmode=disable`)

	privateKey := flag.String(
		"privateKey", defaultPrivateKey,
		"Enter private key. Or use PRIVATE_KEY env")

	publicKey := flag.String(
		"publicKey", defaultPublicKey,
		"Enter public key. Or use PUBLIC_KEY env")

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

	if envPrivateKey := os.Getenv("PRIVATE_KEY"); envPrivateKey != "" {
		*privateKey = envPrivateKey
	}

	if envPublicKey := os.Getenv("PUBLIC_KEY"); envPublicKey != "" {
		*publicKey = envPublicKey
	}

	cfg := &Config{
		URLServer:   *URLServer,
		LoggerLevel: *logLevel,
		DatabaseDsn: *databaseDsn,
		PrivateKey:  *privateKey,
		PublicKey:   *publicKey,
		TokenExp:    time.Hour * 8766,
	}
	return cfg
}

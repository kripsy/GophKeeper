// Package config manages the configuration settings for the GophKeeper application.
// This part includes the definition and parsing of command-line flags.
package config

import "flag"

// Flags defines the structure to hold command-line arguments.
// It includes options for specifying custom paths for storage, upload, server address,
// and the path to the configuration file.
type Flags struct {
	StoragePath   string
	UploadPath    string
	ServerAddress string
	ConfigPath    string
}

// parseFlags parses command-line arguments into a Flags struct.
// It uses the flag package to define and parse command-line options, providing
// users with the flexibility to override default configuration settings.
func parseFlags() Flags {
	var f Flags
	flag.StringVar(&f.ConfigPath, "cfg-path", "",
		"Specify the path to the config.yaml if it is not in the same directory as the program")
	flag.StringVar(&f.StoragePath, "storage-path", "",
		"You can change the default file storage directory")
	flag.StringVar(&f.UploadPath, "upload-path", "",
		"You can change the default file upload directory")
	flag.StringVar(&f.ServerAddress, "server-addr", "",
		"Server address for data synchronization (be careful)")

	flag.Parse()

	return f
}

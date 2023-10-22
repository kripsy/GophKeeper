package config

import "flag"

type Flags struct {
	StoragePath   string
	UploadPath    string
	ServerAddress string
	ConfigPath    string
}

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

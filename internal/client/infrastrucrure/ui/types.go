package ui

const (
	SecretsKey      = "Secrets"
	AddSecretKey    = "Add Secret"
	DeleteSecretKey = "Delete Secret"
	UpdateSecretKey = "Update Secret"
	SyncSecrets     = "Sync Secrets"
	About           = "About"
	ExitKey         = "Exit"
)

var MenuTable = []string{SyncSecrets, SecretsKey, AddSecretKey, UpdateSecretKey, DeleteSecretKey, About, ExitKey}
var LocalMenuTable = []string{SecretsKey, AddSecretKey, UpdateSecretKey, DeleteSecretKey, About, ExitKey}

const (
	Data = "Data"
	Info = "Info"
)

var UpdateItems = []string{Data, Info}

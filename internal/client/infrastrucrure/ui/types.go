// Package ui defines user interface elements and messages for the GophKeeper application.
// It includes constants for various error messages displayed to the user.
//
//nolint:gochecknoglobals
package ui

const (
	// SecretsKey is the menu option for accessing secrets.
	SecretsKey = "Secrets"
	// AddSecretKey is the menu option for adding a new secret.
	AddSecretKey = "Add Secret"
	// DeleteSecretKey is the menu option for deleting a secret.
	DeleteSecretKey = "Delete Secret"
	// UpdateSecretKey is the menu option for updating an existing secret.
	UpdateSecretKey = "Update Secret"
	// SyncSecrets is the menu option for synchronizing secrets across devices.
	SyncSecrets = "Sync Secrets"
	// About is the menu option for displaying information about the application.
	About = "About"
	// ExitKey is the menu option for exiting the application.
	ExitKey = "Exit"
)

// MenuTable defines the set of options available in the main menu.
var MenuTable = []string{SyncSecrets, SecretsKey, AddSecretKey, UpdateSecretKey, DeleteSecretKey, About, ExitKey}

// LocalMenuTable defines the set of options available in the local-only mode menu.
var LocalMenuTable = []string{SecretsKey, AddSecretKey, UpdateSecretKey, DeleteSecretKey, About, ExitKey}

const (
	// Data represents the option to update the data of a secret.
	Data = "Data"
	// Info represents the option to update the metadata/info of a secret.
	Info = "Info"
)

// UpdateItems defines the options available for updating a secret.
var UpdateItems = []string{Data, Info}

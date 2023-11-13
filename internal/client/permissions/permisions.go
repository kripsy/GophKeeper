// Package permissions defines standard file and directory permissions for the GophKeeper application.
// It specifies the access permissions to be used when creating files and directories.
package permissions

import "os"

const (
	// FileMode sets the permissions for files created by the application.
	// The 0700 permission grants read, write, and execute rights to the file's owner only.
	FileMode os.FileMode = 0700
	// DirMode sets the permissions for directories created by the application.
	// The 0700 permission grants read, write, and execute rights to the directory's owner only.
	DirMode os.FileMode = 0700
)

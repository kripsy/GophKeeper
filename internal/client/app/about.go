// Package app provides core functionalities for the GophKeeper application,
// including build information and utility functions.
package app

import "fmt"

// aboutMsg is a constant template for displaying information about the
// GophKeeper application. It includes placeholders for the build version
// and build date.
const aboutMsg = `
「  GophKeeper  」
This is an Application with the ability to store your secrets locally,
as well as synchronize between your clients when registering through the server.
Build version: %s                          Build date: %s
`

// BuildInfo holds information about the build version and date.
// It is used to display the current version of the application.
type BuildInfo struct {
	BuildVersion string // BuildVersion is the current version of the application.
	BuildDate    string // BuildDate is the date when the current version was built.
}

// about returns a formatted string containing the application's build version
// and date, providing users with details about the version of GophKeeper they are using.
func about(info BuildInfo) string {
	return fmt.Sprintf(aboutMsg, info.BuildVersion, info.BuildDate)
}

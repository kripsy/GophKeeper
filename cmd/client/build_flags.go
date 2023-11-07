//nolint:gochecknoglobals
package main

import "github.com/kripsy/GophKeeper/internal/client/app"

var buildVersion = defaultValue
var buildDate = defaultValue

const defaultValue = "N/A"

func getBuildInfo() app.BuildInfo {
	return app.BuildInfo{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
	}
}

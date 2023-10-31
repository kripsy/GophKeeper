package app

import "fmt"

const aboutMsg = `
「  GophKeeper  」
This is an Application with the ability to store your secrets locally,
as well as synchronize between your clients when registering through the server.
Build version: %s                          Build date: %s
`

type BuildInfo struct {
	BuildVersion string
	BuildDate    string
}

func about(info BuildInfo) string {
	return fmt.Sprintf(aboutMsg, info.BuildVersion, info.BuildDate)
}

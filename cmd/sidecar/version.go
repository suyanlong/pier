package main

import (
	"fmt"

	"github.com/link33/sidecar"
	"github.com/urfave/cli"
)

var versionCMD = cli.Command{
	Name:  "version",
	Usage: "Show version about ap",
	Action: func(ctx *cli.Context) error {
		fmt.Print(getVersion(true))

		return nil
	},
}

func getVersion(all bool) string {
	version := fmt.Sprintf("Sidecar version: %s-%s-%s\n", sidecar.CurrentVersion, sidecar.CurrentBranch, sidecar.CurrentCommit)
	if all {
		version += fmt.Sprintf("App build date: %s\n", sidecar.BuildDate)
		version += fmt.Sprintf("System version: %s\n", sidecar.Platform)
		version += fmt.Sprintf("Golang version: %s\n", sidecar.GoVersion)
	}

	return version
}

package main

import (
	"fmt"

	"github.com/link33/sidercar"
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
	version := fmt.Sprintf("Sidercar version: %s-%s-%s\n", sidercar.CurrentVersion, sidercar.CurrentBranch, sidercar.CurrentCommit)
	if all {
		version += fmt.Sprintf("App build date: %s\n", sidercar.BuildDate)
		version += fmt.Sprintf("System version: %s\n", sidercar.Platform)
		version += fmt.Sprintf("Golang version: %s\n", sidercar.GoVersion)
	}

	return version
}

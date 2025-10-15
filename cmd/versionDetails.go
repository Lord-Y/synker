// Package cmd manage all commands required to run synker
package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

// VersionDetails command options
func VersionDetails() (z *cli.Command) {
	return &cli.Command{
		Name:  "version-details",
		Usage: "print version details of the cli",
		Action: func(context.Context, *cli.Command) error {
			fmt.Printf("Version: %s\n", Version)
			fmt.Printf("Revision: %s\n", revision)
			fmt.Printf("Build date: %s\n", buildDate)
			fmt.Printf("Go version: %s\n", goVersion)
			return nil
		},
	}
}

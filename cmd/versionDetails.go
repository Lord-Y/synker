// Package cmd manage all commands required to run synker
package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// VersionDetails command options
func VersionDetails(c *cli.Context) (z *cli.Command) {
	return &cli.Command{
		Name:  "version-details",
		Usage: "print version details of the cli",
		Action: func(c *cli.Context) error {
			fmt.Printf("Version: %s\n", Version)
			fmt.Printf("Revision: %s\n", revision)
			fmt.Printf("Build date: %s\n", buildDate)
			fmt.Printf("Go version: %s\n", goVersion)
			return nil
		},
	}
}

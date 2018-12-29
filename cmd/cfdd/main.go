package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/pirogoeth/cfdd/cmd"
)

var (
	Version   string
	BuildHash string

	commands []cli.Command = []cli.Command{
		cmd.CurrentCmd,
		cmd.ListRecordsCmd,
		cmd.UpdateCmd,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "cfdd"
	app.Usage = "Dynamic DNS Updater for Cloudflare"
	app.Version = fmt.Sprintf("%s (%s)", Version, BuildHash)
	app.HideHelp = true
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, debug",
			Usage: "Be chattier about things",
		},
		cli.StringFlag{
			Name:   "cf-email",
			Usage:  "Cloudflare account email address",
			EnvVar: "CF_EMAIL",
		},
		cli.StringFlag{
			Name:   "cf-token",
			Usage:  "Cloudflare API token",
			EnvVar: "CF_TOKEN",
		},
		cli.StringFlag{
			Name:   "zone",
			Usage:  "DNS Zone domain",
			EnvVar: "DOMAIN",
		},
		cli.StringFlag{
			Name:   "record-name, rname",
			Usage:  "DNS Record name (used as 'record-name'.'zone')",
			EnvVar: "RECORD_NAME",
		},
		cli.StringFlag{
			Name:   "interface, iface",
			Usage:  "Local interface to use for addressing",
			EnvVar: "INTERFACE",
			Value:  "eth0",
		},
		cli.BoolFlag{
			Name:   "filter-unroutable",
			Usage:  "Filters unroutable addresses from the local interface",
			EnvVar: "FILTER_UNROUTABLE",
		},
		cli.BoolFlag{
			Name:   "v4-only",
			Usage:  "Only update record with IPv4 addresses",
			EnvVar: "V4_ONLY",
		},
		cli.BoolFlag{
			Name:   "v6-only",
			Usage:  "Only update record with IPv6 addresses",
			EnvVar: "V6_ONLY",
		},
	}

	app.Before = func(ctx *cli.Context) error {
		verbose := ctx.Bool("verbose")
		if verbose {
			log.SetFormatter(&log.TextFormatter{})
			log.SetOutput(os.Stderr)
			log.SetLevel(log.DebugLevel)
			log.Debug("Verbose logging enabled")
		}

		if ctx.Bool("v4-only") && ctx.Bool("v6-only") {
			return cli.NewExitError("--v4-only and --v6-only may not be specified together", 127)
		}

		return nil
	}

	app.Commands = commands

	app.Run(os.Args)
}

package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/pirogoeth/cfdd/cmd"
)

var (
	Version   string
	BuildHash string

	commands []cli.Command = []cli.Command{
		cmd.CurrentCmd,
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
			Name: "cf-email",
			Usage: "Cloudflare account email address",
			EnvVar: "CF_EMAIL",
		},
		cli.StringFlag{
			Name: "cf-token",
			Usage: "Cloudflare API token",
			EnvVar: "CF_TOKEN",
		},
		cli.StringFlag{
			Name: "zone",
			Usage: "DNS Zone domain",
			EnvVar: "DOMAIN",
		},
		cli.StringFlag{
			Name: "record-name, rname",
			Usage: "DNS Record name (used as 'record-name'.'zone')",
			EnvVar: "RECORD_NAME",
		},
		cli.StringFlag{
			Name: "interface, iface",
			Usage: "Local interface to use for addressing",
			EnvVar: "INTERFACE",
			Value: "eth0",
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

		return nil
	}

	app.Commands = commands

	app.Run(os.Args)
}

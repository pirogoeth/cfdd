package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/mailgun/holster/errors"
	"gopkg.in/urfave/cli.v1"

	"github.com/pirogoeth/cfdd/cfq"
	"github.com/pirogoeth/cfdd/util"
)

var CurrentCmd cli.Command = cli.Command{
	Name:      "current",
	Usage:     "Retrieve the current setting for the domain name",
	Aliases:   []string{"c"},
	ArgsUsage: "",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "check-interface",
			Usage: "Compare the DNS values against local interface values. If no update is needed, returns non-zero.",
		},
	},
	Action: current,
}

func current(ctx *cli.Context) error {
	// Expects four flags:
	//  cf-email
	//  cf-token
	//  zone
	//  rname
	cfEmail := ctx.GlobalString("cf-email")
	if cfEmail == "" {
		return cli.NewExitError("`--cf-email` required for displaying current domain info!", 1)
	}

	cfToken := ctx.GlobalString("cf-token")
	if cfToken == "" {
		return cli.NewExitError("`--cf-token` required for displaying current domain info!", 1)
	}

	zoneName := ctx.GlobalString("zone")
	if zoneName == "" {
		return cli.NewExitError("`--zone` required for displaying current domain info!", 1)
	}

	recordName := ctx.GlobalString("record-name")
	if recordName == "" {
		return cli.NewExitError("`--record-name` required for displaying current domain info!", 1)
	}

	checkIface := ctx.Bool("check-interface")
	ifaceName := ctx.GlobalString("interface")
	filterUnroutable := ctx.GlobalBool("filter-unroutable")

	log.WithField("cfEmail", cfEmail).Debugf("getting cloudflare client")
	cfApi, err := cfq.GetCloudflareClient(cfEmail, cfToken)
	if err != nil {
		return errors.Wrap(err, "while getting cloudflare client")
	}

	fqdn := util.BuildFQDN(recordName, zoneName)

	log.WithFields(log.Fields{
		"zone": zoneName,
		"fqdn": fqdn,
	}).Debugf("getting records for fqdn in zone")
	records, err := cfq.GetAddressesForZone(cfApi, zoneName, recordName)
	if err != nil {
		return errors.Wrap(err, "while getting addresses for zone")
	}

	log.Debugf("getting IPs for dns records")
	actual, err := cfq.DNSRecordsToNetIP(records)
	if err != nil {
		return errors.Wrap(err, "while getting IP address from DNS record")
	}

	fmt.Printf("%s:\n", fqdn)

	if len(actual) != 0 {
		for _, addr := range actual {
			fmt.Printf(" - %s\n", addr.String())
		}
	} else {
		fmt.Println("No addresses exist in Cloudflare DNS for this record")
	}

	fmt.Println()

	needsUpdate := false

	if checkIface {
		expected, err := util.GetAddressesForInterface(ifaceName)
		if err != nil {
			return errors.Wrap(err, "while getting addresses for interface")
		}

		if filterUnroutable {
			expected = util.FilterGloballyUnroutableAddrs(expected)

			if len(expected) == 0 {
				return cli.NewExitError("No expected addresses after filtering unroutable addresses", 127)
			}
		}

		fmt.Printf("Expected addresses [%s]:\n", ifaceName)
		for _, eip := range expected {
			if len(actual) == 0 {
				fmt.Printf(" - %s\n", util.Error(eip.String()))
				needsUpdate = true
				continue
			}

			for _, aip := range actual {
				if util.IsV4(aip) != util.IsV4(eip) {
					continue
				}

				if aip.Equal(eip) {
					fmt.Printf(" - %s\n", util.Okay(eip.String()))
				} else {
					needsUpdate = true
					fmt.Printf(" - %s\n", util.Warn(eip.String()))
				}
			}
		}
		fmt.Println()

		if !needsUpdate {
			return cli.NewExitError("No update needed", 1)
		}
	}

	return nil
}

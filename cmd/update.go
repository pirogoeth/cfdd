package cmd

import (
	"fmt"

	"github.com/mailgun/holster/errors"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/pirogoeth/cfdd/cfq"
	"github.com/pirogoeth/cfdd/util"
)

var UpdateCmd cli.Command = cli.Command{
	Name: "update",
	Usage: "Update the information in Cloudflare for a single domain name",
	Aliases: []string{"up"},
	ArgsUsage: "",
	Action: update,
}

func update(ctx *cli.Context) error {
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

	ifaceName := ctx.GlobalString("interface")

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
	actual, err := cfq.DNSRecordToNetIP(records)
	if err != nil {
		return errors.Wrap(err, "while getting IP address from DNS record")
	}

	expected, err := util.GetAddressesForInterface(ifaceName)
	if err != nil {
		return errors.Wrap(err, "while getting addresses for interface")
	}

	fmt.Printf("complete the impl")

	return nil
}

package cmd

import (
	"fmt"

	"github.com/mailgun/holster/errors"
	log "github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/pirogoeth/cfdd/cfq"
)

var ListRecordsCmd cli.Command = cli.Command{
	Name:      "list-records",
	Usage:     "List records on the configured zone",
	Aliases:   []string{"l"},
	ArgsUsage: "",
	Action:    listRecords,
}

func listRecords(ctx *cli.Context) error {
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

	log.WithField("cfEmail", cfEmail).Debugf("getting cloudflare client")
	cfApi, err := cfq.GetCloudflareClient(cfEmail, cfToken)
	if err != nil {
		return errors.Wrap(err, "while getting cloudflare client")
	}

	log.WithField("zone", zoneName).Debugf("getting DNS records")
	records, err := cfq.GetRecordsForZone(cfApi, zoneName)
	if err != nil {
		return errors.Wrap(err, "while fetching records for zone")
	}

	fmt.Printf("%s: DNS records\n", zoneName)
	for _, record := range records {
		fmt.Printf(" - Name: %s\n", record.Name)
		fmt.Printf("   Type: %s\n", record.Type)
		fmt.Printf("   TTL: %d\n", record.TTL)
		fmt.Printf("   Content: %s\n\n", record.Content)
	}

	return nil
}

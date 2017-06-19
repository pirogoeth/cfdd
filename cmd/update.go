package cmd

import (
	log "github.com/Sirupsen/logrus"
	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/mailgun/holster/errors"
	"gopkg.in/urfave/cli.v1"

	"github.com/pirogoeth/cfdd/cfq"
	"github.com/pirogoeth/cfdd/util"
)

var UpdateCmd cli.Command = cli.Command{
	Name:      "update",
	Usage:     "Update the information in Cloudflare for a single domain name",
	Aliases:   []string{"up"},
	ArgsUsage: "",
	Action:    update,
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
	filterUnroutable := ctx.GlobalBool("filter-unroutable")

	log.WithField("cfEmail", cfEmail).Debugf("getting cloudflare client")
	cfApi, err := cfq.GetCloudflareClient(cfEmail, cfToken)
	if err != nil {
		return errors.Wrap(err, "while getting cloudflare client")
	}

	fqdn := util.BuildFQDN(recordName, zoneName)

	zoneId, err := cfq.GetZoneId(cfApi, zoneName)
	if err != nil {
		return errors.Wrap(err, "while getting zone id during update")
	}

	log.WithFields(log.Fields{
		"zone": zoneName,
		"fqdn": fqdn,
	}).Debugf("getting records for fqdn in zone")
	records, err := cfq.GetAddressesForZone(cfApi, zoneName, recordName)
	if err != nil {
		return errors.Wrap(err, "while getting addresses for zone")
	}

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

	for _, expectAddr := range expected {
		updated := false

		for _, actualRec := range records {
			actualAddr, err := cfq.DNSRecordToNetIP(actualRec)
			if err != nil {
				return errors.Wrap(err, "while converting DNS record to IP addr during update")
			}

			if actualAddr.Equal(expectAddr) {
				updated = true
				break
			}

			if util.IsV4(actualAddr) != util.IsV4(expectAddr) {
				continue
			}

			log.WithFields(log.Fields{
				"zoneId":     zoneId,
				"recordId":   actualRec.ID,
				"recordName": actualRec.Name,
				"recordType": actualRec.Type,
			}).Debugf("Deleting record from zone")

			cfApi.DeleteDNSRecord(zoneId, actualRec.ID)
		}

		if !updated {
			var newRecord = cloudflare.DNSRecord{
				Name:     recordName,
				Content:  expectAddr.String(),
				ZoneID:   zoneId,
				ZoneName: zoneName,
				Proxied:  false,
			}

			if util.IsV4(expectAddr) {
				newRecord.Type = "A"
			} else {
				newRecord.Type = "AAAA"
			}

			log.WithFields(log.Fields{
				"recordName": recordName,
				"zoneName":   zoneName,
				"content":    expectAddr.String(),
			}).Debugf("Creating record in cloudflare")

			resp, err := cfApi.CreateDNSRecord(zoneId, newRecord)
			if err != nil {
				return errors.WithContext{
					"apiResponse": resp,
				}.Wrapf(err, "while creating new record")
			}
		}
	}

	return nil
}

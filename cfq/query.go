package cfq

import (
	"net"

	"github.com/cloudflare/cloudflare-go"
	"github.com/mailgun/holster/errors"
	log "github.com/Sirupsen/logrus"

	"github.com/pirogoeth/cfdd/util"
)

func DNSRecordToNetIP(dnames []cloudflare.DNSRecord) ([]net.IP, error) {
	var actual []net.IP

	for _, rec := range dnames {
		ip := net.ParseIP(rec.Content)
		if ip == nil {
			return nil, errors.WithContext{
				"content": rec.Content,
			}.Error("while parsing IP addr from cloudflare DNS record")
		}

		actual = append(actual, ip)
	}

	return actual, nil
}

func GetAddressesForZone(cfApi *cloudflare.API, zoneName, recordName string) ([]cloudflare.DNSRecord, error) {
	log.WithField("zoneName", zoneName).Debugf("Fetching zone id")
	zoneId, err := cfApi.ZoneIDByName(zoneName)
	if err != nil {
		return nil, errors.WithContext{
			"zoneName": zoneName,
		}.Wrap(err, "while fetching zone id")
	}

	fqdn := util.BuildFQDN(recordName, zoneName)

	log.WithField("recordName", recordName).Debugf("Building DNS records (A/AAAA) for query")
	recordA := cloudflare.DNSRecord{
		Type: "A",
		Name: fqdn,
	}
	recordA4 := cloudflare.DNSRecord{
		Type: "AAAA",
		Name: fqdn,
	}

	log.WithField("recordType", "A").Debugf("Querying for records")
	cfA, err := cfApi.DNSRecords(zoneId, recordA)
	if err != nil {
		return nil, errors.WithContext{
			"zoneName": zoneName,
			"zoneId": zoneId,
			"recordQuery": recordA,
		}.Wrap(err, "while querying for A record")
	}

	log.WithField("recordType", "AAAA").Debugf("Querying for records")
	cfA4, err := cfApi.DNSRecords(zoneId, recordA4)
	if err != nil {
		return nil, errors.WithContext{
			"zoneName": zoneName,
			"zoneId": zoneId,
			"recordQuery": recordA,
		}.Wrap(err, "while querying for AAAA record")
	}

	return append(cfA, cfA4...), nil
}

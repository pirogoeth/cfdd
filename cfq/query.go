package cfq

import (
	"net"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/mailgun/holster/errors"
	log "github.com/sirupsen/logrus"

	"github.com/pirogoeth/cfdd/util"
)

func DNSRecordToNetIP(dname cloudflare.DNSRecord) (net.IP, error) {
	ip := net.ParseIP(dname.Content)
	if ip == nil {
		return nil, errors.WithContext{
			"content": dname.Content,
		}.Error("while parsing IP addr from cloudflare DNS record")
	}

	return ip, nil
}

func DNSRecordsToNetIP(dnames []cloudflare.DNSRecord) ([]net.IP, error) {
	var actual []net.IP

	for _, rec := range dnames {
		ip, err := DNSRecordToNetIP(rec)
		if err != nil {
			return nil, errors.Wrap(err, "while parsing IP addrs from DNS record list")
		}
		actual = append(actual, ip)
	}

	return actual, nil
}

func GetZoneId(cfApi *cloudflare.API, zoneName string) (string, error) {
	log.WithField("zoneName", zoneName).Debugf("Fetching zone id")
	zoneId, err := cfApi.ZoneIDByName(zoneName)
	if err != nil {
		return "", errors.WithContext{
			"zoneName": zoneName,
		}.Wrap(err, "while fetching zone id")
	}

	return zoneId, nil
}

func GetRecordsForZone(cfApi *cloudflare.API, zoneName string) ([]cloudflare.DNSRecord, error) {
	zoneId, err := GetZoneId(cfApi, zoneName)
	if err != nil {
		return nil, errors.WithContext{
			"zoneName": zoneName,
		}.Wrap(err, "while querying records for zone")
	}

	log.WithField("zoneId", zoneId).Debugf("Fetching records for zone")
	recordA := cloudflare.DNSRecord{
		Type: "A",
	}
	recordA4 := cloudflare.DNSRecord{
		Type: "AAAA",
	}

	log.WithField("recordType", "A").Debugf("Querying for records")
	cfA, err := cfApi.DNSRecords(zoneId, recordA)
	if err != nil {
		return nil, errors.WithContext{
			"zoneName":    zoneName,
			"zoneId":      zoneId,
			"recordQuery": recordA,
		}.Wrap(err, "while querying for A record")
	}

	log.WithField("recordType", "AAAA").Debugf("Querying for records")
	cfA4, err := cfApi.DNSRecords(zoneId, recordA4)
	if err != nil {
		return nil, errors.WithContext{
			"zoneName":    zoneName,
			"zoneId":      zoneId,
			"recordQuery": recordA,
		}.Wrap(err, "while querying for AAAA record")
	}

	return append(cfA, cfA4...), nil
}

func GetAddressesForZone(cfApi *cloudflare.API, zoneName, recordName string) ([]cloudflare.DNSRecord, error) {
	zoneId, err := GetZoneId(cfApi, zoneName)
	if err != nil {
		return nil, errors.Wrap(err, "while fetching addresses for zone")
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
			"zoneName":    zoneName,
			"zoneId":      zoneId,
			"recordQuery": recordA,
		}.Wrap(err, "while querying for A record")
	}

	log.WithField("recordType", "AAAA").Debugf("Querying for records")
	cfA4, err := cfApi.DNSRecords(zoneId, recordA4)
	if err != nil {
		return nil, errors.WithContext{
			"zoneName":    zoneName,
			"zoneId":      zoneId,
			"recordQuery": recordA,
		}.Wrap(err, "while querying for AAAA record")
	}

	return append(cfA, cfA4...), nil
}

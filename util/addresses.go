package util

import (
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/mailgun/holster/errors"
)

var unroutable []*net.IPNet

func init() {
	unroutableNetStrs := []string{
		"10.0.0.0/8",
		"100.64.0.0/10",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.168.0.0/16",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"240.0.0.0/4",
		"255.255.255.255/32",
		"::/128",
		"::1/128",
		"100::/64",
		"2001::/32",
		"2001:10::/28",
		"2001:20::/28",
		"2001:db8::/32",
		"fc00::/7",
	}

	for _, netStr := range unroutableNetStrs {
		_, ipn, err := net.ParseCIDR(netStr)
		if err != nil {
			panic(fmt.Sprintf("Could not parse network CIDR %s", netStr))
		}
		unroutable = append(unroutable, ipn)
	}
}

func GetAddressesForInterface(ifaceName string) ([]net.IP, error) {
	addrs, err := GetIfaceAddrs(ifaceName)
	if err != nil {
		return nil, errors.WithContext{
			"ifaceName": ifaceName,
		}.Wrap(err, "while fetching local interface addresses")
	}

	var expected []net.IP

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, errors.Wrap(err, "while parsing IP addr from local interface")
		}

		if ip == nil {
			return nil, errors.WithContext{
				"content": addr.String(),
			}.Error("while parsing IP addr from local interface")
		}

		expected = append(expected, ip)
	}

	return expected, nil
}

func FilterGloballyUnroutableAddrs(addrs []net.IP) []net.IP {
	var addrsGlobal []net.IP = make([]net.IP, 0)

	for _, addr := range addrs {
		if addr.IsInterfaceLocalMulticast() {
			log.Debugf("Found unroutable address %s", addr.String())
			continue
		} else if addr.IsLinkLocalMulticast() {
			log.Debugf("Found unroutable address %s", addr.String())
			continue
		} else if addr.IsLinkLocalUnicast() {
			log.Debugf("Found unroutable address %s", addr.String())
			continue
		} else if addr.IsLoopback() {
			log.Debugf("Found unroutable address %s", addr.String())
			continue
		}

		isUnroutable := false

		for _, network := range unroutable {
			if network.Contains(addr) {
				log.Debugf("Found unroutable address %s in network %s", addr.String(), network.String())
				isUnroutable = true
				break
			}
		}

		if isUnroutable {
			continue
		}

		addrsGlobal = append(addrsGlobal, addr)
	}

	return addrsGlobal
}

func IsV4(addr net.IP) bool {
	return addr.To4() != nil
}

func IsV6(addr net.IP) bool {
	return !IsV4(addr)
}

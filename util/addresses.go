package util

import (
	"net"

	"github.com/mailgun/holster/errors"
)

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

package util

import (
	"net"

	"github.com/mailgun/holster/errors"
)

func GetIfaceAddrs(ifaceName string) ([]net.Addr, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, errors.WithContext{
			"ifaceName": ifaceName,
		}.Wrap(err, "while looking up local interface")
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, errors.WithContext{
			"ifaceName": ifaceName,
		}.Wrap(err, "while looking up addresses for interface")
	}

	return addrs, nil
}

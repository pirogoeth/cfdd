package cfq

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/mailgun/holster/errors"
	log "github.com/Sirupsen/logrus"
)

func GetCloudflareClient(cfEmail, cfToken string) (*cloudflare.API, error) {
	log.WithField("cfEmail", cfEmail).Debugf("Creating cloudflare client")
	cfApi, err := cloudflare.New(cfToken, cfEmail)
	if err != nil {
		return nil, errors.Wrap(err, "while creating cloudflare client")
	}

	return cfApi, nil
}

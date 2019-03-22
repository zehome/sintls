package dns

import (
	"fmt"
	"github.com/zehome/sintls/dns/ovh"
)

type DNSUpdater interface {
	SetRecord(fqdn string, fieldtype string, target string) error
	RemoveRecord(fqdn string, fieldtype string) error
	ExtractRecordName(fqdn, domain string) string
	ExtractAuthZone(fqdn string) (string, error)
}

func NewDNSUpdaterByName(name string) (DNSUpdater, error) {
	switch name {
	case "ovh":
		return ovh.NewDNSProvider()
	default:
		return nil, fmt.Errorf("unrecognized DNS provider: %s", name)
	}
}

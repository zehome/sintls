package dns

import (
	"fmt"
	"github.com/zehome/sintls/dns/ovh"
	"github.com/zehome/sintls/dns/rfc2136"
)

type DNSUpdater interface {
	SetRecord(fqdn string, fieldtype string, target string) error
	RemoveRecords(fqdn string) error
	Refresh(fqdn string) error
}

func NewDNSUpdaterByName(name string) (DNSUpdater, error) {
	switch name {
	case "rfc2136":
		return rfc2136.NewDNSProvider()
	case "ovh":
		return ovh.NewDNSProvider()
	default:
		return nil, fmt.Errorf("unrecognized DNS provider: %s", name)
	}
}

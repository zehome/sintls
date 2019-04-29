package ovh

import (
	"net/http"
	"fmt"
	"log"
	"errors"
	"strings"
	"github.com/go-acme/lego/challenge/dns01"
	"github.com/go-acme/lego/platform/config/env"
	"github.com/ovh/go-ovh/ovh"
)

// Record a DNS record
type Record struct {
	ID        int    `json:"id,omitempty"`
	FieldType string `json:"fieldType,omitempty"`
	SubDomain string `json:"subDomain,omitempty"`
	Target    string `json:"target,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	Zone      string `json:"zone,omitempty"`
}

// Config is used to configure the creation of the DNSProvider
type Config struct {
	APIEndpoint        string
	ApplicationKey     string
	ApplicationSecret  string
	ConsumerKey        string
	TTL                int
	HTTPClient         *http.Client
}

// NewDefaultConfig returns a default configuration for the DNSProvider
func NewDefaultConfig() *Config {
	return &Config{
		TTL: 180,
		HTTPClient: &http.Client{
			Timeout: env.GetOrDefaultSecond("OVH_HTTP_TIMEOUT", ovh.DefaultTimeout),
		},
	}
}

// OVHApi is an implementation of the acme.ChallengeProvider interface
// that uses OVH's REST API to manage TXT records for a domain.
type OVHApi struct {
	config      *Config
	client      *ovh.Client
}

func arrayToString(a []int, delim string) string {
    return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

// NewDNSProvider returns a OVHApi instance configured for OVH
// Credentials must be passed in the environment variable:
// OVH_ENDPOINT : it must be ovh-eu or ovh-ca
// OVH_APPLICATION_KEY
// OVH_APPLICATION_SECRET
// OVH_CONSUMER_KEY
func NewDNSProvider() (*OVHApi, error) {
	values, err := env.Get("OVH_ENDPOINT", "OVH_APPLICATION_KEY", "OVH_APPLICATION_SECRET", "OVH_CONSUMER_KEY")
	if err != nil {
		return nil, fmt.Errorf("ovh: %v", err)
	}
	config := NewDefaultConfig()
	config.APIEndpoint = values["OVH_ENDPOINT"]
	config.ApplicationKey = values["OVH_APPLICATION_KEY"]
	config.ApplicationSecret = values["OVH_APPLICATION_SECRET"]
	config.ConsumerKey = values["OVH_CONSUMER_KEY"]
	return NewDNSProviderConfig(config)
}

// NewDNSProviderConfig return a OVHApi instance configured for OVH.
func NewDNSProviderConfig(config *Config) (*OVHApi, error) {
	if config == nil {
		return nil, errors.New("ovh: the configuration of the DNS provider is nil")
	}

	if config.APIEndpoint == "" || config.ApplicationKey == "" || config.ApplicationSecret == "" || config.ConsumerKey == "" {
		return nil, fmt.Errorf("ovh: credentials missing")
	}
	client, err := ovh.NewClient(
		config.APIEndpoint,
		config.ApplicationKey,
		config.ApplicationSecret,
		config.ConsumerKey,
	)
	if err != nil {
		return nil, fmt.Errorf("ovh: %v", err)
	}
	client.Client = config.HTTPClient
	return &OVHApi{
		config:    config,
		client:    client,
	}, nil
}

func (o *OVHApi) ExtractRecordName(fqdn, domain string) string {
	name := dns01.UnFqdn(fqdn)
	if idx := strings.Index(name, "."+domain); idx != -1 {
		return name[:idx]
	}
	return name
}

func (o *OVHApi) ExtractAuthZone(fqdn string) (string, error) {
	// Parse domain name
	authZone, err := dns01.FindZoneByFqdn(dns01.ToFqdn(fqdn))
	if err != nil {
		return "", fmt.Errorf("ovh: could not determine zone for domain: '%s'. %s", fqdn, err)
	}
	return dns01.UnFqdn(authZone), nil
}

func (api *OVHApi) SetRecord(fqdn, fieldtype, target string) error {
	allowedfieldtypes := map[string]bool {
		"A": true,
		"AAAA": true,
		"CNAME": true,
	}
	if ! allowedfieldtypes[fieldtype] {
		return fmt.Errorf("ovh: fieldtype %s not supported.", fieldtype)
	}
	// Parse domain name
	authZone, err := api.ExtractAuthZone(fqdn)
	if err != nil {
		return fmt.Errorf("ovh: could not determine zone for domain: '%s'. %s", fqdn, err)
	}
	subDomain := api.ExtractRecordName(fqdn, authZone)
	reqURL := fmt.Sprintf("/domain/zone/%s/record", authZone)
	reqData := Record{
		FieldType: fieldtype,
		SubDomain: subDomain,
		Target: target,
		TTL: 300, // seconds
	}
	var respData Record
	err = api.client.Post(reqURL, reqData, &respData)
	if err != nil {
		return fmt.Errorf("ovh: error when call api to add record (%s): %v", reqURL, err)
	}
	return nil
}

func (api *OVHApi) RemoveRecords(fqdn string) error {
	authZone, err := dns01.FindZoneByFqdn(dns01.ToFqdn(fqdn))
	if err != nil {
		return fmt.Errorf("ovh: could not determine zone for domain: '%s'. %s", fqdn, err)
	}
	authZone = dns01.UnFqdn(authZone)
	subDomain := api.ExtractRecordName(fqdn, authZone)
	reqURL := fmt.Sprintf("/domain/zone/%s/record?subDomain=%s", authZone, subDomain)
	recordids := []int{}
	err = api.client.Get(reqURL, &recordids)
	if err != nil {
		return fmt.Errorf("ovh: error when call api to get record (%s): %v", reqURL, err)
	}
	log.Printf("ovh: remove recods on zone=%s subdomain=%s: %s\n", authZone, subDomain, arrayToString(recordids, ","))
	for _, recordid := range recordids {
		// log.Printf("ovh: remove record id=%s", recordid)
		err = api.client.Delete(
			fmt.Sprintf("/domain/zone/%s/record/%d", authZone, recordid), nil)
		if err != nil {
			log.Printf("ovh: unable to remove record %s\n", recordid)
		}
	}
	return nil
}

func (api *OVHApi) Refresh(fqdn string) error {
	authZone, err := dns01.FindZoneByFqdn(dns01.ToFqdn(fqdn))
	if err != nil {
		return fmt.Errorf("ovh: could not determine zone for domain: '%s'. %s", fqdn, err)
	}
	authZone = dns01.UnFqdn(authZone)
	err = api.client.Post(fmt.Sprintf("/domain/zone/%s/refresh", authZone), nil, nil)
	if err != nil {
		log.Printf("ovh: Refresh %s failed: %s\n", authZone, err)
	}
	return err
}
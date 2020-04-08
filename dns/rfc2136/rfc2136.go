package rfc2136

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-acme/lego/v3/challenge/dns01"
	"github.com/go-acme/lego/v3/platform/config/env"
	"github.com/miekg/dns"
)

// Environment variables names.
const (
	envNamespace = "RFC2136_"

	EnvTSIGKey       = envNamespace + "TSIG_KEY"
	EnvTSIGSecret    = envNamespace + "TSIG_SECRET"
	EnvTSIGAlgorithm = envNamespace + "TSIG_ALGORITHM"
	EnvNameserver    = envNamespace + "NAMESERVER"
	EnvDNSTimeout    = envNamespace + "DNS_TIMEOUT"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvSequenceInterval   = envNamespace + "SEQUENCE_INTERVAL"
)

// Config is used to configure the creation of the DNSProvider
type Config struct {
	Nameserver         string
	TSIGAlgorithm      string
	TSIGKey            string
	TSIGSecret         string
	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	SequenceInterval   time.Duration
	DNSTimeout         time.Duration
}

// NewDefaultConfig returns a default configuration for the DNSProvider
func NewDefaultConfig() *Config {
	return &Config{
		TSIGAlgorithm:      env.GetOrDefaultString(EnvTSIGAlgorithm, dns.HmacMD5),
		TTL:                env.GetOrDefaultInt(EnvTTL, dns01.DefaultTTL),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, env.GetOrDefaultSecond("RFC2136_TIMEOUT", 60*time.Second)),
		PollingInterval:    env.GetOrDefaultSecond(EnvPollingInterval, 2*time.Second),
		SequenceInterval:   env.GetOrDefaultSecond(EnvSequenceInterval, dns01.DefaultPropagationTimeout),
		DNSTimeout:         env.GetOrDefaultSecond(EnvDNSTimeout, 10*time.Second),
	}
}

// DNSProvider is an implementation of the challenge.Provider interface that
// uses dynamic DNS updates (RFC 2136) to create TXT records on a nameserver.
type DNSProvider struct {
	config *Config
}

// NewDNSProvider returns a DNSProvider instance configured for rfc2136
// dynamic update. Configured with environment variables:
// RFC2136_NAMESERVER: Network address in the form "host" or "host:port".
// RFC2136_TSIG_ALGORITHM: Defaults to hmac-md5.sig-alg.reg.int. (HMAC-MD5).
// See https://github.com/miekg/dns/blob/master/tsig.go for supported values.
// RFC2136_TSIG_KEY: Name of the secret key as defined in DNS server configuration.
// RFC2136_TSIG_SECRET: Secret key payload.
// RFC2136_PROPAGATION_TIMEOUT: DNS propagation timeout in time.ParseDuration format. (60s)
// To disable TSIG authentication, leave the RFC2136_TSIG* variables unset.
func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvNameserver)
	if err != nil {
		return nil, fmt.Errorf("rfc2136: %w", err)
	}

	config := NewDefaultConfig()
	config.Nameserver = values[EnvNameserver]
	config.TSIGKey = env.GetOrFile(EnvTSIGKey)
	config.TSIGSecret = env.GetOrFile(EnvTSIGSecret)

	return NewDNSProviderConfig(config)
}

// NewDNSProviderConfig return a DNSProvider instance configured for rfc2136.
func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("rfc2136: the configuration of the DNS provider is nil")
	}

	if config.Nameserver == "" {
		return nil, errors.New("rfc2136: nameserver missing")
	}

	if config.TSIGAlgorithm == "" {
		config.TSIGAlgorithm = dns.HmacMD5
	}

	// Append the default DNS port if none is specified.
	if _, _, err := net.SplitHostPort(config.Nameserver); err != nil {
		if strings.Contains(err.Error(), "missing port") {
			config.Nameserver = net.JoinHostPort(config.Nameserver, "53")
		} else {
			return nil, fmt.Errorf("rfc2136: %w", err)
		}
	}

	if len(config.TSIGKey) == 0 && len(config.TSIGSecret) > 0 ||
		len(config.TSIGKey) > 0 && len(config.TSIGSecret) == 0 {
		config.TSIGKey = ""
		config.TSIGSecret = ""
	}

	return &DNSProvider{config: config}, nil
}

// Present creates a record using the specified parameters
func (d *DNSProvider) SetRecord(fqdn, fieldtype, target string) error {
	allowedFieldTypes := map[string]bool{
		"A":     true,
		"AAAA":  true,
		"CNAME": true,
	}
	if !allowedFieldTypes[fieldtype] {
		return fmt.Errorf("rfc2136: fieldtype %s not supported.", fieldtype)
	}

	err := d.changeRecord("INSERT", fieldtype, fmt.Sprintf("%s.", fqdn), target, d.config.TTL)
	if err != nil {
		return fmt.Errorf("rfc2136: failed to insert: %w", err)
	}
	return nil
}

// Noop using rfc2136
func (d *DNSProvider) Refresh(fqdn string) error {
	return nil
}

func (d *DNSProvider) RemoveRecords(fqdn string) error {
	for _, fieldtype := range []string{"A", "AAAA", "CNAME"} {
		_ = d.changeRecord("REMOVE", fieldtype, fmt.Sprintf("%s.", fqdn), "", d.config.TTL)
	}
	return nil
}

func (d *DNSProvider) changeRecord(action, fieldtype string, fqdn, target string, ttl int) error {
	// Find the zone for the given fqdn
	zone, err := dns01.FindZoneByFqdnCustom(fqdn, []string{d.config.Nameserver})
	if err != nil {
		return err
	}

	// Create RR
	var rrs []dns.RR
	switch fieldtype {
	case "A":
		rr := new(dns.A)
		rr.Hdr = dns.RR_Header{Name: fqdn, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: uint32(ttl)}
		rr.A = net.ParseIP(target)
		rrs = append(rrs, rr)
	case "AAAA":
		rr := new(dns.AAAA)
		rr.Hdr = dns.RR_Header{Name: fqdn, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: uint32(ttl)}
		rr.AAAA = net.ParseIP(target)
		rrs = append(rrs, rr)
	case "CNAME":
		rr := new(dns.CNAME)
		rr.Hdr = dns.RR_Header{Name: fqdn, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: uint32(ttl)}
		rr.Target = target
		rrs = append(rrs, rr)
	default:
		return fmt.Errorf("rfc2136: unsupported field type %s", fieldtype)
	}

	// Create dynamic update packet
	m := new(dns.Msg)
	m.SetUpdate(zone)
	switch action {
	case "INSERT":
		// Always remove old challenge left over from who knows what.
		m.RemoveRRset(rrs)
		m.Insert(rrs)
	case "REMOVE":
		m.Remove(rrs)
	default:
		return fmt.Errorf("unexpected action: %s", action)
	}

	// Setup client
	c := &dns.Client{Timeout: d.config.DNSTimeout}
	c.SingleInflight = true

	// TSIG authentication / msg signing
	if len(d.config.TSIGKey) > 0 && len(d.config.TSIGSecret) > 0 {
		m.SetTsig(dns.Fqdn(d.config.TSIGKey), d.config.TSIGAlgorithm, 300, time.Now().Unix())
		c.TsigSecret = map[string]string{dns.Fqdn(d.config.TSIGKey): d.config.TSIGSecret}
	}

	// Send the query
	reply, _, err := c.Exchange(m, d.config.Nameserver)
	if err != nil {
		return fmt.Errorf("DNS update failed: %w", err)
	}
	if reply != nil && reply.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("DNS update failed: server replied: %s", dns.RcodeToString[reply.Rcode])
	}

	return nil
}

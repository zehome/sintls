// Package httpreq implements a DNS provider for solving the DNS-01 challenge through a HTTP server.
package sintlsprovider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
	"github.com/urfave/cli"
)

var UserAgent string

type message struct {
	Domain      string `json:"domain"`
	Token       string `json:"token"`
	KeyAuth     string `json:"keyAuth"`
	TargetA     string `json:"dnstarget_a"`
	TargetAAAA  string `json:"dnstarget_aaaa"`
	TargetCNAME string `json:"dnstarget_cname"`
	TargetMX    string `json:"dnstarget_mx"`
}

// Config is used to configure the creation of the Provider
type Config struct {
	Endpoint           *url.URL
	Username           string
	Password           string
	TargetA            string
	TargetAAAA         string
	TargetCNAME        string
	TargetMX           string
	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	HTTPClient         *http.Client
}

// NewDefaultConfig returns a default configuration for the Provider
func NewDefaultConfig() *Config {
	return &Config{
		PropagationTimeout: 180 * time.Second,
		PollingInterval:    dns01.DefaultPollingInterval,
		HTTPClient: &http.Client{
			Timeout: env.GetOrDefaultSecond("SINTLS_HTTP_TIMEOUT", 30*time.Second),
		},
	}
}

// Provider describes a provider for acme-proxy
type Provider struct {
	config *Config
}

func (d *Provider) GetTargetA() string {
	return d.config.TargetA
}
func (d *Provider) GetTargetAAAA() string {
	return d.config.TargetAAAA
}
func (d *Provider) GetTargetCNAME() string {
	return d.config.TargetCNAME
}
func (d *Provider) GetTargetMX() string {
	return d.config.TargetMX
}

// NewProvider returns a Provider instance.
func NewProvider(ctx *cli.Context) (*Provider, error) {
	endpoint, err := url.Parse(ctx.GlobalString("server"))
	if err != nil {
		return nil, fmt.Errorf("sintls: %v", err)
	}
	config := NewDefaultConfig()
	config.Username = os.Getenv("SINTLS_USERNAME")
	config.Password = os.Getenv("SINTLS_PASSWORD")
	config.TargetA = ctx.GlobalString("target-a")
	config.TargetAAAA = ctx.GlobalString("target-aaaa")
	config.TargetCNAME = ctx.GlobalString("target-cname")
	config.TargetMX = ctx.GlobalString("target-mx")
	config.Endpoint = endpoint
	return NewProviderConfig(config)
}

// NewProviderConfig return a Provider .
func NewProviderConfig(config *Config) (*Provider, error) {
	if config == nil {
		return nil, errors.New("sintls: the configuration of the DNS provider is nil")
	}

	if config.Endpoint == nil {
		return nil, errors.New("sintls: the endpoint is missing")
	}

	return &Provider{config: config}, nil
}

// Timeout returns the timeout and interval to use when checking for DNS propagation.
// Adjusting here to cope with spikes in propagation times.
func (d *Provider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

// Just update DNS target records
func (d *Provider) UpdateDNS(domain string) error {
	msg := &message{
		Domain:      domain,
		Token:       "",
		KeyAuth:     "",
		TargetA:     d.GetTargetA(),
		TargetAAAA:  d.GetTargetAAAA(),
		TargetCNAME: d.GetTargetCNAME(),
		TargetMX:    d.GetTargetMX(),
	}
	err := d.doPost("/updatedns", msg)
	if err != nil {
		return fmt.Errorf("sintls: %v", err)
	}
	return nil
}

// Present creates a TXT record to fulfill the dns-01 challenge
func (d *Provider) Present(domain, token, keyAuth string) error {
	msg := &message{
		Domain:      domain,
		Token:       token,
		KeyAuth:     keyAuth,
		TargetA:     d.GetTargetA(),
		TargetAAAA:  d.GetTargetAAAA(),
		TargetCNAME: d.GetTargetCNAME(),
		TargetMX:    d.GetTargetMX(),
	}
	err := d.doPost("/present", msg)
	if err != nil {
		return fmt.Errorf("sintls: %v", err)
	}
	return nil
}

// CleanUp removes the TXT record matching the specified parameters
func (d *Provider) CleanUp(domain, token, keyAuth string) error {
	msg := &message{
		Domain:      domain,
		Token:       token,
		KeyAuth:     keyAuth,
		TargetA:     d.GetTargetA(),
		TargetAAAA:  d.GetTargetAAAA(),
		TargetCNAME: d.GetTargetCNAME(),
		TargetMX:    d.GetTargetMX(),
	}
	err := d.doPost("/cleanup", msg)
	if err != nil {
		return fmt.Errorf("sintls: %v", err)
	}
	return nil
}

func (d *Provider) doPost(uri string, msg interface{}) error {
	reqBody := &bytes.Buffer{}
	err := json.NewEncoder(reqBody).Encode(msg)
	if err != nil {
		return err
	}

	newURI := path.Join(d.config.Endpoint.EscapedPath(), uri)
	endpoint, err := d.config.Endpoint.Parse(newURI)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)
	if len(d.config.Username) > 0 && len(d.config.Password) > 0 {
		req.SetBasicAuth(d.config.Username, d.config.Password)
	}

	resp, err := d.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%d: failed to read response body: %v", resp.StatusCode, err)
		}

		return fmt.Errorf("%d: request failed: %v", resp.StatusCode, string(body))
	}

	return nil
}

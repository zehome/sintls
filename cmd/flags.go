package cmd

import (
	"github.com/go-acme/lego/lego"
	"github.com/urfave/cli"
)

func CreateFlags(defaultPath string) []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:  "domains, d",
			Usage: "Add a domain to the process. Can be specified multiple times.",
		},
		cli.StringFlag{
			Name:  "server",
			Usage: "SINTLS server",
			Value: "https://auth.clarilab.fr/",
		},
		cli.StringFlag{
			Name:  "target-a",
			Usage: "DNS Target (A Entry)",
		},
		cli.StringFlag{
			Name:  "target-aaaa",
			Usage: "DNS Target (AAAA Entry)",
		},
		cli.StringFlag{
			Name:  "target-cname",
			Usage: "DNS Target (CNAME Entry)",
		},
		cli.StringFlag{
			Name:  "ca-server",
			Usage: "CA hostname (and optionally :port).",
			Value: lego.LEDirectoryProduction,
		},
		cli.BoolFlag{
			Name:  "accept-tos, a",
			Usage: "By setting this flag to true you indicate that you accept the current Let's Encrypt terms of service.",
		},
		cli.StringFlag{
			Name:  "email, m",
			Usage: "Email used for registration and recovery contact.",
		},
		cli.StringFlag{
			Name:  "key-type, k",
			Value: "ec384",
			Usage: "Key type to use for private keys. Supported: rsa2048, rsa4096, rsa8192, ec256, ec384.",
		},
		cli.StringFlag{
			Name:  "path",
			Usage: "Directory to use for storing the data.",
			Value: defaultPath,
		},
		cli.IntFlag{
			Name:  "cert.timeout",
			Usage: "Set the certificate timeout value to a specific value in seconds. Only used when obtaining certificates.",
			Value: 30,
		},
		cli.BoolFlag{
			Name: "disable-selfupdate",
			Usage: "Disable self update mecanism using github",
			EnvVar: "SINTLS_DISABLE_SELFUPDATE",
		},
	}
}

package main

import (
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/log"
	"github.com/urfave/cli"
	"github.com/zehome/sintls/provider"
)

func setupHttpReq(ctx *cli.Context, client *lego.Client) {
	provider, err := sintlsprovider.NewProvider(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		log.Fatal(err)
	}
}

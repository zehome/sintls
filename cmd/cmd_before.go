package cmd

import (
	"github.com/urfave/cli"
	"log"
)


func Before(ctx *cli.Context) error {
	if len(ctx.GlobalString("path")) == 0 {
		log.Fatal("Could not determine current working directory. Please pass --path.")
	}
	err := createNonExistingFolder(ctx.GlobalString("path"))
	if err != nil {
		log.Fatalf("Could not check/create path: %v", err)
	}

	if len(ctx.GlobalString("ca-server")) == 0 {
		log.Fatal("Could not determine current working ca-server. Please pass --ca-server.")
	}
	if len(ctx.GlobalString("server")) == 0 {
		log.Fatal("Could not determine current sintls server. Please passe --server.")
	}
	return nil
}

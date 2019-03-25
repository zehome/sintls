package cmd

import (
	"github.com/go-acme/lego/log"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli"
)

func doSelfUpdate(ctx *cli.Context) {
	v, err := semver.Parse(ctx.App.Version)
	if err != nil {
		log.Printf("Selfupdate of non official builds are not supported: %s", err)
		return
	}
	latest, err := selfupdate.UpdateCommand("sintls", v, "zehome/sintls")
	if err != nil {
		log.Println("Binary update failed:", err)
		return
	}
	if !latest.Version.Equals(v) {
		log.Println("Successfully updated to version", latest.Version)
		log.Println("Release note:\n", latest.ReleaseNotes)
	}
}

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
	if !ctx.GlobalBool("disable-self-update") {
		doSelfUpdate(ctx)
	}
	return nil
}

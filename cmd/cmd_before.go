package cmd

import (
	"bufio"
	"fmt"
	"github.com/go-acme/lego/log"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli"
	"os"
)

func doSelfUpdate(ctx *cli.Context) {
	v, err := semver.Parse(ctx.App.Version)
	if err != nil {
		log.Printf("Selfupdate of non official builds are not supported: %s", err)
		return
	}
    latest, found, err := selfupdate.DetectLatest("zehome/sintls")
    if err != nil {
        log.Println("Error occurred while detecting version:", err)
        return
    }

    if !found || latest.Version.LTE(v) {
        return
    }

    fmt.Print("Do you want to update to v", latest.Version, "? (y/n): ")
    input, err := bufio.NewReader(os.Stdin).ReadString('\n')
    if err != nil || (input != "y\n" && input != "n\n") {
        return
    }
    if input == "n\n" {
        return
    }

	exe, err := os.Executable()
    if err != nil {
        log.Println("Could not locate executable path")
        return
    }
    if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
        log.Println("Error occurred while updating binary:", err)
        return
    }
    log.Println("Successfully updated to version", latest.Version)
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

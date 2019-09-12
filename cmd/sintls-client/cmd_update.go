package main

import (
	"bufio"
	"fmt"
	"github.com/blang/semver"
	"github.com/go-acme/lego/v3/log"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli"
	"os"
	"strings"
)

func doSelfUpdate(ctx *cli.Context) {
	v, err := semver.Parse(ctx.App.Version)
	if err != nil {
		log.Printf("Only official builds are not supported: %s", err)
		return
	}
	if strings.Contains(ctx.App.Version, "-dev") {
		log.Printf("Do not self update development versions!")
		return
	}
	log.Printf("version: %s (%s)", ctx.App.Version, v)
	latest, found, err := selfupdate.DetectLatest("zehome/sintls")
	if err != nil {
		log.Println("Error occurred while detecting version:", err)
		return
	}

	if !found || latest.Version.LTE(v) {
		return
	}
	if !ctx.GlobalBool("unattended") {
		fmt.Print("Do you want to update to v", latest.Version, "? (y/n): ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil || (input != "y\n" && input != "n\n") {
			return
		}
		if input == "n\n" {
			return
		}
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
	log.Println("Successfully updated to version:", latest.Version)
}

func createUpdate() cli.Command {
	return cli.Command{
		Name:   "update",
		Usage:  "Update on github",
		Action: doSelfUpdate,
		Flags:	[]cli.Flag{
			cli.BoolFlag{
				Name:  "unattended",
				Usage: "Do unattended upgrade, without asking",
			},
		},
	}
}

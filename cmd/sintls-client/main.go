package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli"
	"github.com/zehome/sintls/provider"
)

var (
	version = "1.0.0"
)

func main() {
	usr, _ := user.Current()
	app := cli.NewApp()
	app.Name = "sintls"
	app.HelpName = "sintls"
	app.Usage = "Simple Internal TLS certificate helper"
	app.EnableBashCompletion = true
	app.Version = version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("sintls version %s %s/%s\n", c.App.Version, runtime.GOOS, runtime.GOARCH)
	}
	if usr == nil || len(usr.HomeDir) == 0 {
		log.Fatal("Unable to determine home directory. Is USER environment variable defined?")
	}
	defaultPath := filepath.Join(usr.HomeDir, ".config", "sintls")
	app.Flags = CreateFlags(defaultPath)
	app.Before = Before
	app.Commands = CreateCommands()

	selfupdate.EnableLog()
	sintlsprovider.UserAgent = fmt.Sprintf("sintls/%s", version)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

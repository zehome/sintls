package main

import (
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/go-acme/lego/providers/dns"
	"github.com/go-pg/pg"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"io"
	"path"
	"strings"
)

import (
	sintls_dns "github.com/zehome/sintls/dns"
	"github.com/zehome/sintls/sintls"
)

func main() {
	bindaddress := flag.String(
		"bindaddress", "[::1]:8522",
		"sintls listening address")
	autotls := flag.Bool("autotls", false, "enable tls")
	dbuser := flag.String("dbuser", os.Getenv("USER"), "postgresql username")
	dbname := flag.String("dbname", os.Getenv("USER"), "postgresql dbname")
	dbhost := flag.String("dbhost", "/var/run/postgresql", "postgresql host")
	dbport := flag.Int("dbport", 5432, "postgresql port")
	providername := flag.String("provider", os.Getenv("PROVIDER"), "lego DNS provider name")
	logfile := flag.String(
		"logfile", path.Join(os.Getenv("HOME"), "sintls.log"),
		"sintls log")
	debug := flag.Bool("debug", false, "enable debug mode")
	quiet := flag.Bool("quiet", false, "quiet mode does not print anything on the console after initialization")
	initdb := flag.Bool("initdb", false, "initialize database")
	flag.Parse()
	log.SetOutput(os.Stdout)

	var dbaddr string
	var dbnetwork string
	if strings.HasPrefix(*dbhost, "/") {
		dbaddr = fmt.Sprintf("%s/.s.PGSQL.%d", *dbhost, *dbport)
		dbnetwork = "unix"
	} else {
		dbaddr = fmt.Sprintf("%s:%d", *dbhost, *dbport)
		dbnetwork = "tcp"
	}
	db, err := sintls.OpenDB(
		&pg.Options{
			Network:      dbnetwork,
			Addr:         dbaddr,
			User:         *dbuser,
			Database:     *dbname,
			PoolSize:     2,
			MinIdleConns: 0,
		},
		*debug,
		*initdb,
		true,
	)
	if err != nil {
		log.Fatal("OpenDB failed: ", err)
		return
	}

	// CLI
	if flag.NArg() >= 1 {
		RunCLI(db, flag.Args())
		return
	}

	// Server mode
	if len(*providername) == 0 {
		log.Fatal("-provider is mandatory")
		return
	}
	logger := &lumberjack.Logger{
	    Filename:   *logfile,
	    MaxSize:    10,
	    MaxBackups: 5,
	    Compress:   true,
	}
	if ! *quiet {
		log.SetOutput(io.MultiWriter(os.Stdout, logger))
	} else {
		log.SetOutput(logger)
	}

	if err != nil {
		log.Fatal("Unable to open database: ", err)
	}
	defer db.Close()
	if *initdb == true {
		return
	}

	// get lego acme provider
	provider, err := dns.NewDNSChallengeProviderByName(*providername)
	if err != nil {
		log.Fatal("invalid lego provider: ", err)
	}
	dnsupdater, err := sintls_dns.NewDNSUpdaterByName(*providername)
	if err != nil {
		log.Fatal("invalid dns updater: ", err)
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${remote_ip} - - [${time_rfc3339}] "${method} ${path} HTTP/1.1" ${status} ${bytes_out} "${referer}" "${user_agent}"` + "\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "pong"})
	})

	dbuse := e.Group("", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("database", db)
			c.Set("dnsprovider", provider)
			c.Set("dnsupdater", dnsupdater)
			return next(c)
		}
	})

	authorized := dbuse.Group("", middleware.BasicAuth(sintls.BasicAuth(db, false)))
	authorized.POST("/present", sintls.LegoPresent)
	authorized.POST("/cleanup", sintls.LegoCleanup)
	authorized.POST("/updatedns", sintls.UpdateDNSRecords)

	// admin := dbuse.Group("/admin", sintls.BasicAuth(db, true))
	// admin.POST("/auth", sintls.CreateAuth)
	// admin.DELETE("/auth", sintls.DeleteAuth)

	fmt.Printf(
		"Listening on %s (log: %s)\n", *bindaddress, *logfile)
	if *autotls == true {
		e.Logger.Fatal(e.StartAutoTLS(*bindaddress))
	} else {
		e.Logger.Fatal(e.Start(*bindaddress))
	}
}

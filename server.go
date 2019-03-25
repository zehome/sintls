package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-acme/lego/providers/dns"
	"github.com/go-pg/pg"
	"io"
	"log"
	"net/http"
	"os"
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
	dbuser := flag.String("dbuser", os.Getenv("USER"), "postgresql username")
	dbname := flag.String("dbname", os.Getenv("USER"), "postgresql dbname")
	dbhost := flag.String("dbhost", "/var/run/postgresql", "postgresql host")
	dbport := flag.Int("dbport", 5432, "postgresql port")
	providername := flag.String("provider", os.Getenv("PROVIDER"), "lego DNS provider name")
	logfile := flag.String(
		"logfile", path.Join(os.Getenv("HOME"), "sintls.log"),
		"sintls log")
	debug := flag.Bool("debug", false, "enable debug mode")
	initdb := flag.Bool("initdb", false, "initialize database")
	flag.Parse()

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

	// CLI
	if flag.NArg() >= 1 {
		RunCLI(db, flag.Args())
		return
	}

	// Server mode
	if err != nil {
		log.Fatal("Unable to open database", err)
	}
	defer db.Close()
	if *initdb == true {
		return
	}

	if *debug == false {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
		f, err := os.OpenFile(*logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			log.Fatal(err)
		}
		gin.DefaultWriter = io.MultiWriter(f)
		defer f.Close()
		log.SetOutput(f)
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

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	dbuse := r.Group("/", func(c *gin.Context) {
		c.Set("database", db)
		c.Set("dnsprovider", provider)
		c.Set("dnsupdater", dnsupdater)
	})

	authorized := dbuse.Group("/", sintls.BasicAuth(db, false))
	authorized.POST("/present", sintls.LegoPresent)
	authorized.POST("/cleanup", sintls.LegoCleanup)
	authorized.POST("/updatedns", sintls.UpdateDNSRecords)

	// admin := dbuse.Group("/admin", sintls.BasicAuth(db, true))
	// admin.POST("/auth", sintls.CreateAuth)
	// admin.DELETE("/auth", sintls.DeleteAuth)

	fmt.Printf(
		"Listening on %s (log: %s)\n", *bindaddress, *logfile)
	r.Run(*bindaddress)
}

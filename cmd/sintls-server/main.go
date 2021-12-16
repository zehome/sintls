package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/coreos/go-systemd/v22/activation"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

import (
	sintls_dns "github.com/zehome/sintls/dns"
	"github.com/zehome/sintls/sintls"
)

func GetStringDefaultEnv(variable, defaultvalue string) string {
	value := os.Getenv(variable)
	if len(value) == 0 {
		return defaultvalue
	} else {
		return value
	}
}

func GetBooleanDefaultEnv(variable string, defaultvalue bool) bool {
	value := os.Getenv(variable)
	if len(value) == 0 {
		return defaultvalue
	} else {
		v, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatal("invalid %s: %s", variable, err)
		}
		return v
	}
}

func GetIntDefaultEnv(variable string, defaultvalue int) int {
	value := os.Getenv(variable)
	if len(value) == 0 {
		return defaultvalue
	} else {
		v, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal("invalid %s: %s", variable, err)
		}
		return v
	}
}

func main() {
	bindaddress := flag.String(
		"bindaddress",
		GetStringDefaultEnv("SINTLS_BINDADDRESS", "[::1]:8522"),
		"sintls listening address")
	autotls := flag.Bool(
		"autotls",
		GetBooleanDefaultEnv("SINTLS_ENABLETLS", false),
		"enable tls")
	autotls_cachedir := flag.String(
		"autotls_cachedir",
		GetStringDefaultEnv("SINTLS_AUTOTLS_CACHEDIR", path.Join(os.Getenv("HOME"), ".cache")),
		"autotls cache directory")
	autotls_domain := flag.String(
		"autotls_domain",
		GetStringDefaultEnv("SINTLS_AUTOTLS_DOMAIN", ""),
		"autotls domain whitelist")
	dbuser := flag.String(
		"dbuser",
		GetStringDefaultEnv("SINTLS_DBUSER", os.Getenv("USER")),
		"postgresql username")
	dbname := flag.String(
		"dbname",
		GetStringDefaultEnv("SINTLS_DBNAME", os.Getenv("USER")),
		"postgresql dbname")
	dbhost := flag.String(
		"dbhost",
		GetStringDefaultEnv("SINTLS_DBHOST", "/var/run/postgresql"),
		"postgresql host")
	dbport := flag.Int(
		"dbport",
		GetIntDefaultEnv("SINTLS_DBPORT", 5432),
		"postgresql port")
	providername := flag.String(
		"provider",
		GetStringDefaultEnv("SINTLS_PROVIDER", "ovh"),
		"lego DNS provider name")
	debug := flag.Bool("debug", false, "enable debug mode")
	initdb := flag.Bool("initdb", false, "initialize database")
	disable_colors := flag.Bool("disable-colors", false, "disable colors")
	disable_autocleanup := flag.Bool("disable-autocleanup", false, "disable automatic DNS cleanup")
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
		RunCLI(db, *disable_colors, flag.Args())
		return
	}

	// Server mode
	if len(*providername) == 0 {
		log.Fatal("-provider is mandatory")
		return
	}
	logwriter, err := syslog.New(syslog.LOG_NOTICE, "sintls")
	// during initialization, we do provide informations on Stdout
	log.SetOutput(io.MultiWriter(os.Stdout, logwriter))

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
	e.Debug = *debug
	e.HideBanner = true
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

	// Ok, now we log everything using syslog
	if !*debug {
		log.SetOutput(logwriter)
	}

	// Background revoke worker
	if ! *disable_autocleanup {
		go autocleanup(db, dnsupdater, time.Hour * 6)
	}

	if *autotls {
		e.AutoTLSManager.Cache = autocert.DirCache(*autotls_cachedir)
		if len(*autotls_domain) != 0 {
			e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(*autotls_domain)
		}
		s := e.TLSServer
		s.TLSConfig = new(tls.Config)
		s.TLSConfig.GetCertificate = e.AutoTLSManager.GetCertificate
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, acme.ALPNProto)
		listeners, err := activation.TLSListeners(s.TLSConfig)
		if err == nil && len(listeners) == 1 {
			e.HidePort = true
			log.Println("Using systemd TLS listener")
			e.TLSListener = listeners[0]
			*bindaddress = e.TLSListener.Addr().String()
		}
		log.Printf("Listening on %s\n", *bindaddress)
		s.Addr = *bindaddress
		if !e.DisableHTTP2 {
			s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
		}
		e.Logger.Fatal(e.StartServer(e.TLSServer))
	} else {
		listeners, err := activation.Listeners()
		if err == nil && len(listeners) == 1 {
			e.HidePort = true
			log.Println("Using systemd HTTP listener")
			e.Listener = listeners[0]
			*bindaddress = e.Listener.Addr().String()
		}
		log.Printf("Listening on %s\n", *bindaddress)
		e.Logger.Fatal(e.Start(*bindaddress))
	}
}

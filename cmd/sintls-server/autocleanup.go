package main

import (
	"log"
	"time"
	"github.com/go-pg/pg/v10"
)

import "github.com/zehome/sintls/sintls"
import "github.com/zehome/sintls/dns"

// in days
const CLEANUP_TIMEOUT = 5
// in days (after which we remove from the database, even if we can't remove from DNS)
const CLEANUP_FORCETIMEOUT = 10

func autocleanup(db *pg.DB, dnsupdater dns.DNSUpdater, sleeptime time.Duration) {
	for {
		var hosts []sintls.Host
		err := db.Model(&hosts).Column(
			"host.host_id",
			"host.name",
			"host.updated_at").
			Order("host.name ASC").
			Where(`updated_at < (now() - '? hour'::interval)`, (24 * (90 + CLEANUP_TIMEOUT))).
			Select()
		if err != nil {
			log.Println(err)
		}
		for _, host := range hosts {
			log.Printf("Autocleanup %s (%d) (last update: %s)", host.Name, host.HostId, host.UpdatedAt)
			err := dnsupdater.RemoveRecords(host.Name)
			if err != nil {
				log.Printf("Autocleanup %s: unable to remove records: %v", host.Name, err)
			}
			if (err == nil || time.Now().Sub(host.UpdatedAt) > (time.Hour * 24 * (90 + CLEANUP_FORCETIMEOUT))) {
				_, err := db.Model(&host).WherePK().Delete()
				if err != nil {
					log.Println(err)
				}
			}
		}
		time.Sleep(sleeptime)
	}
}

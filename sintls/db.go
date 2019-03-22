package sintls

import (
	"database/sql"
	"fmt"
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"log"
	"net"
	"strings"
	"time"
)

type Authorization struct {
	tableName       struct{}     `sql:"sintls_authorization"`
	AuthorizationId uint64       `sql:"authorization_id,pk"`
	CreatedAt       time.Time    `sql:"created_at,notnull,default:now()"`
	UpdatedAt       time.Time    `sql:"updated_at,notnull,default:now()"`
	Name            string       `sql:"name,notnull,unique"`
	Secret          string       `sql:"secret,notnull"`
	Admin           sql.NullBool `sql:"admin,notnull,default:false"`
}

func (a Authorization) String() string {
	return fmt.Sprintf("Name: %s Secret: %s Id: %d", a.Name, a.Secret, a.AuthorizationId)
}

func (a *Authorization) CanUseHost(db *pg.DB, host string) bool {
	// Get subdomain
	subdomain := strings.Join(strings.Split(host, ".")[1:], ".")
	log.Println("host: ", host, "subdomain: ", subdomain, "a.Name: ", a.Name)
	var subdomains []SubDomain
	count, err := db.Model(&subdomains).
		Column("Authorization").
		Where(`"authorization".name = ?`, a.Name).
		Where(`"sub_domain".name = ?`, subdomain).
		Count()
	if err != nil {
		log.Println("Count error: ", err)
		return false
	}
	return count > 0
}

func (a *Authorization) CreateOrUpdateHost(
		db *pg.Tx, fqdn string,
		target_a, target_aaaa net.IP, target_cname string) error {
	// Find subdomain
	var subdomain *SubDomain
	var subdomains []SubDomain
	err := db.Model(&subdomains).
		Column("subdomain_id").
		Relation("Authorization").
		Where(`"authorization".authorization_id = ?`, a.AuthorizationId).
		Select()
	if err != nil {
		return fmt.Errorf("db: get subdomain failed: %s", err)
	}
	if len(subdomains) == 0 {
		return fmt.Errorf("db: no subdomain available")
	}
	for _, _subdomain := range subdomains {
		if strings.HasSuffix(fqdn, _subdomain.Name) {
			subdomain = &_subdomain
			break
		}
	}
	if subdomain == nil {
		return fmt.Errorf("db: hostname does not match any subdomain")
	}
	host := Host{
		Name: fqdn,
		SubDomainId: subdomain.SubDomainId,
		DnsTargetA: target_a,
		DnsTargetAAAA: target_aaaa,
		DnsTargetCNAME: target_cname,
	}
	qs := db.Model(&host).
		OnConflict("(name, subdomain_id) DO UPDATE").
		Set("updated_at = now()").
		Set("dns_target_cname = ?", target_cname)
	if target_a.String() != "<nil>" {
		qs = qs.Set("dns_target_a = ?", target_a)
	} else {
		qs = qs.Set("dns_target_a = ?", nil)
	}
	if target_aaaa.String() != "<nil>" {
		qs = qs.Set("dns_target_aaaa = ?", target_aaaa)
	} else {
		qs = qs.Set("dns_target_aaaa = ?", nil)
	}
	_, err = qs.Insert()
	if err != nil {
		return fmt.Errorf("db: update host failed: %s", err)
	}
	return nil
}

type SubDomain struct {
	tableName       struct{}  `sql:"sintls_subdomain"`
	SubDomainId     uint64    `sql:"subdomain_id,pk"`
	CreatedAt       time.Time `sql:"created_at,notnull,default:now()"`
	UpdatedAt       time.Time `sql:"updated_at,notnull,default:now()"`
	Authorization   *Authorization
	AuthorizationId uint64 `sql:"authorization_id,notnull,on_delete:CASCADE"`
	Name            string `sql:"name,notnull,unique"`
}

type Host struct {
	tableName     struct{}  `sql:"sintls_host"`
	HostId        uint64    `sql:"host_id,pk"`
	CreatedAt     time.Time `sql:"created_at,notnull,default:now()"`
	UpdatedAt     time.Time `sql:"updated_at,notnull,default:now()"`
	SubDomain     *SubDomain
	SubDomainId   uint64 `sql:"subdomain_id,notnull,on_delete:CASCADE"`
	Name          string `sql:"name,notnull"`
	DnsTargetA    net.IP `sql:"dns_target_a"`
	DnsTargetAAAA net.IP `sql:"dns_target_aaaa"`
	DnsTargetCNAME string `sql:"dns_target_cname"`
}


type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	log.Println(q.FormattedQuery())
}

func OpenDB(options *pg.Options, debug bool, initdb bool, runmigrations bool) (db *pg.DB, err error) {
	db = pg.Connect(options)
	if debug {
		db.AddQueryHook(dbLogger{})
	}
	if initdb {
		_, _, err := migrations.Run(db, "init")
		if err != nil {
			return db, err
		}
	}
	if runmigrations {
		oldVersion, newVersion, err := migrations.Run(db, "up")
		if newVersion != oldVersion {
			fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
		}
		if err != nil {
			return db, err
		}
	}
	return db, err
}

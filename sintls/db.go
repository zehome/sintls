package sintls

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"log"
	"net"
	"strings"
	"time"
)

type Authorization struct {
	tableName       struct{}     `pg:"sintls_authorization"`
	AuthorizationId uint64       `pg:"authorization_id,pk"`
	CreatedAt       time.Time    `pg:"created_at,notnull,default:now()"`
	UpdatedAt       time.Time    `pg:"updated_at,notnull,default:now()"`
	Name            string       `pg:"name,notnull,unique"`
	Secret          string       `pg:"secret,notnull"`
	Admin           sql.NullBool `pg:"admin,notnull,default:false"`
}

func (a Authorization) String() string {
	return fmt.Sprintf("Name: %s Secret: %s Id: %d", a.Name, a.Secret, a.AuthorizationId)
}

func (a *Authorization) CanUseHost(db *pg.DB, host string) bool {
	// We will try to get down to the root subdomain by splitting and stripping the first part
	// until we can find a valid subdomain... Or not!
	namesplit := strings.Split(host, ".")
	for stripindex := range namesplit {
		if (stripindex > 5) {
			log.Println("host:", host, "a.Name:", a.Name, "strip was going too far (>5)")
			break
		}
		subdomain := strings.Join(namesplit[stripindex:], ".")
		log.Println("host:", host, "subdomain:", subdomain, "a.Name:", a.Name)
		var subdomains []SubDomain
		count, err := db.Model(&subdomains).
			Relation("Authorization").
			Where(`"authorization".name = ?`, a.Name).
			Where(`"sub_domain".name = ?`, subdomain).
			Count()
		if err != nil {
			log.Println("Count error:", err)
			return false
		}
		if (count > 0) {
			return true
		}
	}
	return false
}

func (a *Authorization) CreateOrUpdateHost(
	db *pg.Tx, fqdn string,
	target_a, target_aaaa net.IP, target_cname string, target_mx string) error {
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
		Name:           fqdn,
		SubDomainId:    subdomain.SubDomainId,
		DnsTargetA:     target_a,
		DnsTargetAAAA:  target_aaaa,
		DnsTargetCNAME: target_cname,
		DnsTargetMX:	target_mx,
	}
	qs := db.Model(&host).
		OnConflict("(name, subdomain_id) DO UPDATE").
		Set("updated_at = now()").
		Set("dns_target_cname = ?", target_cname).
		Set("dns_target_mx = ?", target_mx)
	if len(target_a) != 0 {
		qs = qs.Set("dns_target_a = ?", target_a)
	} else {
		qs = qs.Set("dns_target_a = ?", nil)
	}
	if len(target_aaaa) != 0 {
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
	tableName       struct{}  `pg:"sintls_subdomain"`
	SubDomainId     uint64    `pg:"subdomain_id,pk"`
	CreatedAt       time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt       time.Time `pg:"updated_at,notnull,default:now()"`
	Authorization   *Authorization
	AuthorizationId uint64 `pg:"authorization_id,notnull,on_delete:CASCADE"`
	Name            string `pg:"name,notnull,unique"`
}

type Host struct {
	tableName      struct{}  `pg:"sintls_host"`
	HostId         uint64    `pg:"host_id,pk"`
	CreatedAt      time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt      time.Time `pg:"updated_at,notnull,default:now()"`
	SubDomain      *SubDomain
	SubDomainId    uint64 `pg:"subdomain_id,notnull,on_delete:CASCADE"`
	Name           string `pg:"name,notnull"`
	DnsTargetA     net.IP `pg:"dns_target_a"`
	DnsTargetAAAA  net.IP `pg:"dns_target_aaaa"`
	DnsTargetCNAME string `pg:"dns_target_cname"`
	DnsTargetMX    string `pg:"dns_target_mx"`
}

type dbLogger struct {
}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	log.Println(q.FormattedQuery())
	return nil
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

package sintls

import (
	"database/sql"
	"fmt"
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"log"
	"net"
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

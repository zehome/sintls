package sintls

import (
	"github.com/go-pg/migrations/v7"
)

// This is just to avoid a bug in go-pg/migrations
// in order to detect only .sql files
func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`
CREATE TABLE sintls_authorization (
  authorization_id bigserial primary key,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  name text unique not null,
  secret text not null,
  admin boolean not null default false
);
CREATE INDEX ON sintls_authorization(name);

CREATE TABLE sintls_subdomain (
  subdomain_id bigserial primary key,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  name text unique not null,
  authorization_id bigint not null references sintls_authorization(authorization_id) on delete cascade
);
CREATE INDEX ON sintls_subdomain(name);
CREATE INDEX ON sintls_subdomain(authorization_id);

CREATE TABLE sintls_host (
  host_id bigserial primary key,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  name text not null,
  subdomain_id bigint not null references sintls_subdomain(subdomain_id) on delete cascade,
  dns_target_a inet,
  dns_target_aaaa inet,
  UNIQUE (name, subdomain_id),
  CHECK (dns_target_a is not null or dns_target_aaaa is not null)
);
CREATE INDEX on sintls_host(name);
CREATE INDEX on sintls_host(subdomain_id);`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`DROP TABLE sintls_host; DROP TABLE sintls_subdomain; DROP TABLE sintls_authorization;`)
		return err
	})
}

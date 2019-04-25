package sintls

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`
ALTER TABLE sintls_host ADD COLUMN dns_target_cname text,
DROP CONSTRAINT sintls_host_check,
ADD CONSTRAINT sintls_host_check CHECK (dns_target_a is not null or dns_target_aaaa is not null or dns_target_cname is not null);`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`
ALTER TABLE sintls_host
DROP CONSTRAINT sintls_host_check,
ADD CONSTRAINT sintls_host_check CHECK (dns_target_a is not null or dns_target_aaaa is not null),
DROP COLUMN sintls_target_cname;`)
		return err
	})
}
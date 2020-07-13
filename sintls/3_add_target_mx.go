package sintls

import (
	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`
ALTER TABLE sintls_host ADD COLUMN dns_target_mx text;`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`
ALTER TABLE sintls_host
DROP COLUMN sintls_target_cname;`)
		return err
	})
}

package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		db.ExecContext(ctx, `
			ALTER TABLE resources 
			ADD COLUMN gitleaks jsonb DEFAULT '[]'::jsonb;
		`)

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		return nil
	})
}

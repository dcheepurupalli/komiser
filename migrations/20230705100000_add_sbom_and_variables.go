package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		db.ExecContext(ctx, `
			ALTER TABLE resources
			ADD COLUMN variables jsonb default '[]'::jsonb;
		`)
		db.ExecContext(ctx, `
			ALTER TABLE resources
			ADD COLUMN sbom jsonb default '{}'::jsonb;
		`)
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		return nil
	})
}

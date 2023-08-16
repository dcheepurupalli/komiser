package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		db.ExecContext(ctx, `
			CREATE TABLE nodes (
				id uuid NOT NULL DEFAULT gen_random_uuid(),
				name varchar(255) NOT NULL,
				"type" varchar(255) NOT NULL,
				fields json NOT NULL,
				created_at timestamptz NOT NULL DEFAULT now(),
				updated_at timestamptz NOT NULL DEFAULT now(),
				CONSTRAINT pkey_nodes PRIMARY KEY (id, type)
			);
		`)
		db.ExecContext(ctx, `
			CREATE TABLE edges (
				source bigint NOT NULL ,
				dest bigint NOT NULL,
				name VARCHAR(255) NOT NULL,
				CONSTRAINT pkey_edges PRIMARY KEY (source, dest, name)
			);
		`)
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		// db.ExecContext(ctx, `
		// 		DROP TABLE edges;
		// 	`)
		// db.ExecContext(ctx, `
		// 		DROP TABLE nodes;
		// 	`)
		return nil
	})
}

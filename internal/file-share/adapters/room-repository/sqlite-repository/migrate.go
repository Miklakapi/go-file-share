package sqliterepository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MakeMigrations(ctx context.Context, conn *sql.DB, migrationsDir string) error {
	if err := createSchemaMigrations(ctx, conn); err != nil {
		return err
	}

	em, err := getExecutedMigrations(ctx, conn)
	if err != nil {
		return err
	}

	m, err := getNotExecutedMigrations(ctx, migrationsDir, em)
	if err != nil {
		return err
	}

	for _, f := range m {
		if err := executeMigration(ctx, conn, migrationsDir, f); err != nil {
			return err
		}
	}
	return nil
}

func getExecutedMigrations(ctx context.Context, conn *sql.DB) (map[string]struct{}, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	appliedSet := make(map[string]struct{}, 10)

	rows, err := conn.QueryContext(ctx, "SELECT id FROM schema_migrations ORDER BY id;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		appliedSet[id] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appliedSet, nil
}

func getNotExecutedMigrations(ctx context.Context, migrationsDir string, executedMigrations map[string]struct{}) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	migrations := make([]string, 0, 5)

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if _, ok := executedMigrations[name]; ok {
			continue
		}

		if filepath.Ext(name) != ".sql" {
			continue
		}

		migrations = append(migrations, name)
	}

	return migrations, nil
}

func executeMigration(ctx context.Context, conn *sql.DB, migrationsDir string, migration string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	fullPath := filepath.Join(migrationsDir, migration)

	sqlBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	migrationSQL := strings.TrimSpace(string(sqlBytes))
	if migrationSQL == "" {
		return nil
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, migrationSQL); err != nil {
		return fmt.Errorf("migration %s failed: %w", migration, err)
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (id) VALUES (?);", migration); err != nil {
		return fmt.Errorf("migration %s: insert schema_migrations failed: %w", migration, err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func createSchemaMigrations(ctx context.Context, conn *sql.DB) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id         TEXT PRIMARY KEY,
			applied_at INTEGER NOT NULL DEFAULT (CAST(strftime('%s','now') AS INTEGER))
		);
	`)
	return err
}

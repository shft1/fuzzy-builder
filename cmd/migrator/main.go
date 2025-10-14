package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/fuzzy?sslmode=disable")
	migrationsDir := getenv("MIGRATIONS_DIR", "migrations")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	if err := ensureSchemaMigrations(ctx, pool); err != nil {
		log.Fatalf("ensure table: %v", err)
	}

	files, err := readMigrationFiles(migrationsDir)
	if err != nil {
		log.Fatalf("read migrations: %v", err)
	}
	if len(files) == 0 {
		log.Printf("no migrations found in %s", migrationsDir)
		return
	}

	applied, err := loadApplied(ctx, pool)
	if err != nil {
		log.Fatalf("load applied: %v", err)
	}

	for _, f := range files {
		if applied[f] {
			continue
		}
		log.Printf("applying %s", f)
		sqlBytes, err := os.ReadFile(filepath.Join(migrationsDir, f))
		if err != nil {
			log.Fatalf("read %s: %v", f, err)
		}
		if err := applyMigration(ctx, pool, f, string(sqlBytes)); err != nil {
			log.Fatalf("apply %s: %v", f, err)
		}
		log.Printf("applied %s", f)
	}
	log.Printf("migrations complete")
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func ensureSchemaMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
        version TEXT PRIMARY KEY,
        applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    )`)
	return err
}

func readMigrationFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errorsIs(err, fs.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		files = append(files, name)
	}
	sort.Strings(files)
	return files, nil
}

func loadApplied(ctx context.Context, pool *pgxpool.Pool) (map[string]bool, error) {
	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out[v] = true
	}
	return out, rows.Err()
}

func applyMigration(ctx context.Context, pool *pgxpool.Pool, version, sql string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, sql); err != nil {
		return fmt.Errorf("exec migration: %w", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// errorsIs avoids importing errors in main body earlier
func errorsIs(err, target error) bool {
	return err != nil && target != nil && strings.Contains(err.Error(), target.Error())
}

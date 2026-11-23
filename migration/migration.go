package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Migrate(db *sql.DB, dir string) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating schema_migrations table: %w", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading migration directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		var exists bool
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)`, file).Scan(&exists)
		if err != nil {
			return fmt.Errorf("checking migration status for %s: %w", file, err)
		}
		if exists {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return fmt.Errorf("reading migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("executing migration %s: %w", file, err)
		}

		if _, err := db.Exec(`INSERT INTO schema_migrations (filename) VALUES ($1)`, file); err != nil {
			return fmt.Errorf("recording migration %s: %w", file, err)
		}
	}

	return nil
}
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	migrationsDir = "migrations"
)

func main() {
	// Parse command line flags
	direction := flag.String("direction", "up", "Migration direction: up or down")
	flag.Parse()

	log.Println("🔄 Starting database migrations...")

	// Get database connection string from environment
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "superuser"),
		getEnv("DB_PASSWORD", "superpass"),
		getEnv("DB_NAME", "super_salary_db"),
		getEnv("DB_SSLMODE", "disable"),
	)

	// Connect to database
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database")

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("❌ Failed to create migrations table: %v", err)
	}

	// Run migrations
	if *direction == "up" {
		if err := migrateUp(db); err != nil {
			log.Fatalf("❌ Migration up failed: %v", err)
		}
		log.Println("✅ Migrations completed successfully")
	} else if *direction == "down" {
		if err := migrateDown(db); err != nil {
			log.Fatalf("❌ Migration down failed: %v", err)
		}
		log.Println("✅ Rollback completed successfully")
	} else {
		log.Fatalf("❌ Invalid direction: %s (must be 'up' or 'down')", *direction)
	}
}

// createMigrationsTable creates the schema_migrations table to track applied migrations
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	_, err := db.Exec(query)
	return err
}

// migrateUp applies pending migrations
func migrateUp(db *sql.DB) error {
	// Get all .up.sql migration files
	files, err := getMigrationFiles(".up.sql")
	if err != nil {
		return err
	}

	if len(files) == 0 {
		log.Println("ℹ️  No migration files found")
		return nil
	}

	// Get already applied migrations
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	// Apply pending migrations
	for _, file := range files {
		version := getVersionFromFilename(file)
		if applied[version] {
			log.Printf("⏭️  Skipping already applied migration: %s", file)
			continue
		}

		log.Printf("📝 Applying migration: %s", file)
		if err := applyMigration(db, file, version); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}
		log.Printf("✅ Successfully applied: %s", file)
	}

	return nil
}

// migrateDown rolls back the last applied migration
func migrateDown(db *sql.DB) error {
	// Get all .down.sql migration files
	files, err := getMigrationFiles(".down.sql")
	if err != nil {
		return err
	}

	if len(files) == 0 {
		log.Println("ℹ️  No migration files found")
		return nil
	}

	// Get already applied migrations
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	// Reverse the order for down migrations (rollback from newest to oldest)
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	// Rollback the last applied migration
	for _, file := range files {
		version := getVersionFromFilename(file)
		if !applied[version] {
			log.Printf("⏭️  Skipping non-applied migration: %s", file)
			continue
		}

		log.Printf("🔙 Rolling back migration: %s", file)
		if err := rollbackMigration(db, file, version); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", file, err)
		}
		log.Printf("✅ Successfully rolled back: %s", file)
		
		// Only rollback one migration at a time
		break
	}

	return nil
}

// getMigrationFiles returns all migration files with the given suffix, sorted
func getMigrationFiles(suffix string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)
	return files, nil
}

// getAppliedMigrations returns a map of already applied migration versions
func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// applyMigration executes an up migration and records it
func applyMigration(db *sql.DB, filename, version string) error {
	// Read migration file
	content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	// Record migration in schema_migrations table
	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return err
	}

	return tx.Commit()
}

// rollbackMigration executes a down migration and removes the record
func rollbackMigration(db *sql.DB, filename, version string) error {
	// Read migration file
	content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute rollback SQL
	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	// Remove migration record from schema_migrations table
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
		return err
	}

	return tx.Commit()
}

// getVersionFromFilename extracts the version number from a migration filename
// Example: "001_create_users_table.up.sql" -> "001"
func getVersionFromFilename(filename string) string {
	parts := strings.Split(filename, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return filename
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

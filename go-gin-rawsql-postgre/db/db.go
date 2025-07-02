package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection and creates the database if it doesn't exist
func InitDB(dsn string) {
	// Parse the DSN to get database name
	dbName := extractDBName(dsn)

	// Connect to postgres database to create our target database
	postgresDSN := strings.Replace(dsn, "/"+dbName, "/postgres", 1)
	postgresDB, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatal("Error opening postgres DB:", err)
	}
	defer postgresDB.Close()

	// Create database if it doesn't exist
	createDatabaseIfNotExists(postgresDB, dbName)

	// Connect to the target database
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	// Run schema if tables don't exist
	runSchemaIfNeeded(DB)

	log.Println("Database initialized successfully")
}

// extractDBName extracts the database name from the DSN
func extractDBName(dsn string) string {
	// Parse DSN like: "postgres://username:password@localhost:5432/dbname?sslmode=disable"
	parts := strings.Split(dsn, "/")
	if len(parts) >= 3 {
		dbPart := parts[len(parts)-1]
		// Remove query parameters if any
		if strings.Contains(dbPart, "?") {
			dbPart = strings.Split(dbPart, "?")[0]
		}
		return dbPart
	}
	return "film_rental" // default fallback
}

// createDatabaseIfNotExists creates the database if it doesn't exist
func createDatabaseIfNotExists(db *sql.DB, dbName string) {
	// Check if database exists
	var exists int
	query := `SELECT 1 FROM pg_database WHERE datname = $1`
	err := db.QueryRow(query, dbName).Scan(&exists)

	if err == sql.ErrNoRows {
		// Database doesn't exist, create it
		createQuery := fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)
		_, err = db.Exec(createQuery)
		if err != nil {
			log.Fatalf("Error creating database %s: %v", dbName, err)
		}
		log.Printf("Database %s created successfully", dbName)
	} else if err != nil {
		log.Fatalf("Error checking database existence: %v", err)
	} else {
		log.Printf("Database %s already exists", dbName)
	}
}

// runSchemaIfNeeded runs the schema if tables don't exist
func runSchemaIfNeeded(db *sql.DB) {
	// Check if any of our tables exist
	var count int
	query := `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name IN ('staff', 'film', 'customer')`
	err := db.QueryRow(query).Scan(&count)

	if err != nil {
		log.Fatalf("Error checking tables: %v", err)
	}

	if count == 0 {
		log.Println("No tables found, running schema...")
		runSchema(db)
	} else {
		log.Println("Tables already exist, skipping schema creation")
	}
}

// runSchema executes the schema.sql file
func runSchema(db *sql.DB) {
	// Use the original schema file with improved parsing
	runSchemaManually(db)

	// Optionally run seed data
	runSeedData(db)
}

// runSchemaManually executes the schema file using psql command
func runSchemaManually(db *sql.DB) {
	// Use psql command to execute the schema file properly
	// This is the recommended way to handle PostgreSQL dump files
	schemaPath := "schema/schema.sql"

	// Get database connection info from the DSN
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Parse DSN to get connection details
	host, port, user, password, dbname := parseDSN(dsn)

	// Set environment variables for psql
	env := os.Environ()
	env = append(env, fmt.Sprintf("PGHOST=%s", host))
	env = append(env, fmt.Sprintf("PGPORT=%s", port))
	env = append(env, fmt.Sprintf("PGUSER=%s", user))
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", password))
	env = append(env, fmt.Sprintf("PGDATABASE=%s", dbname))

	// Execute psql command
	cmd := exec.Command("psql", "-f", schemaPath)
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Warning: Error executing schema with psql: %v", err)
		log.Printf("Output: %s", string(output))

		// Fallback to manual parsing if psql fails
		log.Println("Falling back to manual parsing...")
		runSchemaManuallyFallback(db)
	} else {
		log.Println("Schema executed successfully with psql")
	}
}

// parseDSN parses the database connection string
func parseDSN(dsn string) (host, port, user, password, dbname string) {
	// Parse DSN like: "postgres://username:password@localhost:5432/dbname?sslmode=disable"

	// Remove the postgres:// prefix
	dsn = strings.TrimPrefix(dsn, "postgres://")

	// Split by @ to separate credentials from host
	parts := strings.Split(dsn, "@")
	if len(parts) != 2 {
		log.Printf("Warning: Invalid DSN format: %s", dsn)
		return "localhost", "5432", "postgres", "", "postgres"
	}

	// Parse credentials
	creds := strings.Split(parts[0], ":")
	if len(creds) >= 2 {
		user = creds[0]
		password = creds[1]
	}

	// Parse host and database
	hostPart := parts[1]
	hostPortDB := strings.Split(hostPart, "/")
	if len(hostPortDB) >= 2 {
		hostPort := strings.Split(hostPortDB[0], ":")
		if len(hostPort) >= 2 {
			host = hostPort[0]
			port = hostPort[1]
		} else {
			host = hostPort[0]
			port = "5432"
		}

		dbname = strings.Split(hostPortDB[1], "?")[0] // Remove query parameters
	}

	return
}

// runSchemaManuallyFallback is a fallback method for schema execution
func runSchemaManuallyFallback(db *sql.DB) {
	// Read schema file
	schemaPath := "schema/schema.sql"
	schemaContent, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Error reading schema file: %v", err)
	}

	// Use improved PostgreSQL dump parser
	statements := parsePostgresDumpImproved(string(schemaContent))

	// Execute each statement
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		_, err := db.Exec(statement)
		if err != nil {
			log.Printf("Warning: Error executing statement: %v", err)
			log.Printf("Statement: %s", statement[:min(len(statement), 100)])
		}
	}

	log.Println("Schema executed successfully manually")
}

// parsePostgresDumpImproved parses a PostgreSQL dump file with better handling of dollar-quoted strings
func parsePostgresDumpImproved(content string) []string {
	var statements []string
	lines := strings.Split(content, "\n")
	var currentStatement strings.Builder
	inStatement := false
	inComment := false
	inDollarQuote := false
	dollarTag := ""

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines
		if trimmedLine == "" {
			continue
		}

		// Handle comment blocks
		if strings.HasPrefix(trimmedLine, "/*") {
			inComment = true
			continue
		}
		if strings.Contains(trimmedLine, "*/") {
			inComment = false
			continue
		}
		if inComment {
			continue
		}

		// Skip single line comments
		if strings.HasPrefix(trimmedLine, "--") {
			continue
		}

		// Skip metadata lines
		if isMetadataLine(trimmedLine) {
			continue
		}

		// Handle dollar-quoted strings
		if !inDollarQuote {
			// Look for start of dollar quote
			if strings.Contains(trimmedLine, "$$") {
				dollarTag = "$$"
				inDollarQuote = true
			} else if strings.Contains(trimmedLine, "$_$") {
				dollarTag = "$_$"
				inDollarQuote = true
			}
		} else {
			// Look for end of dollar quote
			if strings.Contains(trimmedLine, dollarTag) {
				inDollarQuote = false
				dollarTag = ""
			}
		}

		// If we're in a dollar quote, just add the line and continue
		if inDollarQuote {
			if inStatement {
				currentStatement.WriteString(" ")
				currentStatement.WriteString(trimmedLine)
			}
			continue
		}

		// Check if this line starts a new statement
		if strings.HasPrefix(strings.ToUpper(trimmedLine), "CREATE") ||
			strings.HasPrefix(strings.ToUpper(trimmedLine), "INSERT") ||
			strings.HasPrefix(strings.ToUpper(trimmedLine), "ALTER") ||
			strings.HasPrefix(strings.ToUpper(trimmedLine), "DROP") ||
			strings.HasPrefix(strings.ToUpper(trimmedLine), "GRANT") ||
			strings.HasPrefix(strings.ToUpper(trimmedLine), "REVOKE") ||
			strings.HasPrefix(strings.ToUpper(trimmedLine), "SET") {

			// If we were building a statement, save it
			if inStatement {
				stmt := currentStatement.String()
				if strings.TrimSpace(stmt) != "" {
					statements = append(statements, stmt)
				}
				currentStatement.Reset()
			}

			inStatement = true
			currentStatement.WriteString(trimmedLine)
		} else if inStatement {
			// Continue building the current statement
			currentStatement.WriteString(" ")
			currentStatement.WriteString(trimmedLine)
		}

		// Check if statement ends with semicolon
		if inStatement && strings.HasSuffix(trimmedLine, ";") {
			stmt := currentStatement.String()
			if strings.TrimSpace(stmt) != "" {
				statements = append(statements, stmt)
			}
			currentStatement.Reset()
			inStatement = false
		}
	}

	// Add any remaining statement
	if inStatement {
		stmt := currentStatement.String()
		if strings.TrimSpace(stmt) != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}

// isMetadataLine checks if a line is a metadata line from pg_dump
func isMetadataLine(line string) bool {
	// Skip lines that are just metadata comments or empty
	if strings.TrimSpace(line) == "" {
		return true
	}

	// Skip lines that are just comments
	if strings.HasPrefix(strings.TrimSpace(line), "--") {
		return true
	}

	// Skip lines that contain only metadata keywords
	line = strings.TrimSpace(line)
	metadataKeywords := []string{
		"Type:", "Schema:", "Owner:", "Name:", "ALTER TYPE", "ALTER DOMAIN",
		"SET statement_timeout", "SET lock_timeout", "SET idle_in_transaction_session_timeout",
		"SET client_encoding", "SET standard_conforming_strings", "SELECT pg_catalog.set_config",
		"SET check_function_bodies", "SET xmloption", "SET client_min_messages", "SET row_security",
		"PostgreSQL database dump", "Dumped from database version", "Dumped by pg_dump version",
	}

	for _, keyword := range metadataKeywords {
		if strings.Contains(line, keyword) {
			return true
		}
	}

	// Skip lines that are just metadata separators
	if line == "--" {
		return true
	}

	return false
}

// runSeedData runs the seed data if needed
func runSeedData(db *sql.DB) {
	// Check if admin user already exists
	var count int
	query := `SELECT COUNT(*) FROM staff WHERE username = 'admin'`
	err := db.QueryRow(query).Scan(&count)

	if err != nil {
		log.Printf("Warning: Error checking admin user: %v", err)
		return
	}

	if count == 0 {
		log.Println("No admin user found, running seed data...")
		runSeedScript(db)
	} else {
		log.Println("Admin user already exists, skipping seed data")
	}
}

// runSeedScript executes the seed_admin.sql file
func runSeedScript(db *sql.DB) {
	seedPath := "schema/seed_admin.sql"
	seedContent, err := os.ReadFile(seedPath)
	if err != nil {
		log.Printf("Warning: Error reading seed file: %v", err)
		return
	}

	// Split the seed into individual statements
	statements := strings.Split(string(seedContent), ";")

	// Execute each statement
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		_, err := db.Exec(statement)
		if err != nil {
			log.Printf("Warning: Error executing seed statement: %v", err)
		}
	}

	log.Println("Seed data executed successfully")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

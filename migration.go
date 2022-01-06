package pge

var (
	// "1" is the number we arbitrarily choose for the schema migrations lock
	LockMigrations = NewQuery("acquire advisory lock for schema migrations", `
		SELECT pg_advisory_xact_lock(1)
	`)

	CreateTableSchemaVersion = NewQuery("create table schema_versions", `
		CREATE TABLE IF NOT EXISTS schema_versions (
			version int,
			migrated timestamptz NOT NULL DEFAULT NOW(),
			PRIMARY KEY (version, migrated)
		)
	`)

	SelectLatestSchemaVersion = NewQuery("select latest schema version", `
		SELECT version
		FROM schema_versions
		ORDER BY migrated DESC
		LIMIT 1
	`)

	InsertSchemaVersion = NewQuery("insert schema version", `
		INSERT INTO schema_versions(version)
		VALUES ($1)
	`)
)

type Migration struct {
	Name    string
	Queries []Query
}

package config

// DB Config for a single DB (SQL, opensearch)
type DBConfig struct {
	// IP/domain
	Host string

	// Port number
	Port uint16

	// Creds
	Username string
	Password string

	// Database name
	Database string
}

package config

type LogLevelT uint8

const (
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_ERROR
)

func (ll LogLevelT) IsValid() bool {
	return ll <= LOG_LEVEL_ERROR
}

// Server config, contains everything that can be modified
type ServerConfig struct {
	// Server port
	Port uint16

	// Logs
	LogLevel LogLevelT
}

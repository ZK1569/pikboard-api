package util

type config struct {
	Port     string
	Env      string
	Version  string
	Database databaseConfig
}

type databaseConfig struct {
	Username string
	Name     string
	Password string
	Host     string
	Port     string
}

var EnvVariable *config

func init() {
	EnvVariable = &config{
		Port:    getString("PORT", "8080"),
		Env:     getString("ENV", "development"),
		Version: getString("VERSION", "0.0.1"),
		Database: databaseConfig{
			Username: getString("DB_USER", ""),
			Name:     getString("DB_NAME", ""),
			Password: getString("DB_PASSWORD", ""),
			Host:     getString("DB_HOST", ""),
			Port:     getString("DB_PORT", ""),
		},
	}

}

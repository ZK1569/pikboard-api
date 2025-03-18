package util

type config struct {
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

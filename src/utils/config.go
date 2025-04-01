package util

type config struct {
	Env      string
	Version  string
	Database databaseConfig
	S3       s3Config
	Ai       aiConfig
}

type databaseConfig struct {
	Username string
	Name     string
	Password string
	Host     string
	Port     string
}

type s3Config struct {
	AccessKey       string
	SecretAccessKey string
	Region          string
	BucketName      string
}

type aiConfig struct {
	OpenAiKey string
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
		S3: s3Config{
			AccessKey:       getString("AWS_ACCESS_KEY", ""),
			SecretAccessKey: getString("AWS_SECRET_ACCESS_KEY", ""),
			Region:          getString("AWS_REGION", ""),
			BucketName:      getString("AWS_BUCKET_NAME", ""),
		},
		Ai: aiConfig{
			OpenAiKey: getString("OPENAI_KEY", ""),
		},
	}

}

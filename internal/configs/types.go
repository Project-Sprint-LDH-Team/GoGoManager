package configs

type (
	Config struct {
		Service  Service  `mapstructure:"service"`
		Database Database `mapstructure:"database"`
		AWS      AWSConfig
	}

	Service struct {
		Port      string `mapstructure:"port"`
		SecretJWT string `mapstructure:"secretJWT"`
	}

	Database struct {
		DataSourceName string `mapstructure:"dataSourceName"`
	}

	AWSConfig struct {
		Region          string
		Bucket          string
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
	}
)

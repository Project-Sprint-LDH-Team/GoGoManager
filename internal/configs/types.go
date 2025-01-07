package configs

type (
	Config struct {
		Service       Service
		Database      DatabaseConfig
		SpotifyConfig SpotifyConfig
	}

	Service struct {
		Port      string
		SecretKey string
	}

	DatabaseConfig struct {
		DataSourceName string
	}

	SpotifyConfig struct {
		ClientID     string
		ClientSecret string
	}
)

package configs

import (
	"github.com/spf13/viper"
)

var (
	config *Config
)

// option defines configuration option.
type option struct {
	configFolders []string
	configFile    string
	configType    string
}

// Init initializes `config` from the default config file.
// use `WithConfigFile` to specify the location of the config file.
func Init(opts ...Option) error {
	opt := &option{
		configFolders: getDefaultConfigFolder(),
		configFile:    getDefaultConfigFile(),
		configType:    getDefaultConfigType(),
	}

	for _, optFunc := range opts {
		optFunc(opt)
	}
	for _, configFolder := range opt.configFolders {
		viper.AddConfigPath(configFolder)
	}

	viper.SetConfigName(opt.configFile)
	viper.SetConfigType(opt.configType)
	viper.AutomaticEnv()

	config = new(Config)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&config)
}

// Option define an option for config package.
type Option func(*option)

// WithConfigFolder set `configs` to use the given configs folder.
func WithConfigFolder(configFolders []string) Option {
	return func(opt *option) {
		opt.configFolders = configFolders
	}
}

// WithConfigFile set `config` to use the given config file.
func WithConfigFile(configFile string) Option {
	return func(opt *option) {
		opt.configFile = configFile
	}
}

// WithConfigType set `config` to use the given config type.
func WithConfigType(configType string) Option {
	return func(opt *option) {
		opt.configType = configType
	}
}

// getDefaultConfigFolder get default config folder.
func getDefaultConfigFolder() []string {
	return []string{"./configs/"}
}

// getDefaultConfigFile get default config file.
func getDefaultConfigFile() string {
	return "config"
}

// getDefaultConfigType get default config type.
func getDefaultConfigType() string {
	return "yaml"
}

// Get config.
func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}

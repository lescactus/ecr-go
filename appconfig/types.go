package appconfig

// Config stores configuration values
var Config *config

type config struct {

	// Application provides the application configuration
	Application struct {
		Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
		ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
		LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
		DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
		Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
	}
}

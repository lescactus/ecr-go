package appconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidLogLevel(t *testing.T) {
	tests := []struct {
		desc  string
		input string
		want  bool
	}{
		{
			desc:  "Log level set to error",
			input: "error",
			want:  true,
		},
		{
			desc:  "Log level set to info",
			input: "info",
			want:  true,
		},
		{
			desc:  "Log level set to debug",
			input: "debug",
			want:  true,
		},
		{
			desc:  "Log level set to invalid value",
			input: "invalid",
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			b := isValidLogLevel(test.input)
			assert.Equal(t, test.want, b)
		})
	}
}

func TestLoadConfig(t *testing.T) {
	testsWithoutError := []struct {
		desc  string
		osEnv map[string]string
		input *config
		want  *config
	}{
		{
			desc: "Override all environment variables",
			osEnv: map[string]string{
				"APPLICATION_NAME":    "foo",
				"CONFIG_DIR":          "dir/",
				"LOG_LEVEL":           "error",
				"DRY_RUN":             "true",
				"APPLICATION_VERSION": "99.99.99",
			},
			input: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{},
			},
			want: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{
					Name:      "foo",
					ConfigDir: "dir/",
					LogLevel:  "error",
					DryRun:    true,
					Version:   "99.99.99",
				},
			},
		},
		{
			desc:  "No environment variables override",
			osEnv: map[string]string{},
			input: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{},
			},
			want: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{
					Name:      "ecr-go",
					ConfigDir: "files/",
					LogLevel:  "info",
					DryRun:    false,
					Version:   "0.1.2",
				},
			},
		},
		{
			desc: "Override log level",
			osEnv: map[string]string{
				"LOG_LEVEL": "debug",
			},
			input: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{},
			},
			want: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{
					Name:      "ecr-go",
					ConfigDir: "files/",
					LogLevel:  "debug",
					DryRun:    false,
					Version:   "0.1.2",
				},
			},
		},
		{
			desc: "Override dry run mode",
			osEnv: map[string]string{
				"DRY_RUN": "true",
			},
			input: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{},
			},
			want: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{
					Name:      "ecr-go",
					ConfigDir: "files/",
					LogLevel:  "info",
					DryRun:    true,
					Version:   "0.1.2",
				},
			},
		},
	}

	testsWithErrors := []struct {
		desc  string
		osEnv map[string]string
		input *config
		want  *config
	}{
		{
			desc: "Incorrect loglevel",
			osEnv: map[string]string{
				"LOG_LEVEL": "invalid",
			},
			input: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{},
			},
			want: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{
					Name:      "foo",
					ConfigDir: "dir/",
					LogLevel:  "error",
					Version:   "99.99.99",
				},
			},
		},
		{
			desc: "Incorrect dry run mode",
			osEnv: map[string]string{
				"DRY_RUN": "invalid",
			},
			input: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{},
			},
			want: &config{
				Application: struct {
					Name      string `env:"APPLICATION_NAME" envDefault:"ecr-go"`
					ConfigDir string `env:"CONFIG_DIR" envDefault:"files/"`
					LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
					DryRun    bool   `env:"DRY_RUN" envDefault:"false"`
					Version   string `env:"APPLICATION_VERSION" envDefault:"0.1.2"`
				}{
					Name:      "foo",
					ConfigDir: "dir/",
					LogLevel:  "error",
					DryRun:    false,
					Version:   "99.99.99",
				},
			},
		},
	}

	for _, test := range testsWithoutError {
		t.Run(test.desc, func(t *testing.T) {
			// set environment variables
			for k, v := range test.osEnv {
				os.Setenv(k, v)
			}
			LoadConfig(test.input)
			assert.Equal(t, test.want, test.input)

			// reset environment variables
			for k := range test.osEnv {
				os.Unsetenv(k)
			}
		})
	}

	for _, test := range testsWithErrors {
		t.Run(test.desc, func(t *testing.T) {
			// set environment variables
			for k, v := range test.osEnv {
				os.Setenv(k, v)
			}
			err := LoadConfig(test.input)
			assert.Error(t, err)

			// reset environment variables
			for k := range test.osEnv {
				os.Unsetenv(k)
			}
		})
	}
}

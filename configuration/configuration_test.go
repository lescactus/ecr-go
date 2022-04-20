package configuration

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	Logger, _ = zap.NewProduction()
	defer Logger.Sync()
}

func TestNewConfigurationFile(t *testing.T) {
	var c1 ConfigurationFile
	var c2 = ConfigurationFile{
		logger: Logger,
	}

	c1 = NewConfigurationFile(Logger)

	assert.Equal(t, c1, c2)
}

func TestGetYamlConfigurationFiles(t *testing.T) {
	tests := []struct {
		desc         string
		mockFilesDir string
		want         []string
	}{
		{
			desc:         "Yaml files exists in an existing directory",
			mockFilesDir: "testdata/files/",
			want:         []string{"testdata/files/test_1.yaml", "testdata/files/test_10.yaml", "testdata/files/test_11.yaml", "testdata/files/test_2.yml", "testdata/files/test_5.yaml", "testdata/files/test_6.yaml", "testdata/files/test_7.yaml", "testdata/files/test_8.yaml", "testdata/files/test_9.yaml"},
		},
		{
			desc:         "Yaml files doesn't exists in an existing directory",
			mockFilesDir: "testdata/no_files/",
			want:         []string{},
		},
		{
			desc:         "Directory doesn't exists",
			mockFilesDir: "nothing/",
			want:         []string(nil),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			a, _ := GetYamlConfigurationFiles(test.mockFilesDir)
			assert.Equal(t, test.want, a)
		})
	}
}

func TestLoadYamlConfiguration(t *testing.T) {
	policy_1_json, _ := ioutil.ReadFile("testdata/files/policies/policy_1.json")

	testsWithoutError := []struct {
		desc     string
		mockFile string
		want     ConfigurationFile
	}{
		{
			desc:     "Yaml file exists, policy exists, policy is valid",
			mockFile: "testdata/files/test_1.yaml",
			want: ConfigurationFile{
				RepositoryName:       "repository_1",
				RepositoryPolicyFile: "testdata/files/policies/policy_1.json",
				RepositoryPolicy:     policy_1_json,
				logger:               Logger,
			},
		},
	}

	testsWithError := []struct {
		desc     string
		mockFile string
		want     error
	}{
		{
			desc:     "Yaml file exists, policy exists, policy is invalid",
			mockFile: "testdata/files/test_2.yml",
			want:     errors.New("not a valid json file"),
		},
		{
			desc:     "Yaml file exists, policy does not exists",
			mockFile: "testdata/files/test_5.yaml",
			want:     errors.New("open policydoesnotexists.json: no such file or directory"),
		},
		{
			desc:     "Yaml file doesn't exists",
			mockFile: "testdata/files/filedoesnotexists.yaml",
			want:     errors.New("open testdata/files/filedoesnotexists.yaml: no such file or directory"),
		},
		{
			desc:     "Yaml file is invalid (1)",
			mockFile: "testdata/files/test_6.yaml",
			want:     errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `this is...` into configuration.ConfigurationFile"),
		},
		{
			desc:     "Yaml file is invalid (2)",
			mockFile: "testdata/files/test_7.yaml",
			want:     errors.New("yaml: unmarshal errors:\n  line 3: field nonExistingField not found in type configuration.ConfigurationFile"),
		},
		{
			desc:     "Yaml file exists, repositoryName is missing",
			mockFile: "testdata/files/test_8.yaml",
			want:     errors.New("RepositoryName must be present and not empty"),
		},
		{
			desc:     "Yaml file exists, repositoryName is empty",
			mockFile: "testdata/files/test_9.yaml",
			want:     errors.New("RepositoryName must be present and not empty"),
		},
		{
			desc:     "Yaml file exists, RepositoryPolicyFile is missing",
			mockFile: "testdata/files/test_10.yaml",
			want:     errors.New("RepositoryPolicyFile must be present and not empty"),
		},
		{
			desc:     "Yaml file exists, RepositoryPolicyFile is empty",
			mockFile: "testdata/files/test_11.yaml",
			want:     errors.New("RepositoryPolicyFile must be present and not empty"),
		},
	}

	for _, test := range testsWithoutError {
		t.Run(test.desc, func(t *testing.T) {
			a := NewConfigurationFile(Logger)
			err := a.LoadYamlConfiguration(test.mockFile)
			assert.NoError(t, err)
			assert.Equal(t, test.want, a)
		})
	}

	for _, test := range testsWithError {
		t.Run(test.desc, func(t *testing.T) {
			a := NewConfigurationFile(Logger)
			err := a.LoadYamlConfiguration(test.mockFile)
			assert.Error(t, test.want, err)
			if assert.Errorf(t, err, test.want.Error()) {
				assert.EqualError(t, err, test.want.Error())
			}
		})
	}
}

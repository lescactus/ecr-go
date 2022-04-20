package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type ConfigurationFile struct {
	RepositoryName       string `yaml:"repositoryName"`
	RepositoryPolicyFile string `yaml:"repositoryPolicyFile"`
	RepositoryPolicy     []byte
	logger               *zap.Logger
}

// GetYamlConfigurationFiles will recursively look for all yaml files in the root directory passed as argument
// Only the files ending with .yaml or .yml will be accepted
// It returns a list of all yaml files found or any error encountered
func GetYamlConfigurationFiles(root string) ([]string, error) {
	yamlConfigurationFiles := []string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, e error) error {
		// In case of any error, return
		if e != nil {
			return e
		}
		// Only look for .yaml or .yml files
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			yamlConfigurationFiles = append(yamlConfigurationFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return yamlConfigurationFiles, nil
}

// NewConfigurationFile instanciate a ConfigurationFile struct
// It returns an empty ConfigurationFile
func NewConfigurationFile(logger *zap.Logger) ConfigurationFile {
	return ConfigurationFile{
		logger: logger,
	}
}

// LoadYamlConfiguration will load the yaml file into a ConfigurationFile struct
// Additionally, it will load the json policy defined in RepositoryPolicyFile
// It returns any error encountered
func (c *ConfigurationFile) LoadYamlConfiguration(yamlFile string) error {
	c.logger.Debug(fmt.Sprintf("%s - Reading yaml file ...", yamlFile))
	d, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	c.logger.Debug(fmt.Sprintf("%s - Unmarshalling ...", yamlFile))
	err = yaml.UnmarshalStrict(d, &c)
	if err != nil {
		return err
	}
	// Ensure RepositoryName and RepositoryPolicyFile are not empty
	if c.RepositoryName == "" {
		return errors.New("RepositoryName must be present and not empty")
	}
	if c.RepositoryPolicyFile == "" {
		return errors.New("RepositoryPolicyFile must be present and not empty")
	}

	c.logger.Debug(fmt.Sprintf("%s - Unmarshalled: [%v, %v]", yamlFile, c.RepositoryName, c.RepositoryPolicyFile))

	// Validate the RepositoryPolicyFile is a valid json file
	c.logger.Debug(fmt.Sprintf("%s - Validating json policy: %s ...", yamlFile, c.RepositoryPolicyFile))
	j, err := ioutil.ReadFile(c.RepositoryPolicyFile)
	if err != nil {
		return err
	}
	if !json.Valid(j) {
		return errors.New("not a valid json file")
	}
	c.logger.Debug(fmt.Sprintf("%s - Json policy validated: %s ...", yamlFile, c.RepositoryPolicyFile))

	// Load the associated json policy from the file defined in RepositoryPolicyFile
	c.RepositoryPolicy = j

	return nil
}

package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/lescactus/ecr-go/appconfig"
	"github.com/lescactus/ecr-go/configuration"
	"github.com/lescactus/ecr-go/ecrupdater"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func main() {

	logger, err := NewLogger(appconfig.Config.Application.LogLevel)
	if err != nil {
		log.Fatalln("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info(fmt.Sprintf("Staring %s v%s", appconfig.Config.Application.Name, appconfig.Config.Application.Version))
	logger.Info(fmt.Sprintf("Configuration directory is set to %s", appconfig.Config.Application.ConfigDir))
	logger.Info(fmt.Sprintf("Running in dry-mode: %v", appconfig.Config.Application.DryRun))

	var ConfigurationFiles []configuration.ConfigurationFile

	// Look recursively for all yaml configuration files
	yamlConfigurationFilesList, err := configuration.GetYamlConfigurationFiles(appconfig.Config.Application.ConfigDir)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error: cannot get the configuration files: %v", err))
	}

	// For each yaml configuration file, load the associated json policy defined in ConfigurationFile.RepositoryPolicyFile
	// in ConfigurationFile.RepositoryPolicy
	// Ensure there is no duplicates
	for _, yamlFile := range yamlConfigurationFilesList {
		c := configuration.NewConfigurationFile(logger)
		if err := c.LoadYamlConfiguration(yamlFile); err == nil {
			for _, y := range ConfigurationFiles {
				if c.RepositoryName == y.RepositoryName {
					logger.Fatal(fmt.Sprint("Error: Duplicate RepositoryName", c.RepositoryName, "found in", yamlFile))
				}
			}
			ConfigurationFiles = append(ConfigurationFiles, c)
		} else {
			logger.Fatal(fmt.Sprintf("Error: Loading %s: %v", yamlFile, err))
		}
	}

	// Instanciate a new aws session
	awssession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	e := ecrupdater.ECRUpdaterClient{
		Client: ecr.New(awssession, &aws.Config{}),
		Logger: logger,
	}
	e.Init()

	// Skip the ECR update if in dry run mode
	if !appconfig.Config.Application.DryRun {
		var wg sync.WaitGroup

		// Update the ECR repositories policies
		for i := range ConfigurationFiles {
			wg.Add(1)
			go e.Work(ConfigurationFiles[i], &wg)
		}

		wg.Wait()

		// Summarize how it went
		logger.Info("")
		logger.Info("Repository update completed. Summary:")
		logger.Info(fmt.Sprintf("\tNumber of successful repositories updates: %v", len(e.RepositorySuccededUpdate.RepositoryNames)))
		for i := range e.RepositorySuccededUpdate.RepositoryNames {
			logger.Info(fmt.Sprintf("\t\t- %v", e.RepositorySuccededUpdate.RepositoryNames[i]))
		}
		logger.Info(fmt.Sprintf("\tNumber of failed repositories updates: %v", len(e.RepositoryFailedUpdate.GetAll())))
		for i := range e.RepositoryFailedUpdate.GetAll() {
			logger.Info(fmt.Sprintf("\t\t- %v: %v", i, e.RepositoryFailedUpdate.GetAll()[i]))
		}

		if len(e.RepositoryFailedUpdate.GetAll()) > 0 {
			os.Exit(1)
		}
	}
	if appconfig.Config.Application.DryRun {
		logger.Info("Dry-run completed ... all configuration files are valid")
	}
}

package ecrupdater

import (
	"fmt"
	"sync"

	"github.com/lescactus/ecr-go/configuration"
	"github.com/lescactus/ecr-go/summary"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

type ECRUpdaterClient struct {
	Client                   ecriface.ECRAPI
	RepositoryFailedUpdate   summary.RepositoryFailedUpdate
	RepositorySuccededUpdate summary.RepositorySuccededUpdate
	Logger                   *zap.Logger
}

// Init will initialize the ECR client
func (e *ECRUpdaterClient) Init() {
	e.RepositoryFailedUpdate = summary.NewRepositoryFailedUpdate()
	e.RepositorySuccededUpdate = summary.NewRepositorySuccededUpdate()
}

// Work will update the given ECR repository policy
// It will update the status of the update (success or fail) in a summary.RepositoryFailedUpdate and a summary.RepositorySuccededUpdate
func (e *ECRUpdaterClient) Work(config configuration.ConfigurationFile, wg *sync.WaitGroup) {
	defer wg.Done()

	e.Logger.Info(fmt.Sprintf("Updating repository %s ...", config.RepositoryName))

	// Actual AWS call to update the ECR repository policy
	_, err := e.Client.SetRepositoryPolicy(&ecr.SetRepositoryPolicyInput{
		PolicyText:     aws.String(string(config.RepositoryPolicy)),
		RepositoryName: &config.RepositoryName,
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			e.Logger.Error(fmt.Sprintf("Error: An error occured while updating the repository %v: \"%v\"", config.RepositoryName, awsErr))

		} else {
			e.Logger.Error(fmt.Sprintf("Error: An error occured while updating the repository %v: \"%v\"", config.RepositoryName, awsErr))
		}
		e.RepositoryFailedUpdate.Add(config.RepositoryName, err)
	} else {
		e.Logger.Info(fmt.Sprintf("Policy updated for repository %s", config.RepositoryName))
		e.RepositorySuccededUpdate.RepositoryNames = append(e.RepositorySuccededUpdate.RepositoryNames, config.RepositoryName)
	}
}

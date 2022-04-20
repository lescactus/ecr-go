package ecrupdater

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/lescactus/ecr-go/configuration"
	"github.com/lescactus/ecr-go/summary"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/stretchr/testify/assert"
)

var Logger *zap.Logger

type mockedECRUpdatedPolicy struct {
	ecriface.ECRAPI
	Output ecr.SetRepositoryPolicyOutput
}

func (m mockedECRUpdatedPolicy) SetRepositoryPolicy(input *ecr.SetRepositoryPolicyInput) (*ecr.SetRepositoryPolicyOutput, error) {
	switch strings.Split(aws.StringValue(input.RepositoryName), "_")[0] {
	case "notfound":
		return &ecr.SetRepositoryPolicyOutput{}, awserr.New("RepositoryNotFoundException", "The repository with name "+aws.StringValue(input.RepositoryName)+" does not exist in the registry", errors.New("RepositoryNotFoundException"))
	case "expiredtoken":
		return &ecr.SetRepositoryPolicyOutput{}, awserr.New("ExpiredTokenException", "The security token included in the request is expired", errors.New("ExpiredTokenException"))
	case "nocredentials":
		return &ecr.SetRepositoryPolicyOutput{}, awserr.New("NoCredentialProviders", "no valid providers in chain. Deprecated.", errors.New("NoCredentialProviders"))
	case "accessdenied":
		return &ecr.SetRepositoryPolicyOutput{}, awserr.New("AccessDeniedException", "User is not authorized to perform", errors.New("AccessDeniedException"))
	case "generic":
		return &ecr.SetRepositoryPolicyOutput{}, errors.New("Generic error")
	}
	return &ecr.SetRepositoryPolicyOutput{
		RepositoryName: aws.String("foo"),
		RegistryId:     aws.String("foo"),
		PolicyText: aws.String(`{
			"Version": "2008-10-17",
			"Statement": [
				{
					"Sid": "CrossAccount",
					"Effect": "Allow",
					"Principal": {
						"AWS": [
							"arn:aws:iam::123456789123:root",
						]
					},
					"Action": [
						"ecr:GetDownloadUrlForLayer",
					]
				}
			]
		}`),
	}, nil
}

func init() {
	cfg := zap.Config{
		Encoding: "console",
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
	}
	var err error
	Logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync()
}

func TestWork(t *testing.T) {
	var validRepositoryPolicy = []byte(`{
		"Version": "2008-10-17",
		"Statement": [
			{
				"Sid": "CrossAccount",
				"Effect": "Allow",
				"Principal": {
					"AWS": [
						"arn:aws:iam::123456789123:root",
					]
				},
				"Action": [
					"ecr:GetDownloadUrlForLayer",
				]
			}
		]
	}`)

	testsWithoutError := []struct {
		desc string
		cf   configuration.ConfigurationFile
		want ECRUpdaterClient
	}{
		{
			desc: "ECR repository updated successfully",
			cf: configuration.ConfigurationFile{
				RepositoryName:   "foo",
				RepositoryPolicy: validRepositoryPolicy,
			},
			want: ECRUpdaterClient{
				RepositoryFailedUpdate: summary.NewRepositoryFailedUpdate(),
				RepositorySuccededUpdate: summary.RepositorySuccededUpdate{
					RepositoryNames: []string{"foo"},
				},
				Logger: Logger,
			},
		},
	}

	testsWithErrors := []struct {
		desc string
		cf   configuration.ConfigurationFile
		want ECRUpdaterClient
	}{
		{
			desc: "ECR repository not found error",
			cf: configuration.ConfigurationFile{
				RepositoryName:   "notfound_foo",
				RepositoryPolicy: validRepositoryPolicy,
			},
			want: ECRUpdaterClient{
				RepositoryFailedUpdate: summary.RepositoryFailedUpdate{
					ErrorRepositoryName: map[string]error{
						"notfound_foo": awserr.New("RepositoryNotFoundException", "The repository with name notfound_foo does not exist in the registry", errors.New("RepositoryNotFoundException")),
					},
				},
				RepositorySuccededUpdate: summary.NewRepositorySuccededUpdate(),
			},
		},
		{
			desc: "Expired AWS token",
			cf: configuration.ConfigurationFile{
				RepositoryName:   "expiredtoken_foo",
				RepositoryPolicy: validRepositoryPolicy,
			},
			want: ECRUpdaterClient{
				RepositoryFailedUpdate: summary.RepositoryFailedUpdate{
					ErrorRepositoryName: map[string]error{
						"expiredtoken_foo": awserr.New("ExpiredTokenException", "The security token included in the request is expired", errors.New("ExpiredTokenException")),
					},
				},
				RepositorySuccededUpdate: summary.NewRepositorySuccededUpdate(),
			},
		},
		{
			desc: "No AWS credentials providers",
			cf: configuration.ConfigurationFile{
				RepositoryName:   "nocredentials_foo",
				RepositoryPolicy: validRepositoryPolicy,
			},
			want: ECRUpdaterClient{
				RepositoryFailedUpdate: summary.RepositoryFailedUpdate{
					ErrorRepositoryName: map[string]error{
						"nocredentials_foo": awserr.New("NoCredentialProviders", "no valid providers in chain. Deprecated.", errors.New("NoCredentialProviders")),
					},
				},
				RepositorySuccededUpdate: summary.NewRepositorySuccededUpdate(),
			},
		},
		{
			desc: "Access denied",
			cf: configuration.ConfigurationFile{
				RepositoryName:   "accessdenied_foo",
				RepositoryPolicy: validRepositoryPolicy,
			},
			want: ECRUpdaterClient{
				RepositoryFailedUpdate: summary.RepositoryFailedUpdate{
					ErrorRepositoryName: map[string]error{
						"accessdenied_foo": awserr.New("AccessDeniedException", "User is not authorized to perform", errors.New("AccessDeniedException")),
					},
				},
				RepositorySuccededUpdate: summary.NewRepositorySuccededUpdate(),
			},
		},
		{
			desc: "Generic error",
			cf: configuration.ConfigurationFile{
				RepositoryName:   "generic_foo",
				RepositoryPolicy: validRepositoryPolicy,
			},
			want: ECRUpdaterClient{
				RepositoryFailedUpdate: summary.RepositoryFailedUpdate{
					ErrorRepositoryName: map[string]error{
						"generic_foo": errors.New("Generic error"),
					},
				},
				RepositorySuccededUpdate: summary.NewRepositorySuccededUpdate(),
			},
		},
	}

	for _, test := range testsWithoutError {
		t.Run(test.desc, func(t *testing.T) {
			e := ECRUpdaterClient{
				Client: mockedECRUpdatedPolicy{
					Output: ecr.SetRepositoryPolicyOutput{},
				},
				Logger: Logger,
			}
			e.Init()

			var wg sync.WaitGroup
			wg.Add(1)
			go e.Work(test.cf, &wg)
			wg.Wait()

			assert := assert.New(t)

			assert.NoError(e.RepositoryFailedUpdate.Get("foo"))
			assert.Nil(e.RepositoryFailedUpdate.Get("foo"))
			assert.EqualValues(test.want.RepositorySuccededUpdate, e.RepositorySuccededUpdate)
			assert.EqualValues(test.want.RepositoryFailedUpdate, e.RepositoryFailedUpdate)
		})
	}

	for _, test := range testsWithErrors {
		t.Run(test.desc, func(t *testing.T) {
			e := ECRUpdaterClient{
				Client: mockedECRUpdatedPolicy{
					Output: ecr.SetRepositoryPolicyOutput{},
				},
				Logger: Logger,
			}
			e.Init()

			var wg sync.WaitGroup
			wg.Add(1)
			go e.Work(test.cf, &wg)
			wg.Wait()

			assert := assert.New(t)
			assert.Error(e.RepositoryFailedUpdate.Get(test.cf.RepositoryName))
			assert.NotNil(e.RepositoryFailedUpdate.Get(test.cf.RepositoryName))
			assert.EqualValues(test.want.RepositoryFailedUpdate.Get(test.cf.RepositoryName), e.RepositoryFailedUpdate.Get(test.cf.RepositoryName))
			assert.EqualValues(test.want.RepositorySuccededUpdate, e.RepositorySuccededUpdate)
			assert.EqualValues(test.want.RepositoryFailedUpdate, e.RepositoryFailedUpdate)
			assert.Empty(e.RepositorySuccededUpdate.RepositoryNames)

		})
	}
}

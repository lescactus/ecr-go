# ecr-go

This repository contains a simple cli written in go to update an AWS ECR repository with a given policy.

## Usage

`ecr-go` works by reading all YAML config file(s) (`.yaml` or `.yml`), by default in the `files/` directory that define a repository and its associated policy.

Example of valid config file:

`myrepo.yaml`:
```yaml
repositoryName: alma # this is the name of the repository we want to apply the policy to
repositoryPolicyFile: policy.json # this is the policy to apply
```

`policy.json`:
```json
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Sid": "CrossAccountPull",
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::000000000000:root",
                    "arn:aws:iam::111111111111:root",
                    "arn:aws:iam::222222222222:root"
                ]
            },
            "Action": [
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability"
            ]
        }
    ]
}
```

### Usage with docker

If you've build the docker image embedding this cli, you need to pass your aws keys or credentials file to the container:

```sh
# Use environment variables
docker run --rm \
  -it \
  -e AWS_ACCESS_KEY_ID=xxxx \
  -e AWS_SECRET_ACCESS_KEY=xxxx \
  -e AWS_DEFAULT_REGION=xxxx \
  -v "<path to the config directory>:/app/files/"
  ecr-go

# Mount your config & credentials files in the container
docker run \
  --rm \
  -it \
  -v ~/.aws:/root/.aws \
  -v "<path to the config directory>:/app/files/"
  -e AWS_PROFILE=<your profile if needed> \
  ecr-go
```

### Configuration

`ecr-go` is looking for the following environment variables:

| Name     | Type | Default value    | Description |
| --------|---------|---------|-------|
| `APPLICATION_NAME`  | `string` | `ecr-go`   | Name of the application   |
| `CONFIG_DIR` | `string` |`files/` | Directory where `ecr-go` will look for config files   |
| `DRY_RUN` | `bool` |`false` | Enable dry run mode. Accepted values are go `bool` values: `1|0`, `t|f`, `T|F`, `true|false`, `TRUE|FALSE`, `True|False`|
| `LOG_LEVEL` | `string` |`info` | Verbosity level. Accepted values are `error`, `info` (default) and `debug`    |
| `APPLICATION_VERSION` | `string` |`0.0.2` | Version of `ecr-go`    |

#### Dry Run mode

Running in Dry Run mode will on verify that the yaml files are valid. It will not modify the ECR repository policies.

### Examples

#### Simple example
Given the following file tree:
```sh
$ tree 
.
├── ecr-go*
└─── files
     ├── alma-keel.json
     └── alma-keel.yaml

$ cat files/alma-keel.yaml
repositoryName: alma-keel
repositoryPolicyFile: files/alma-keel.json

$ cat files/alma-keel.json
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Sid": "CrossAccountPull",
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::000000000000:root",
                    "arn:aws:iam::111111111111:root",
                    "arn:aws:iam::222222222222:root"
                ]
            },
            "Action": [
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability"
            ]
        }
    ]
}       
```

Simply run:
```sh
$ ./ecr-go
2021-05-04T23:06:59+02:00	info	Staring ecr-go v0.1.0
2021-05-04T23:06:59+02:00	info	Configuration directory is set to files/
2021-05-04T23:06:59+02:00	info	Running in dry-mode: false
2021-05-04T23:06:59+02:00	info	Updating repository alma-keel ...
2021-05-04T23:06:59+02:00	info	Policy updated for repository alma-keel
2021-05-04T23:06:59+02:00	info	
2021-05-04T23:06:59+02:00	info	Repository update completed. Summary:
2021-05-04T23:06:59+02:00	info		Number of successful repositories updates: 1
2021-05-04T23:06:59+02:00	info			- alma-keel
2021-05-04T23:06:59+02:00	info		Number of failed repositories updates: 0
```

You can have several yaml files in the `files/` directory.

#### Example with failed policies update

In case of mistake in the configuration (repository name inexistant or insufficient permissions for instance), `ecr-go` will summarize:
```sh
$ ./ecr-go
2021-05-04T23:08:23+02:00	info	Staring ecr-go v0.1.0
2021-05-04T23:08:23+02:00	info	Configuration directory is set to files/
2021-05-04T23:08:23+02:00	info	Running in dry-mode: false
2021-05-04T23:08:23+02:00	info	Updating repository alma-keel ...
2021-05-04T23:08:23+02:00	info	Updating repository alma-2 ...
2021-05-04T23:08:23+02:00	info	Updating repository alma-11 ...
2021-05-04T23:08:23+02:00	info	Updating repository alma-0 ...
2021-05-04T23:08:23+02:00	info	Updating repository alma-1 ...
2021-05-04T23:08:23+02:00	error	Error: An error occured while updating the repository alma-2: "RepositoryNotFoundException: The repository with name 'alma-2' does not exist in the registry with id '000000000000'"
2021-05-04T23:08:23+02:00	error	Error: An error occured while updating the repository alma-0: "RepositoryNotFoundException: The repository with name 'alma-0' does not exist in the registry with id '000000000000'"
2021-05-04T23:08:23+02:00	error	Error: An error occured while updating the repository alma-11: "RepositoryNotFoundException: The repository with name 'alma-11' does not exist in the registry with id '000000000000'"
2021-05-04T23:08:23+02:00	info	Policy updated for repository alma-keel
2021-05-04T23:08:23+02:00	error	Error: An error occured while updating the repository alma-1: "RepositoryNotFoundException: The repository with name 'alma-1' does not exist in the registry with id '000000000000'"
2021-05-04T23:08:23+02:00	info	
2021-05-04T23:08:23+02:00	info	Repository update completed. Summary:
2021-05-04T23:08:23+02:00	info		Number of successful repositories updates: 1
2021-05-04T23:08:23+02:00	info			- alma-keel
2021-05-04T23:08:23+02:00	info		Number of failed repositories updates: 4
2021-05-04T23:08:23+02:00	info			- alma-11: RepositoryNotFoundException: The repository with name 'alma-11' does not exist in the registry with id '000000000000'
2021-05-04T23:08:23+02:00	info			- alma-1: RepositoryNotFoundException: The repository with name 'alma-1' does not exist in the registry with id '000000000000'
2021-05-04T23:08:23+02:00	info			- alma-2: RepositoryNotFoundException: The repository with name 'alma-2' does not exist in the registry with id '000000000000'
2021-05-04T23:08:23+02:00	info			- alma-0: RepositoryNotFoundException: The repository with name 'alma-0' does not exist in the registry with id '000000000000'
exit status 1
```

## Installation

### From source with go

You need a working [go](https://golang.org/doc/install) toolchain (It has been developped and tested with go 1.14 and go 1.16 only, but should work with go >= 1.12 ). Refer to the official documentation for more information (or from your Linux/Mac/Windows distribution documentation to install it from your favorite package manager).

```sh
# Clone this repository
git clone https://github.com/lescactus/ecr-go.git && cd ecr-go/

# Build from sources. Use the '-o' flag to change the compiled binary name
go build

# Default compiled binary is ecr-go
# You can optionnaly move it somewhere in your $PATH to access it shell wide
./ecr-go
```

### From source with docker

If you don't have [go](https://golang.org/) installed but have docker, run the following command to build inside a docker container:

```sh
# Build from sources inside a docker container. Use the '-o' flag to change the compiled binary name
# Warning: the compiled binary belongs to root:root
docker run --rm -it -v "$PWD":/app -w /app golang:1.14 go build

# Default compiled binary is ecr-go
# You can optionnaly move it somewhere in your $PATH to access it shell wide
./ecr-go
```

### From source with docker but built inside a docker image

If you don't want to pollute your computer with another program, with cli comes with its own docker image:

```sh
docker build -t ecr-go .
```

## Testing

To run the test suite, run the following commands:

```sh
# Run the unit tests. Remove the '-v' flag to reduce verbosity
go test -v ./... 

# Get coverage to html format
go test -coverprofile -v /tmp/cover.out ./...
go tool cover -html=/tmp/cover.out -o /tmp/cover.out.html
```
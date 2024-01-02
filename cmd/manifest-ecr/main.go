package main

import (
	"encoding/base64"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/drone-plugins/drone-manifest/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const defaultRegion = "us-east-1"

func main() {
	// Load env-file if it exists first
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		err := godotenv.Load(env)
		if err != nil {
			panic(err)
		}
	}

	var (
		registry   = getEnv("PLUGIN_REGISTRY")
		spec       = getEnv("PLUGIN_SPEC")
		region     = getEnv("PLUGIN_REGION", "ECR_REGION", "AWS_REGION")
		key        = getEnv("PLUGIN_ACCESS_KEY", "ECR_ACCESS_KEY", "AWS_ACCESS_KEY_ID")
		secret     = getEnv("PLUGIN_SECRET_KEY", "ECR_SECRET_KEY", "AWS_SECRET_ACCESS_KEY")
		assumeRole = getEnv("PLUGIN_ASSUME_ROLE")
		externalId = getEnv("PLUGIN_EXTERNAL_ID")
	)

	// set the region
	if region == "" {
		region = defaultRegion
	}

	setEnvOrPanic("AWS_REGION", region)

	if key != "" && secret != "" {
		setEnvOrPanic("AWS_ACCESS_KEY_ID", key)
		setEnvOrPanic("AWS_SECRET_ACCESS_KEY", secret)
	}

	sess, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		log.Fatalf("error creating aws session: %v", err)
	}

	svc := getECRClient(sess, assumeRole, externalId)
	username, password, defaultRegistry, err := getAuthInfo(svc)

	if registry == "" {
		registry = defaultRegistry
	}

	if err != nil {
		log.Fatalf("error getting ECR auth: %v", err)
	}

	setEnvOrPanic("PLUGIN_REGISTRY", registry)
	setEnvOrPanic("DOCKER_USERNAME", username)
	setEnvOrPanic("DOCKER_PASSWORD", password)
	setEnvOrPanic("PLUGIN_SPEC", spec)

	// invoke the base docker plugin binary
	cmd := exec.Command(util.GetDroneManifestExecCmd())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logrus.Fatal(err)
	}
}

func getAuthInfo(svc *ecr.ECR) (username, password, registry string, err error) {
	var result *ecr.GetAuthorizationTokenOutput
	var decoded []byte

	result, err = svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return
	}

	auth := result.AuthorizationData[0]
	token := *auth.AuthorizationToken
	decoded, err = base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}

	registry = strings.TrimPrefix(*auth.ProxyEndpoint, "https://")
	creds := strings.Split(string(decoded), ":")
	username = creds[0]
	password = creds[1]
	return
}

// func parseBoolOrDefault(defaultValue bool, s string) (result bool) {
// 	var err error
// 	result, err = strconv.ParseBool(s)
// 	if err != nil {
// 		result = false
// 	}
//
// 	return
// }

func getEnv(key ...string) (s string) {
	for _, k := range key {
		s = os.Getenv(k)
		if s != "" {
			return
		}
	}
	return
}

func setEnvOrPanic(key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		panic(err)
	}
}

func getECRClient(sess *session.Session, role string, externalId string) *ecr.ECR {
	if role == "" {
		return ecr.New(sess)
	}
	if externalId != "" {
		return ecr.New(sess, &aws.Config{
			Credentials: stscreds.NewCredentials(sess, role, func(p *stscreds.AssumeRoleProvider) {
				p.ExternalID = &externalId
			}),
		})
	} else {
		return ecr.New(sess, &aws.Config{
			Credentials: stscreds.NewCredentials(sess, role),
		})
	}
}

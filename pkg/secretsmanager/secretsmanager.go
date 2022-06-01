package secretsmanager

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/patrickmn/go-cache"
	"gitlab.com/pos_malaysia/golib/env"
	"gitlab.com/pos_malaysia/golib/logs"
)

// NOTE : The fields below are mapped according to the secret's JSON string in AWS Secrets Manager
type AWSSecretsManagerResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// For unit testing purpose. See https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/
type SecretsManagerAPI interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

var (
	localCache *cache.Cache
	asmResp    AWSSecretsManagerResponse
	secretName = env.Get("AWS_SECRET_NAME")
)

func init() {
	// Set to no expiration since we will not be pulling data from
	// AWS Secrets Manager from time to time.
	localCache = cache.New(cache.NoExpiration, cache.NoExpiration)
}

// GetClientIDAndSecret returns the client_id and client_secret from AWS Secrets Manager.
//
// * NOTE : Secrets will be cached locally to reduce the number of requests to AWS Secrets Manager.
func GetClientIDAndSecret(ctx context.Context, api SecretsManagerAPI) (AWSSecretsManagerResponse, error) {

	// check if these items are available in local memory cache
	client_id, client_id_found := localCache.Get("client_id")
	client_secret, client_secret_found := localCache.Get("client_secret")

	if !client_id_found || !client_secret_found {

		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(secretName),
			VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
		}

		result, err := api.GetSecretValue(ctx, input)

		if err != nil {
			logs.Error().Err(err).Send()
			return AWSSecretsManagerResponse{}, err
		}

		// Parse the request body into the `OAuthAccessResponse` struct
		if err := json.NewDecoder(strings.NewReader(*result.SecretString)).Decode(&asmResp); err != nil {
			logs.Error().Err(err).Send()
			return AWSSecretsManagerResponse{}, err
		}

		// create the go-cache for USC client_id and client_secret
		localCache.Set("client_id", asmResp.ClientID, cache.NoExpiration)
		localCache.Set("client_secret", asmResp.ClientSecret, cache.NoExpiration)

		return asmResp, nil

	} else {
		// take from localCache
		asmResp.ClientID = client_id.(string)
		asmResp.ClientSecret = client_secret.(string)

		return asmResp, nil
	}

}

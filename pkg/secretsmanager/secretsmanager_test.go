package secretsmanager

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/patrickmn/go-cache"
)

type mockSecretsManagerAPI func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)

func (m mockSecretsManagerAPI) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetFromSecretManager(t *testing.T) {

	OKcases := []struct {
		testName       string
		client         func(t *testing.T) SecretsManagerAPI
		expectClientID string
		expectSecret   string
	}{
		{
			testName: "GetClientIDAndSecret without local cache OK",
			client: func(t *testing.T) SecretsManagerAPI {
				return mockSecretsManagerAPI(func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
					t.Helper()

					arnString := "arn:aws:secretsmanager:ap-southeast-1:223206298738:secret:prod/SDS/test-pHTIyd"
					nameString := "prod/SDS/test"
					secretStr := "{\"client_id\":\"628efb01d9b08800b9b6fa9c\",\"client_secret\":\"ESusZN2EGwhcdn+RbSJiO3QBCIvshk+KUYbPxH5yh1A=\"}"
					versionIDString := "0ee1fcb2-534a-4b3e-8242-db9cb49f9edf"

					temp := secretsmanager.GetSecretValueOutput{
						ARN:          &arnString,
						Name:         &nameString,
						SecretString: &secretStr,
						VersionId:    &versionIDString,
					}

					return &temp, nil

				})
			},
			expectClientID: "628efb01d9b08800b9b6fa9c",
			expectSecret:   "ESusZN2EGwhcdn+RbSJiO3QBCIvshk+KUYbPxH5yh1A=",
		},

		{
			testName: "GetClientIDAndSecret with local cache OK",
			client: func(t *testing.T) SecretsManagerAPI {
				return mockSecretsManagerAPI(func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
					t.Helper()

					localCache.Set("client_id", "628efb01d9b08800b9b6fa9c", cache.NoExpiration)
					localCache.Set("client_secret", "ESusZN2EGwhcdn+RbSJiO3QBCIvshk+KUYbPxH5yh1A=", cache.NoExpiration)

					arnString := "arn:aws:secretsmanager:ap-southeast-1:223206298738:secret:prod/SDS/test-pHTIyd"
					nameString := "prod/SDS/test"
					secretStr := "{\"client_id\":\"628efb01d9b08800b9b6fa9c\",\"client_secret\":\"ESusZN2EGwhcdn+RbSJiO3QBCIvshk+KUYbPxH5yh1A=\"}"
					versionIDString := "0ee1fcb2-534a-4b3e-8242-db9cb49f9edf"

					temp := secretsmanager.GetSecretValueOutput{
						ARN:          &arnString,
						Name:         &nameString,
						SecretString: &secretStr,
						VersionId:    &versionIDString,
					}

					return &temp, nil

				})
			},
			expectClientID: "628efb01d9b08800b9b6fa9c",
			expectSecret:   "ESusZN2EGwhcdn+RbSJiO3QBCIvshk+KUYbPxH5yh1A=",
		},

		{
			testName: "GetClientIDAndSecret with Error",
			client: func(t *testing.T) SecretsManagerAPI {
				return mockSecretsManagerAPI(func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
					t.Helper()

					arnString := "arn:aws:secretsmanager:ap-southeast-1:223206298738:secret:prod/SDS/test-pHTIyd"
					nameString := "prod/SDS/test"
					secretStr := ""
					versionIDString := "0ee1fcb2-534a-4b3e-8242-db9cb49f9edf"

					temp := secretsmanager.GetSecretValueOutput{
						ARN:          &arnString,
						Name:         &nameString,
						SecretString: &secretStr,
						VersionId:    &versionIDString,
					}

					return &temp, nil

				})
			},
			expectClientID: "628efb01d9b08800b9b6fa9c",
			expectSecret:   "ESusZN2EGwhcdn+RbSJiO3QBCIvshk+KUYbPxH5yh1A=",
		},
	}

	for _, tt := range OKcases {
		t.Run(tt.testName, func(t *testing.T) {
			ctx := context.TODO()
			content, err := GetClientIDAndSecret(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if content.ClientID != tt.expectClientID {
				t.Errorf("got client ID %s, but expected %s", content.ClientID, tt.expectClientID)
			}

			if content.ClientSecret != tt.expectSecret {
				t.Errorf("got client secret %s, but expected %s", content.ClientSecret, tt.expectSecret)
			}
		})
	}

}

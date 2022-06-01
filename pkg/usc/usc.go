package usc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"gitlab.com/pos_malaysia/golib/aws"
	"gitlab.com/pos_malaysia/golib/env"

	sm "track-and-trace-api-server/pkg/secretsmanager"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"gitlab.com/pos_malaysia/golib/logs"
)

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

var getClientIDAndSecret = func(ctx context.Context) (sm.AWSSecretsManagerResponse, error) {
	// get client_id and client_secret from AWS Secrets Manager
	secretsManagerAPI := secretsmanager.NewFromConfig(aws.GetAWSConfig())
	return sm.GetClientIDAndSecret(ctx, secretsManagerAPI)
}

var decodeResponse = func(resp *http.Response) (OAuthAccessResponse, error) {
	var oaResp OAuthAccessResponse

	// Parse the request body into the `OAuthAccessResponse` struct
	err := json.NewDecoder(resp.Body).Decode(&oaResp)
	return oaResp, err
}

// GetAccessToken returns the API token, token's expiration time and error(if any) from USC(Universal Service Connector)
func GetAccessToken(ctx context.Context) (string, int, error) {
	GetUSCAccessTokenURL := env.Get("USC_GET_ACCESS_TOKEN_URL")
	method := "POST"

	data := url.Values{}

	uscSecret, err := getClientIDAndSecret(ctx)

	if err != nil {
		logs.Error().Err(err).Send()
		return "", 0, err
	}

	data.Set("client_id", uscSecret.ClientID)
	data.Add("client_secret", uscSecret.ClientSecret)
	data.Add("grant_type", "client_credentials")
	data.Add("scope", "as2corporate.v2trackntracewebapijson.all")

	payload := strings.NewReader(data.Encode())

	client := &http.Client{}

	req, err := http.NewRequest(method, GetUSCAccessTokenURL, payload)

	if err != nil {
		logs.Error().Err(err).Send()
		return "", 0, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		logs.Error().Err(err).Send()
		return "", 0, err
	}

	defer resp.Body.Close()

	oaResp, err := decodeResponse(resp)

	if err != nil {
		logs.Error().Err(err).Send()
		return "", 0, err
	}

	return oaResp.AccessToken, oaResp.ExpiresIn, nil

}

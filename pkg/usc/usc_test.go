package usc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	sm "track-and-trace-api-server/pkg/secretsmanager"

	"gitlab.com/pos_malaysia/golib/aws"
)

func TestGetAccessToken(t *testing.T) {

	os.Setenv("USC_GET_ACCESS_TOKEN_URL", "https://gateway-usc.pos.com.my/security/connect/token")

	// test GetAccessToken OK
	t.Run("GetAccessTokenOK", func(t *testing.T) {

		aws.InitAWS("ap-southeast-1")

		ctx := context.Background()

		// mock getClientIDAndSecret function
		getClientIDAndSecret = func(ctx context.Context) (sm.AWSSecretsManagerResponse, error) {
			return sm.AWSSecretsManagerResponse{
				ClientID:     "babdsbasbdasd",
				ClientSecret: "alksdladkaldsldllsdkslk",
			}, nil
		}

		// mock decodeResponse function
		decodeResponse = func(resp *http.Response) (OAuthAccessResponse, error) {

			jsonData := struct {
				AccessToken string `json:"access_token"`
				ExpiresIn   int    `json:"expires_in"`
			}{
				AccessToken: "abc12311231abc",
				ExpiresIn:   300,
			}

			body, _ := json.Marshal(jsonData)

			resp = &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
			}

			// Parse the request body into the `OAuthAccessResponse` struct
			var oaResp OAuthAccessResponse
			err := json.NewDecoder(resp.Body).Decode(&oaResp)
			return oaResp, err

		}

		testAccessToken, testExpiresIn, err := GetAccessToken(ctx)

		if err != nil {
			t.Error("error : ", err.Error())
		}

		if testAccessToken != "abc12311231abc" {
			t.Errorf("expected abc12311231abc, but got %s", testAccessToken)
		}

		if testExpiresIn != 300 {
			t.Errorf("expected 300, but got %d", testExpiresIn)
		}

	})

	// test getClientIDAndSecret Error
	t.Run("decodeResponseError", func(t *testing.T) {

		aws.InitAWS("ap-southeast-1")

		ctx := context.Background()

		// mock getClientIDAndSecret function and return error on purpose
		getClientIDAndSecret = func(ctx context.Context) (sm.AWSSecretsManagerResponse, error) {
			return sm.AWSSecretsManagerResponse{
				ClientID:     "babdsbasbdasd",
				ClientSecret: "alksdladkaldsldllsdkslk",
			}, errors.New("some strange error")
		}

		// mock decodeResponse function
		decodeResponse = func(resp *http.Response) (OAuthAccessResponse, error) {

			jsonData := struct {
				AccessToken string `json:"access_token"`
				ExpiresIn   int    `json:"expires_in"`
			}{
				AccessToken: "abc12311231abc",
				ExpiresIn:   300,
			}

			body, _ := json.Marshal(jsonData)

			resp = &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
			}

			// Parse the request body into the `OAuthAccessResponse` struct
			var oaResp OAuthAccessResponse
			err := json.NewDecoder(resp.Body).Decode(&oaResp)
			return oaResp, err

		}

		_, _, err := GetAccessToken(ctx)

		if err == nil {
			t.Error("expected error but got err == nil")
		}

	})

	// test getClientIDAndSecret Error
	t.Run("getClientIDAndSecretError", func(t *testing.T) {

		aws.InitAWS("ap-southeast-1")

		ctx := context.Background()

		// mock getClientIDAndSecret function
		getClientIDAndSecret = func(ctx context.Context) (sm.AWSSecretsManagerResponse, error) {
			return sm.AWSSecretsManagerResponse{
				ClientID:     "babdsbasbdasd",
				ClientSecret: "alksdladkaldsldllsdkslk",
			}, nil
		}

		// mock decodeResponse function and return error on purpose
		decodeResponse = func(resp *http.Response) (OAuthAccessResponse, error) {

			jsonData := struct {
				AccessToken string `json:"access_token"`
				ExpiresIn   int    `json:"expires_in"`
			}{
				AccessToken: "abc12311231abc",
				ExpiresIn:   300,
			}

			body, _ := json.Marshal(jsonData)

			resp = &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
			}

			// Parse the request body into the `OAuthAccessResponse` struct
			var oaResp OAuthAccessResponse
			json.NewDecoder(resp.Body).Decode(&oaResp)

			return oaResp, errors.New("some strange error")

		}

		_, _, err := GetAccessToken(ctx)

		if err == nil {
			t.Error("expected error but got err == nil")
		}

	})

}

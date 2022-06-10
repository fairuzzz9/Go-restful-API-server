package responses

import (
	"encoding/json"
	"go-skeleton-rest-app/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type getReponseMessageByCodeTestCase struct {
	name          string
	input         string
	expectedJSON  string
	expectedError string
}

// NOTE : test coverage is at 80% because unable to simulate json.Marshal and Unmarshal error.

// to test json.Unmarshal error
// see stackoverflow.com/questions/32708717
//{name: "test json unmarshal error", input: 666666, expectedJSON: `null`, expectedError: "invalid character 'w' looking for beginning of value"},

var getReponseMessageByCodeTestCases = []getReponseMessageByCodeTestCase{
	{name: `test negative one`, input: "-1", expectedJSON: `null`, expectedError: `invalid code`},
	{name: `test negative ten`, input: "-10", expectedJSON: `null`, expectedError: `invalid code`},
	{name: `test success code`, input: "S0000", expectedJSON: `{"code":"S0000","message":"success","client_request_id":"","server_trace_id":"","data":null}`, expectedError: ``},
	{name: `test database error code`, input: "E1000", expectedJSON: `{"code":"E1000","message":"db error","client_request_id":"","server_trace_id":"","data":null}`, expectedError: ``},
	{name: `test no such code`, input: "9999", expectedJSON: `null`, expectedError: `invalid code`},
}

const (
	TestSuccessCode            = "S0000"
	TestDBErrorCode            = "E1000"
	TestJSONUnMarshalErrorCode = "E2000"
)

var testResponseMessages = map[string]string{
	TestSuccessCode:            "success",
	TestDBErrorCode:            "db error",
	TestJSONUnMarshalErrorCode: "json marshal error",
}

func assertGetReponseMessageByCode(t *testing.T, input string, expectedJSON, expectedError string) {
	actual, actualError := GetReponseMessageByCode(input)

	actualTemp, _ := json.Marshal(actual)
	actualString := string(actualTemp)

	var actualErrorString string
	if actualError != nil {
		actualErrorString = actualError.Error()
	}

	if actualString != expectedJSON {
		t.Fatalf("error! expected JSON string is \"%s\", but got \"%s\" \n", expectedJSON, actualString)
	}

	if actualErrorString != expectedError {
		t.Fatalf("error! expected error string is \"%s\", but got \"%s\" \n", expectedError, actualErrorString)
	}
}

func setupTestResponseMessages() func(tb testing.TB) {

	for code, message := range testResponseMessages {
		if responseData, err := json.Marshal(&models.StandardJSONResponse{Code: code, Message: message}); err != nil {
			panic(err)
		} else {
			// load into map
			responseMap.Store(code, responseData)
		}
	}

	return func(tb testing.TB) {

		// don't use these two methods to clear sync.Map as it can be dangerous
		// and introduce race condition
		//responseMap = nil
		//responseMap = sync.Map{}

		// safe way to clear sync.Map
		responseMap.Range(func(key, value interface{}) bool {
			responseMap.Delete(key)
			return true
		})
	}

}

func TestGetReponseMessageByCode(t *testing.T) {

	teardownTestResponseMessages := setupTestResponseMessages()
	defer teardownTestResponseMessages(t)

	for _, c := range getReponseMessageByCodeTestCases {
		t.Run(c.name, func(t *testing.T) {
			assertGetReponseMessageByCode(t, c.input, c.expectedJSON, c.expectedError)
		})
	}
}

//------------------------------------------------------------------------

type responseWithErrorTestCase struct {
	name              string
	clientRequestID   string
	serverTraceID     string
	responseErrorCode string
	expectedError     error
}

var responseWithErrorTestCases = []responseWithErrorTestCase{
	{name: `test ok`, clientRequestID: "abc123", serverTraceID: "adij234jdhbguuiw434sd", responseErrorCode: TestSuccessCode, expectedError: nil},
	{name: `test expect error`, clientRequestID: "abc123", serverTraceID: "adij234jdhbguuiw434sd", responseErrorCode: "abc", expectedError: echo.NewHTTPError(http.StatusInternalServerError, "invalid code")},
}

func assertResponseWithError(t *testing.T, clientRequestID, serverTraceID, responseErrorCode string, expectedError error) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	echoContext := e.NewContext(req, rec)

	actualError := ResponseWithError(clientRequestID, serverTraceID, echoContext, responseErrorCode)

	t.Log(actualError)

	//if actualError.Error() != expectedError.Error() {
	//	t.Log(expectedError)
	//t.Errorf("expected error %v, but got %v instead.", expectedError, actualError)
	//}
}

func TestResponseWithError(t *testing.T) {

	teardownTestResponseMessages := setupTestResponseMessages()
	defer teardownTestResponseMessages(t)

	for _, c := range responseWithErrorTestCases {
		t.Run(c.name, func(t *testing.T) {
			assertResponseWithError(t, c.clientRequestID, c.serverTraceID, c.responseErrorCode, c.expectedError)
		})
	}
}

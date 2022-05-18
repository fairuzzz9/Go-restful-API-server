package responses

import (
	"encoding/json"
	"testing"
)

type getReponseMessageByCodeTestCase struct {
	name          string
	input         int
	expectedJSON  string
	expectedError string
}

// NOTE : test coverage is at 80% because unable to simulate json.Marshal and Unmarshal error.

// to test json.Unmarshal error
// see stackoverflow.com/questions/32708717
//{name: "test json unmarshal error", input: 666666, expectedJSON: `null`, expectedError: "invalid character 'w' looking for beginning of value"},

var getReponseMessageByCodeTestCases = []getReponseMessageByCodeTestCase{
	{name: `test negative one`, input: -1, expectedJSON: `null`, expectedError: `invalid code`},
	{name: `test negative ten`, input: -10, expectedJSON: `null`, expectedError: `invalid code`},
	{name: `test success code`, input: 0, expectedJSON: `{"code":0,"message":"success","data":null}`, expectedError: ``},
	{name: `test database error code`, input: 1000, expectedJSON: `{"code":1000,"message":"db error","data":null}`, expectedError: ``},
	{name: `test no such code`, input: 9999, expectedJSON: `null`, expectedError: `invalid code`},
}

const (
	TestSuccessCode            = 0
	TestDBErrorCode            = 1000
	TestJSONUnMarshalErrorCode = 2000
)

var testResponseMessages = map[int]string{
	TestSuccessCode:            "success",
	TestDBErrorCode:            "db error",
	TestJSONUnMarshalErrorCode: "json marshal error",
}

func assertGetReponseMessageByCode(t *testing.T, input int, expectedJSON, expectedError string) {
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
		if responseData, err := json.Marshal(&JSONResponse{Code: code, Message: message}); err != nil {
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

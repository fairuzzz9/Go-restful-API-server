package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TestServerRequestID(t *testing.T) {

	testServerRequestID := "<test-home-server-request-id>"

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Add(echo.HeaderXRequestID, testServerRequestID)

	recorder := httptest.NewRecorder()

	e := echo.New()

	echoContext := e.NewContext(request, recorder)

	if recorder.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected %d", recorder.Code, http.StatusOK)
	}

	// see https://github.com/labstack/echo/blob/master/middleware/request_id_test.go
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	rid := middleware.RequestIDWithConfig(middleware.RequestIDConfig{})
	h := rid(handler)
	h(echoContext)

	serverRequestID := recorder.Header().Get(echo.HeaderXRequestID)
	t.Log(serverRequestID)

	if serverRequestID == "" {
		t.Errorf("expected %s, but got %s", testServerRequestID, serverRequestID)
	}

}

func TestHome(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	e := echo.New()

	echoContext := e.NewContext(request, recorder)

	Home(echoContext)

	if recorder.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected %d", recorder.Code, http.StatusOK)
	}

	requestBody := recorder.Body.String()

	if !strings.Contains(requestBody, "home") {
		t.Errorf("response body \"%s\" does not contain \"I'm ok\" ", requestBody)
	}

	t.Log(requestBody)
}

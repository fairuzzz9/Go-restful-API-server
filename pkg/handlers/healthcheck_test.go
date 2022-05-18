package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHealthCheck(t *testing.T) {
	e := echo.New()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)

	echoContext := e.NewContext(request, recorder)

	HealthCheck(echoContext)

	if recorder.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected %d", recorder.Code, http.StatusOK)
	}

	if !strings.Contains(recorder.Body.String(), "I'm ok") {
		t.Errorf("response body \"%s\" does not contain \"I'm ok\" ", recorder.Body.String())
	}
}

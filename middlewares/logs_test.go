package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kokhno-nikolay/letsgochat/middlewares"
)

func TestLogging(t *testing.T) {
	mw := middlewares.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

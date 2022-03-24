package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kokhno-nikolay/letsgochat/middlewares"
)

func TestRecovery(t *testing.T) {
	mw := middlewares.Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

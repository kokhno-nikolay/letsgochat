package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func TestLogger(t *testing.T) {
	router := gin.Default()
	log := logrus.New()
	router.Use(middlewares.Logger(log), gin.Recovery())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, resp.Code, http.StatusNotFound)
}

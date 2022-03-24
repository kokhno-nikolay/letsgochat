package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kokhno-nikolay/letsgochat/api"
)

func TestSignUp(t *testing.T) {
	router := api.NewHandler(api.Deps{}).Init()
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	jsonParam := `{"username":"test", "password":"test"}`
	req, err := http.NewRequest("POST", "/user", strings.NewReader(string(jsonParam)))
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

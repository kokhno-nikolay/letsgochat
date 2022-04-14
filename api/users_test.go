package api_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

func TestSignUp(t *testing.T) {
	repos := repository.NewRepositories(&sql.DB{})
	router := api.NewHandler(api.Deps{repos}).Init()
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	jsonParam := `{"username":"test123123", "password":"test123123123"}`
	req, err := http.NewRequest("POST", "/user", strings.NewReader(string(jsonParam)))
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

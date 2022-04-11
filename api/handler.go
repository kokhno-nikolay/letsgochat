package api

import (
	"github.com/kokhno-nikolay/letsgochat/middlewares"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/kokhno-nikolay/letsgochat/repository"
)

type Handler struct {
	UserRepo  repository.Users
	TokenRepo repository.Token
	Host      string
}

type Deps struct {
	Repos *repository.Repositories
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		UserRepo:  deps.Repos.Users,
		TokenRepo: deps.Repos.Token,
		Host:      os.Getenv("HOST_NAME"),
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.HandleMethodNotAllowed = true
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	log := logrus.New()

	router.Use(middlewares.Logger(log), gin.Recovery())
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	/* User handlers */
	router.POST("/user", h.SignUp)
	router.POST("/user/login", h.SignIn)
	router.GET("/user/active", h.GetActiveUsers)

	router.GET("/chat", func(c *gin.Context) {
		token, ok := c.GetQuery("token")
		if !ok {
			c.String(http.StatusUnauthorized, "missing auth token")
			return
		}

		ok, err := h.TokenRepo.CheckToken(token)
		if err != nil {
			c.String(http.StatusInternalServerError, "token cannot be validated")
			return
		}

		if !ok {
			c.String(http.StatusBadRequest, "token invalid")
			return
		}

		defer func() {
			if err := h.TokenRepo.DeleteToken(token); err != nil {
				c.String(http.StatusInternalServerError, "token cannot be validated")
				return
			}

			log.Println("token deleted successfully")
		}()

		h.Chat(c.Writer, c.Request, token)
	})

	return router
}

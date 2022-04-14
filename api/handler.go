package api

import (
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/kokhno-nikolay/letsgochat/middlewares"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

type Handler struct {
	userRepo repository.Users
	sessions map[string]int
	host     string
	mu       sync.Mutex
}

type Deps struct {
	Repos *repository.Repositories
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		userRepo: deps.Repos.Users,
		sessions: make(map[string]int),
		host:     os.Getenv("HOST_NAME"),
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

		if ok := h.checkToken(token); !ok {
			c.String(http.StatusBadRequest, "token invalid")
			return
		}

		defer func() {
			h.deleteToken(token)
			log.Println("token deleted successfully")
		}()

		h.Chat(c.Writer, c.Request, token)
	})

	return router
}

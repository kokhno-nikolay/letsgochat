package api

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/kokhno-nikolay/letsgochat/middlewares"
	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository"
	"github.com/kokhno-nikolay/letsgochat/services"
)

var (
	workersNum = 100
)

type Job struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Handler struct {
	services    *services.Services
	Sessions    map[string]int
	clients     map[*websocket.Conn]bool
	messageCh   chan message
	broadcaster chan models.ChatMessage
	host        string
	mu          sync.Mutex
}

type Deps struct {
	Repos *repository.Repositories
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services:    services,
		Sessions:    make(map[string]int),
		clients:     make(map[*websocket.Conn]bool),
		messageCh:   make(chan message, workersNum),
		broadcaster: make(chan models.ChatMessage),
		host:        os.Getenv("HOST_NAME"),
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

		if ok := h.CheckUserToken(token); !ok {
			c.String(http.StatusBadRequest, "token invalid")
			return
		}

		defer func() {
			if err := h.services.Users.SwitchToInactive(h.Sessions[token]); err != nil {
				return
			}
			h.DeleteSession(token)
		}()

		go h.handleMessages(token)
		h.handleConnections(c.Writer, c.Request)
	})

	return router
}

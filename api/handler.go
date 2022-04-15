package api

import (
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/kokhno-nikolay/letsgochat/middlewares"
	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

type Handler struct {
	userRepo    repository.Users
	messageRepo repository.Messages
	sessions    map[string]int
	clients     map[*websocket.Conn]bool
	broadcaster chan models.ChatMessage
	host        string
	mu          sync.Mutex
}

type Deps struct {
	Repos *repository.Repositories
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		userRepo:    deps.Repos.Users,
		messageRepo: deps.Repos.Messages,
		sessions:    make(map[string]int),
		clients:     make(map[*websocket.Conn]bool),
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

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
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
			h.DeleteSession(token)
			log.Println("token deleted successfully")
		}()

		go h.handleMessages(token)
		h.handleConnections(c.Writer, c.Request)
	})

	return router
}

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kokhno-nikolay/letsgochat/models"
)

type userInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Sign up
// @Tags Users
// @Description Registration a new user in the system
// @Accept  json
// @Produce  json
// @Param input body userInput true "Please enter your username and password to register"
// @Success 200 {object} models.JSONResult{data=string} "Successful server response"
// @Failure 400 {object} models.JSONResult{data=string} "Invalid input request"
// @Failure 500 {object} models.JSONResult{data=string} "Internal server error"
// @Router /user [post]
func (h *Handler) SignUp(c *gin.Context) {
	var inp userInput

	if err := c.BindJSON(&inp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "invalid input request",
		})

		return
	}

	if len(inp.Username) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "username is too short",
		})

		return
	}

	if len(inp.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "password is too short",
		})

		return
	}

	userExists, err := h.services.Users.UserExists(inp.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "something went wrong, please try again. Error: " + err.Error(),
		})

		return
	}

	if userExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("user with username %s already exists", inp.Username),
		})

		return
	}

	user := models.User{Username: inp.Username, Password: inp.Password, Active: false}
	if err := h.services.Users.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "something went wrong, please try again. Error: " + err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": fmt.Sprintf("user with username %s successfully created", inp.Username),
	})
}

// @Summary Sign in
// @Tags Users
// @Description User account login
// @Accept  json
// @Produce  json
// @Param input body userInput true "Please enter your username and password to login"
// @Success 200 {object} models.JSONResult{data=string} "Successful server response"
// @Failure 400 {object} models.JSONResult{data=string} "Invalid input request"
// @Failure 500 {object} models.JSONResult{data=string} "Internal server error"
// @Router /user/login [post]
func (h *Handler) SignIn(c *gin.Context) {
	var inp userInput

	if err := c.BindJSON(&inp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "input decoding error",
		})

		return
	}

	if len(inp.Username) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "username is too short",
		})

		return
	}

	if len(inp.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "password is too short",
		})

		return
	}

	userExists, err := h.services.Users.UserExists(inp.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "something went wrong, please try again. Error: " + err.Error(),
		})

		return
	}

	if !userExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "user does not exist",
		})

		return
	}

	user, err := h.services.Users.FindByUsername(inp.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "something went wrong, please try again. Error: " + err.Error(),
		})

		return
	}

	if user.Password != inp.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "password is invalid",
		})

		return
	}

	var token string
	tokenExists, t := h.CheckUserSession(user.ID)
	if !tokenExists {
		token = uuid.New().String()
		h.mu.Lock()
		h.Sessions[token] = user.ID
		h.mu.Unlock()
	} else {
		token = t
	}

	if err := h.services.Users.SwitchToActive(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "something went wrong, please try again. Error: " + err.Error(),
		})

		return
	}

	url := fmt.Sprintf("wss://%s/chat?token=%s", h.host, token)
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// @Summary Active users
// @Tags Users
// @Description Number of active users in a chat
// @Accept  json
// @Produce  json
// @Success 200 {string} []string "Returns all active users in the chat"
// @Failure 500 {string} string "Internal server error"
// @Router /user/active [get]
func (h *Handler) GetActiveUsers(c *gin.Context) {
	var res []string

	users, err := h.services.Users.GetActiveUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "something went wrong, please try again. Error: " + err.Error(),
		})
	}

	for _, user := range users {
		res = append(res, user.Username)
	}

	c.JSON(http.StatusOK, gin.H{
		"active_users": res,
	})
}

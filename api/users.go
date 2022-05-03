package api

import (
	"fmt"
	"github.com/kokhno-nikolay/letsgochat/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Sign up
// @Tags users
// @Description New user registration
// @Accept  json
// @Produce  json
// @Param input body userInput true "user data"
// @Success 200 {integer} integer 1
// @Failure 400 {string} string "invalid input request"
// @Failure 500 {string} string "something went wrong"
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

	userExists, err := h.userRepo.UserExists(inp.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})

		return
	}

	if userExists >= 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("user with username %s already exists", inp.Username),
		})

		return
	}

	user := models.User{Username: inp.Username, Password: inp.Password, Active: false}
	if err := h.userRepo.Create(user); err != nil {
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
// @Tags User handlers
// @Description User account login
// @Accept  json
// @Produce  json
// @Param input body userInput true "user data"
// @Success 200 {integer} integer 1
// @Failure 400 {string} string "invalid input request"
// @Failure 500 {string} string "something went wrong"
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

	userExists, err := h.userRepo.UserExists(inp.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})

		return
	}

	if userExists <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "user does not exist",
		})

		return
	}

	user, err := h.userRepo.FindByUsername(inp.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
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

	if err := h.userRepo.SwitchToActive(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})

		return
	}

	url := fmt.Sprintf("wss://%s/chat?token=%s", h.host, token)
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func (h *Handler) GetActiveUsers(c *gin.Context) {
	var res []string

	users, err := h.userRepo.GetAllActive()
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

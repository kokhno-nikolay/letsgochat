package api

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kokhno-nikolay/letsgochat/models"
)

func (h *Handler) SignUp(c *gin.Context) {
	var inp models.UserInput

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

	userExists, err := h.UserRepo.UserExists(inp.Username)
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
	if err := h.UserRepo.Create(user); err != nil {
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

func (h *Handler) SignIn(c *gin.Context) {
	var inp models.UserInput

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

	userExists, err := h.UserRepo.UserExists(inp.Username)
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

	user, err := h.UserRepo.FindByUsername(inp.Username)
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
	tokenExists, t := h.checkUserSession(user.ID)
	if !tokenExists {
		token = uuid.New().String()
		h.Sessions[token] = user.ID
	} else {
		token = t
	}

	url := fmt.Sprintf("ws://%s/chat?token=%s", h.Host, token)
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func (h *Handler) GetActiveUsers(c *gin.Context) {
	var res []string

	users, err := h.UserRepo.GetAllActiveUsers()
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

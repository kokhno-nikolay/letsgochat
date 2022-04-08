package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"

	"github.com/kokhno-nikolay/letsgochat/models"
)

func (h *Handler) SignUp(c *gin.Context) {
	var inp models.UserInput

	if err := c.BindJSON(&inp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "invalid input body",
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
		"message": "user created successfully",
	})
}

func (h *Handler) SignIn(c *gin.Context) {
	var inp models.UserInput

	if err := c.BindJSON(&inp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "invalid input body",
		})

		return
	}

	userExists, err := h.UserRepo.UserExists(inp.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": err.Error(),
		})

		return
	}

	if !userExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "user does not exist",
		})

		return
	}

	user, err := h.UserRepo.FindByUsername(inp.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": err.Error(),
		})

		return
	}

	t := models.Token{UUID: uuid.New(), UserId: user.ID}
	token, err := h.TokenRepo.CreateToken(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "token creation error: " + err.Error(),
		})

		return
	}

	url := fmt.Sprintf("ws://%s/chat?token=%s", h.Host, token.UUID.String())
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

package http

import (
	"core-services/domains/users"
	"core-services/domains/users/models/requests"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHttp struct {
	uc users.UserUseCase
}

func NewUserHttp(uc users.UserUseCase) *UserHttp {
	return &UserHttp{uc}
}

func (h *UserHttp) Register(c *gin.Context) {
	var req requests.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.uc.Register(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *UserHttp) Login(c *gin.Context) {
	var req requests.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.uc.Login(&req)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHttp) GetProfile(c *gin.Context) {
	username := c.Param("username")
	user, err := h.uc.GetProfile(username)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

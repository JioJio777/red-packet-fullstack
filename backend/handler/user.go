package handler

import (
	"net/http"

	"red-packet/pkg/response"
	"red-packet/service"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	user, err := service.Register(req.Username, req.Password)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	token, err := service.Login(req.Username, req.Password)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token":      token,
		"expires_in": 86400,
	})
}

func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := service.GetProfile(userID.(uint64))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, "internal error")
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"balance":  user.Balance,
	})
}

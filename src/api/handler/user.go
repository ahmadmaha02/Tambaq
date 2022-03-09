package handler

import (
	"net/http"
	"project/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Register(c *gin.Context) {

	var userRequest user.PostRegisterBody
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Body is invalid.",
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	newUser, err := h.userService.Register(userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error when inserting into the database.",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil Membuat Akun",
		"status":  "Sukses",
		"data": gin.H{
			"id": newUser.ID,
		},
	})
}

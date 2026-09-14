package controllers

import (
	"net/http"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthInputRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthInputLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(context *gin.Context) {
	var input AuthInputRegister

	// validation
	err := context.ShouldBindBodyWithJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal enkripsi password",
		})
		return
	}

	// simpan ke db
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	userCreated := config.DB.Create(&user).Error

	if userCreated != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Email Sudah terdaftar",
		})
		return
	}

	context.JSON(http.StatusCreated,
		gin.H{
			"message": "Berhasil register",
			"user": gin.H{
				"id":     user.ID,
				"name":   user.Name,
				"email":  user.Email,
				"events": user.Events,
			},
		})
}

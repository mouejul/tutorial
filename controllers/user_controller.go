package controllers

import (
	"net/http"
	"os"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func LoginUser(context *gin.Context) {
	var input AuthInputLogin
	// validation
	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			})
		return
	}

	var user models.User
	userData := config.DB.Where("email = ?", input.Email).First(&user).Error
	if userData != nil {
		context.JSON(http.StatusUnauthorized,
			gin.H{
				"error": "email belum terdaftar",
			})
		return
	}

	errMatchPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if errMatchPassword != nil {
		context.JSON(http.StatusUnauthorized,
			gin.H{
				"error": "Password salahh",
			})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{
				"Error": "gagal membuat token",
			})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"events": user.Events,
		},
	})

}

func GetCurrentUser(context *gin.Context) {
	// ambil userId
	userId, exists := context.Get("userID")
	if !exists {
		context.JSON(http.StatusInternalServerError,
			gin.H{
				"error": "tidak authentication",
			})
		return
	}

	var user models.User

	userData := config.DB.Select("id", "name", "email").First(&user, userId).Error
	if userData != nil {
		context.JSON(http.StatusNotFound,
			gin.H{
				"error": "User tidak ditemukan",
			})
		return
	}

	context.JSON(http.StatusOK,
		gin.H{
			"user": gin.H{
				"name":   user.Name,
				"email":  user.Email,
				"id":     user.ID,
				"events": user.Events,
			},
		})

}

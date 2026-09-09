package controllers

	import "github.com/gin-gonic/gin"


type AuthInputRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthInputLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(context *gin.Context){
	var input AuthInputRegister

	// validation
	err := context.ShouldBindBodyWithJSON()
}

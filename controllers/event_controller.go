package controllers

import (
	"net/http"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
)

func CreateEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	event.UserID = 1

	config.DB.Create(&event)
	context.JSON(http.StatusCreated, gin.H{
		"message": "Data berhasil dibuat",
		"event":   event,
	})
}

func GetEvents(context *gin.Context) {
	var events []models.Event

	config.DB.Find(&events)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil tampil",
		"event":   events,
	})
}

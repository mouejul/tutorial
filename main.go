package main

import (
	"net/http"
	"time"

	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	// Route
	api := server.Group("/api")
	{
		api.POST("/events", createEvent)
		api.GET("/events", getEvents)
	}

	server.Run(":8080")

}

// func handler
func getEvents(context *gin.Context) {
	events := models.GetAllEvents()

	context.JSON(http.StatusOK, events)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data",
			"error":   err.Error(),
		})
		return
	}

	tm := time.Now()
	// dummy
	event.Id = 1
	event.UserId = 1
	event.DatTime = tm

	// save inputan
	event.Save()

	context.JSON(http.StatusCreated, gin.H{
		"message": "create event",
		"event":   event,
	})
}

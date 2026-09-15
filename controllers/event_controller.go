package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

func initImageKit() *imagekit.Client {
	client := imagekit.NewClient(
		option.WithPrivateKey(os.Getenv("IAMGEKIT_PRIVATE_KEY")),
	)
	return &client
}

func CreateEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	// menerima file form data
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "gambar wajib di upload",
		})
		return
	}
	defer file.Close()

	// 1. ipload file ke imageKit
	fileName := header.Filename
	ik := initImageKit()
	uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: fileName,
	})

	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal upload gambar imageKit",
		})
	}

	parseTime, _ := time.Parse(time.RFC3339, c.PostForm("datetime"))

	// simpan ke database
	event := models.Event{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Location:    c.PostForm("location"),
		Datetime:    parseTime,
		Image:       uploadRes.URL,
		ImageID:     uploadRes.FileID,
		UserID:      userID.(int),
	}
	config.DB.Create(&event)
	c.JSON(http.StatusCreated, gin.H{
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

func GetEventbyId(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error

	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "event tidak ditemukan",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Data detail event berhasil tampil",
		"event":   event,
	})
}

func UpdateEvent(context *gin.Context) {
	userID, _ := context.Get("userID")

	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error

	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Event tidak ditemukan",
		})
		return
	}
	if event.UserID != userID.(int) {
		context.JSON(http.StatusForbidden,
			gin.H{
				"error": "Kamu tidak update event user lain",
			})
		return
	}
	var input models.Event
	err := context.ShouldBindBodyWithJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	config.DB.Model(&event).Updates(input)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data Berhasil diupdate",
		"event":   event,
	})
}

func DeleteEvent(context *gin.Context) {
	userID, _ := context.Get("userID")

	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Event tidak ditemukan",
		})
		return
	}

	if event.UserID != userID.(int) {
		context.JSON(http.StatusForbidden,
			gin.H{
				"error": "Kamu tidak update event user lain",
			})
		return
	}

	config.DB.Unscoped().Delete(&event)

	context.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil di delete",
	})

}

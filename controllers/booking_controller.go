package controllers

import (
	"fmt"
	"net/http"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingInput struct {
	Phone   string `json:"phone" binding:"required"`
	EventId int    `json:"eventId" binding:"required"`
}

func CreateBookingEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input BookingInput
	var booking models.Booking

	errValidation := c.ShouldBindJSON(&input)
	if errValidation != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Error(),
		})
		return
	}

	//  kondisi jika user sudah pernah booking event id
	Booking := config.DB.Where("user_id = ? AND event_id = ?",
		userID.(int), input.EventId).First(&booking).Error

	if Booking == nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Anda sudah booking event ini",
			})
		return
	}

	// generate Booking code
	CodeBooking := fmt.Sprintf("BK-%sE%dU%d",
		time.Now().Format("20060102"), input.EventId, userID.(int))

	// simpan ke database

	bookingData := models.Booking{
		Phone:       input.Phone,
		EventId:     input.EventId,
		BookingCode: CodeBooking,
		UserID:      userID.(int),
	}

	errCreateBooking := config.DB.Create(&bookingData).Error
	if errCreateBooking != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "create Booking gagal",
			})
		return
	}

	c.JSON(http.StatusCreated,
		gin.H{
			"message": "Berhasil Daftar Event",
		})
}

func GetBookingbyUser(c *gin.Context) {
	var booking []models.Booking
	userID, _ := c.Get("userID")

	errBookingData := config.DB.Preload("Event").Preload("Event.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Where("user_id = ?", userID).Find(&booking).Error

	if errBookingData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "event tidak ditemukan",
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"booking": booking,
	})
}

func DeleteBooking(c *gin.Context) {
	userID, _ := c.Get("userID")

	var booking models.Booking

	paramsId := c.Param("id")

	bookingData := config.DB.First(&booking, paramsId).Error

	if bookingData != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "event tidak ditemukan",
		})
		return
	}

	if booking.UserID != userID.(int) {

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Kamu tidak bisa menghapus booking user lain",
		})
		return
	}

	config.DB.Unscoped().Delete(&booking)

	c.JSON(http.StatusOK, gin.H{
		"message": "event booking event berhasil dihapus",
	})
}

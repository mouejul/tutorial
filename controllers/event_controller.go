package controllers

import (
	"context"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
	"gorm.io/gorm"
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

func GetEvents(c *gin.Context) {
	var events []models.Event

	// 1 inisiasi dasar query di gorm
	query := config.DB.Model(&models.Event{})

	// 2 tangkap fungsi filter by query
	search := c.Query("search")

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 3. hitung jumlah data sebelum dilimit
	var totalRows int64
	query.Count(&totalRows)

	// 4 tangkap parameter query dan masukan nilai default
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("page", "6")

	page, errPage := strconv.Atoi(pageStr)
	if errPage != nil || page < 1 {
		page = 1
	}

	limit, errLimit := strconv.Atoi(limitStr)
	if errLimit != nil || limit < 1 {
		limit = 6
	}

	// 5 hitung offset
	offset := (page - 1) * limit

	// 6 hitung data perhalaman/page
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	// 7 eksekusi semua fitur yang dibuat diatas
	if err := query.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Limit(limit).Offset(offset).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal mengambil data event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil tampil",
		"event":   events,
		"meta": gin.H{
			"page":      page,
			"limit":     limit,
			"totalRow":  totalRows,
			"totalPage": totalPages,
		},
	})
}

func GetEventbyId(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).First(&event, paramsId).Error

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

func UpdateEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var event models.Event
	paramsId := c.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error

	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event tidak ditemukan",
		})
		return
	}
	if event.UserID != userID.(int) {
		c.JSON(http.StatusForbidden,
			gin.H{
				"error": "Kamu tidak update event user lain",
			})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		ik := initImageKit()
		// upload file gambar baru
		fileName := header.Filename

		uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
			File:     file,
			FileName: fileName,
		})

		if errUpload == nil {
			// hapus gambar lama
			if event.ImageID != "" {
				ik.Files.Delete(context.Background(), event.ImageID)
			}
			// upload field gambar
			event.Image = uploadRes.URL
			event.ImageID = uploadRes.FileID

		}

	}

	if name := c.PostForm("name"); name != "" {
		event.Name = name
	}

	if description := c.PostForm("description"); description != "" {
		event.Description = description
	}
	if location := c.PostForm("name"); location != "" {
		event.Location = location
	}
	if dateTimeStr := c.PostForm("datetime"); dateTimeStr != "" {

		parseTime, errParse := time.Parse(time.RFC3339, dateTimeStr)
		if errParse == nil {
			event.Datetime = parseTime
		}
	}

	config.DB.Save(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Berhasil diupdate",
		"event":   event,
	})
}

func DeleteEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var event models.Event
	paramsId := c.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event tidak ditemukan",
		})
		return
	}

	if event.UserID != userID.(int) {
		c.JSON(http.StatusForbidden,
			gin.H{
				"error": "Kamu tidak update event user lain",
			})
		return
	}

	if event.ImageID != "" {
		ik := initImageKit()
		ik.Files.Delete(context.Background(), event.ImageID)
	}
	config.DB.Unscoped().Delete(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil di delete",
	})

}

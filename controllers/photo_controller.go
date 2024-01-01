// controllers/photo_controller.go

package controllers

import (
	"finpro-golang2/database"
	"finpro-golang2/helpers"
	"finpro-golang2/models"
	"net/http"
	"strconv"
    "log"
	"github.com/gin-gonic/gin"
)

func CreatePhoto(c *gin.Context) {
    // Baca data yang diterima dari body request
    var photo models.Photo
    if err := c.ShouldBindJSON(&photo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Simpan foto ke database
    if err := database.DB.Create(&photo).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    // Kirim respons
    c.JSON(http.StatusOK, photo)
}

// GetPhotos digunakan untuk mendapatkan daftar semua foto
func GetPhotos(c *gin.Context) {
    // Dapatkan token dari header permintaan
    tokenString := c.Request.Header.Get("Authorization")

    // Verifikasi dan ekstrak ID pengguna dari token
    userID, err := helpers.ExtractUserIDFromToken(tokenString)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Dapatkan daftar foto dari database berdasarkan ID pengguna
    var userPhotos []models.Photo
    if err := database.DB.Where("user_id = ?", userID).Find(&userPhotos).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    log.Printf("UserID from token: %d", userID)

    // Kirim respons
    c.JSON(http.StatusOK, userPhotos)
}


func UpdatePhoto(c *gin.Context) {
	// Mendapatkan ID foto dari parameter URL
	photoID, err := strconv.Atoi(c.Param("photoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	// Mendapatkan foto dari database berdasarkan ID
	var photo models.Photo
	if err := database.DB.First(&photo, photoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	// Baca data yang diterima dari body request
	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan perubahan ke database
	if err := database.DB.Save(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Kirim respons
	c.JSON(http.StatusOK, photo)
}

func DeletePhoto(c *gin.Context) {
	// Mendapatkan ID foto dari parameter URL
	photoID, err := strconv.Atoi(c.Param("photoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	// Hapus foto dari database berdasarkan ID
	if err := database.DB.Delete(&models.Photo{}, photoID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Kirim respons
	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}
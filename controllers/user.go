package controllers

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	jwtV3 "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
)

type StrukturUserTambah struct {
	Nama     string `json:"nama" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type StrukturLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UserTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataUser StrukturUserTambah

	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	// Enkripsi password pakai SHA1
	sha := sha1.New()
	sha.Write([]byte(dataUser.Password))
	encryptedString := fmt.Sprintf("%x", sha.Sum(nil))

	modelUser := models.User{
		Nama:     dataUser.Nama,
		Username: dataUser.Username,
		Password: encryptedString,
	}

	db.Create(&modelUser)
	c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil daftar user", "data": modelUser})
}

func UserLogin(c *gin.Context) (interface{}, error) {
	db := c.MustGet("db").(*gorm.DB)
	var dataUser StrukturLogin

	if err := c.ShouldBindJSON(&dataUser); err != nil {
		return nil, jwtV3.ErrMissingLoginValues
	}

	sha := sha1.New()
	sha.Write([]byte(dataUser.Password))
	encryptedString := fmt.Sprintf("%x", sha.Sum(nil))

	var modelUser models.User
	cekUser := db.Where("username = ?", dataUser.Username).Where("password = ?", encryptedString).First(&modelUser)
	
	if cekUser.Error == nil {
		return modelUser, nil
	}
	return nil, jwtV3.ErrFailedAuthentication
}
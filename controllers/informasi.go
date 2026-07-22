package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
)

// Binding JSON buat Postman
type StrukturInformasi struct {
	ID         uint   `json:"id"`
	Judul      string `json:"judul" binding:"required"`
	Konten     string `json:"konten" binding:"required"`
	UrlDokumen string `json:"url_dokumen" binding:"required"`
}

func TampilInformasi(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var modelInfo []models.Informasi
	
	hasil := db.Find(&modelInfo)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil tampil data informasi", "data": modelInfo})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal tampil data", "kesalahan": kesalahan.Error()})
	}
}

func TambahInformasi(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataInfo StrukturInformasi
	
	if err := c.ShouldBindJSON(&dataInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	modelInfo := models.Informasi{
		Judul:      dataInfo.Judul,
		Konten:     dataInfo.Konten,
		UrlDokumen: dataInfo.UrlDokumen,
	}

	hasil := db.Create(&modelInfo)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil tambah informasi", "data": modelInfo})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal tambah informasi", "kesalahan": kesalahan.Error()})
	}
}

func UbahInformasi(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataInfo StrukturInformasi
	
	if err := c.ShouldBindJSON(&dataInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	var modelInfo models.Informasi
	db.First(&modelInfo, dataInfo.ID)

	modelInfo.Judul = dataInfo.Judul
	modelInfo.Konten = dataInfo.Konten
	modelInfo.UrlDokumen = dataInfo.UrlDokumen

	hasil := db.Save(&modelInfo)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil ubah informasi", "data": modelInfo})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal ubah informasi", "kesalahan": kesalahan.Error()})
	}
}

func HapusInformasi(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataInfo StrukturInformasi
	
	if err := c.ShouldBindJSON(&dataInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	var modelInfo models.Informasi
	
	hasil := db.Delete(&modelInfo, dataInfo.ID)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil hapus informasi", "data": dataInfo})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal hapus informasi", "kesalahan": kesalahan.Error()})
	}
}
package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
)

type StrukturSuhu struct {
	Id     uint    `json:"id"` // Tambahan untuk target Ubah & Hapus
	Lokasi string  `json:"lokasi" binding:"required"`
	Suhu   float32 `json:"suhu" binding:"required"`
}

func Tampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var modelSuhu []models.Suhu
	hasil := db.Find(&modelSuhu)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil Tampil data", "kesalahan": nil, "data": modelSuhu})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal Tampil Data", "kesalahan": kesalahan.Error(), "data": nil})
	}
}

func Tambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataSuhu StrukturSuhu
	if err := c.ShouldBindJSON(&dataSuhu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	modelSuhu := models.Suhu{
		Lokasi:    dataSuhu.Lokasi,
		Suhu:      dataSuhu.Suhu,
		CreatedAt: time.Now(),
	}

	hasil := db.Create(&modelSuhu)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil tambah data", "kesalahan": nil, "data": modelSuhu})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal Tambah Data", "kesalahan": kesalahan.Error(), "data": modelSuhu})
	}
}


func Ubah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataSuhu StrukturSuhu
	if err := c.ShouldBindJSON(&dataSuhu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	var modelSuhu models.Suhu
	db.First(&modelSuhu, dataSuhu.Id)
	modelSuhu.Lokasi = dataSuhu.Lokasi
	modelSuhu.Suhu = dataSuhu.Suhu

	hasil := db.Save(&modelSuhu)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil ubah data", "kesalahan": nil, "data": modelSuhu})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal ubah Data", "kesalahan": kesalahan.Error(), "data": modelSuhu})
	}
}


func Hapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataSuhu StrukturSuhu
	if err := c.ShouldBindJSON(&dataSuhu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	var modelSuhu models.Suhu
	hasil := db.Delete(&modelSuhu, dataSuhu.Id)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil hapus data", "kesalahan": nil, "data": dataSuhu})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal hapus Data", "kesalahan": kesalahan.Error(), "data": dataSuhu})
	}
}
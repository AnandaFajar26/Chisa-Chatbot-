package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PesanTampil(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"pesan": "Fitur tampil pesan siap"})
}

func PesanTambah(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"pesan": "Fitur tambah pesan siap"})
}

func PesanUbah(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"pesan": "Fitur ubah pesan siap"})
}

func PesanHapus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"pesan": "Fitur hapus pesan siap"})
}
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestKirimWA struct untuk menerima data dari Postman
type RequestKirimWA struct {
	NoTujuan string `json:"no_tujuan" binding:"required"`
	Pesan    string `json:"pesan" binding:"required"`
}

func KirimWA(c *gin.Context) {
	var req RequestKirimWA

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON salah. Pastikan menggunakan key 'no_tujuan' dan 'pesan'."})
		return
	}

	// TODO: Panggil fungsi dari package wa untuk mengirim pesan.
	// Contoh: wa.KirimPesanAPI(req.NoTujuan, req.Pesan)
	// Pastikan Anda membuat fungsi KirimPesanAPI(no, pesan) di folder wa/wa.go

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"pesan":  "Pesan berhasil diterima API (Silakan hubungkan dengan modul WA Anda)",
		"data":   req,
	})
}

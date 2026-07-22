package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
)

// Binding JSON dari Postman (User cuma ngirim Nama, Harga, Jumlah)
type StrukturPesanan struct {
	ID         uint    `json:"id"`
	NamaProduk string  `json:"nama_produk" binding:"required"`
	Harga      float64 `json:"harga" binding:"required"`
	Jumlah     int     `json:"jumlah" binding:"required"`
}

func TampilPesanan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var modelPesanan []models.Pesanan
	db.Find(&modelPesanan)

	c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil tampil data pesanan", "data": modelPesanan})
}

func TambahPesanan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesanan StrukturPesanan

	if err := c.ShouldBindJSON(&dataPesanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}


	totalAwal := dataPesanan.Harga * float64(dataPesanan.Jumlah)
	var diskon float64 = 0
	tipePelanggan := "Regular"

	if totalAwal > 500000 {
		diskon = totalAwal * 0.10 // Diskon 10%
		tipePelanggan = "Gold"
	}
	
	totalHargaAkhir := totalAwal - diskon
	// --------------------------------------------------------

	pesananBaru := models.Pesanan{
		NamaProduk:    dataPesanan.NamaProduk,
		Harga:         dataPesanan.Harga,
		Jumlah:        dataPesanan.Jumlah,
		Diskon:        diskon,
		TotalHarga:    totalHargaAkhir,
		TipePelanggan: tipePelanggan,
	}

	db.Create(&pesananBaru)
	c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil tambah pesanan", "data": pesananBaru})
}

func UbahPesanan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesanan StrukturPesanan

	if err := c.ShouldBindJSON(&dataPesanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	var pesananUbah models.Pesanan
	db.First(&pesananUbah, dataPesanan.ID)

	// Hitung ulang jika harga/jumlah diubah
	totalAwal := dataPesanan.Harga * float64(dataPesanan.Jumlah)
	var diskon float64 = 0
	tipePelanggan := "Regular"

	if totalAwal > 500000 {
		diskon = totalAwal * 0.10
		tipePelanggan = "Gold"
	}

	pesananUbah.NamaProduk = dataPesanan.NamaProduk
	pesananUbah.Harga = dataPesanan.Harga
	pesananUbah.Jumlah = dataPesanan.Jumlah
	pesananUbah.Diskon = diskon
	pesananUbah.TotalHarga = totalAwal - diskon
	pesananUbah.TipePelanggan = tipePelanggan

	db.Save(&pesananUbah)
	c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil ubah pesanan", "data": pesananUbah})
}

func HapusPesanan(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesanan StrukturPesanan

	if err := c.ShouldBindJSON(&dataPesanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	var modelPesanan models.Pesanan
	db.Delete(&modelPesanan, dataPesanan.ID)

	c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil hapus pesanan", "data": dataPesanan})
}
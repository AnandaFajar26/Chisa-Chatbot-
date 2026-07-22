package models

import "gorm.io/gorm"

type Pesanan struct {
	gorm.Model
	NamaProduk    string
	Harga         float64
	Jumlah        int
	Diskon        float64
	TotalHarga    float64
	TipePelanggan string
}
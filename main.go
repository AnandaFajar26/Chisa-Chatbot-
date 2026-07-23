package main

import (
	"log"
	"main/controllers"
	"main/models"
	"main/wa"
	"os"
	"time"

	jwtV3 "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db := koneksi()

	// Automigrate database biar tetep aman
	db.AutoMigrate(&models.Suhu{}, &models.Informasi{}, &models.Pesanan{}, &models.User{}, &models.Dokumen{})
	db.AutoMigrate(&models.Pesan{})

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Konfigurasi JWT Middleware
	key_jwt := os.Getenv("KEY_JWT")
	authMiddleware, err := jwtV3.New(&jwtV3.GinJWTMiddleware{
		Realm:       "fikom UDB",
		Key:         []byte(key_jwt),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "identity",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"identity": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwtV3.ExtractClaims(c)
			return &models.User{
				Username: claims["identity"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals models.User
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwtV3.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			var user models.User
			result := db.Where("username = ? AND password = ?", username, password).First(&user)
			if result.Error == nil {
				return &models.User{
					Username: user.Username,
				}, nil
			}

			return nil, jwtV3.ErrFailedAuthentication
		},
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// Route endpoints publik
	r.POST("/login", authMiddleware.LoginHandler)

	// ---> INI ENDPOINT UTAMA UNTUK UAS KAMU <---
	r.POST("/tanya-gemini", controllers.TanyaGeminiAPI)
	r.POST("/kirim-wa", controllers.KirimWA)

	auth := r.Group("/backend")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/informasi", controllers.TampilInformasi)
		auth.POST("/informasi", controllers.TambahInformasi)
		auth.PUT("/informasi", controllers.UbahInformasi)
		auth.DELETE("/informasi", controllers.HapusInformasi)

		auth.GET("/pesanan", controllers.TampilPesanan)
		auth.POST("/pesanan", controllers.TambahPesanan)
		auth.PUT("/pesanan", controllers.UbahPesanan)
		auth.DELETE("/pesanan", controllers.HapusPesanan)

		auth.GET("/pesan", controllers.PesanTampil)
		auth.POST("/pesan", controllers.PesanTambah)
		auth.PUT("/pesan", controllers.PesanUbah)
		auth.DELETE("/pesan", controllers.PesanHapus)

		auth.POST("/drive", controllers.DriveUpload)
		auth.GET("/drive", controllers.DriveTampil)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8111"
	}

	// Jalankan WhatsApp Bot di background (Goroutine)
	go wa.InitWa(db)
	// Jalankan Server Gin
	log.Println("Server berjalan di port: " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}

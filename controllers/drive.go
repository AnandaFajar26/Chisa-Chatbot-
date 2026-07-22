package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
)

// URL Google Script lu
var urlGoogleScript = "https://script.google.com/macros/s/AKfycbyGKH55vq09DJX9TAO-1vKAt8ywJsRFsLuS3eZkVqQ8idNuxyKqNTTNmtLhZUQrfdcO/exec"

func DriveUpload(c *gin.Context) {
	fileName := c.PostForm("fileName")
	file, errFile := c.FormFile("file")
	
	// Pengaman 1: Cek kalau file gagal masuk dari Postman
	if errFile != nil {
		c.JSON(400, gin.H{"pesan": "File tidak terdeteksi oleh Golang", "error": errFile.Error()})
		return
	}
	
	mimeType := file.Header.Get("Content-Type")

	fileOpen, _ := file.Open()
	defer fileOpen.Close()
	fileData, _ := ioutil.ReadAll(fileOpen)

	data := base64.StdEncoding.EncodeToString(fileData)
	postBody, _ := json.Marshal(map[string]string{
		"fileName": fileName,
		"mimeType": mimeType,
		"data":     data,
	})

	requestBody := bytes.NewBuffer(postBody)
	res, err := http.Post(urlGoogleScript, "application/json; charset=UTF-8", requestBody)

	if err != nil {
		c.JSON(500, gin.H{"kode_error": "ERR-DRIVE", "pesan": "Gagal konek ke Google"})
		return
	}

	hasilBody, _ := ioutil.ReadAll(res.Body)
	hasilString := string(hasilBody)
	res.Body.Close()

	var hasilJson map[string]interface{}
	errJson := json.Unmarshal([]byte(hasilString), &hasilJson)

	// === PENGAMAN 2: CEK BALASAN GOOGLE ===
	// Kalau Google ngebales Error HTML, tampilkan pesannya di Postman biar ketahuan!
	if errJson != nil || hasilJson["filename"] == nil {
		c.JSON(500, gin.H{
			"status": false,
			"pesan": "Google Script menolak! Cek tulisan di bawah ini:",
			"balasan_google": hasilString,
		})
		return
	}
	// ======================================

	// Kalau aman, baru Simpan ke database MySQL
	db := c.MustGet("db").(*gorm.DB)
	dokumenBaru := models.Dokumen{
		NamaDokumen: hasilJson["filename"].(string),
		FileId:      hasilJson["fileId"].(string),
		FileUrl:     hasilJson["fileUrl"].(string),
	}
	hasilDokumen := db.Create(&dokumenBaru)

	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"pesan":     "Berhasil Upload",
		"data":      hasilJson,
		"tersimpan": hasilDokumen.RowsAffected,
	})
}

func DriveTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dokumen []models.Dokumen
	db.Find(&dokumen)
	c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil Tampil", "data": dokumen})
}

func DriveUnduh(c *gin.Context) {
	id := c.Param("id")
	res, err := http.Get(urlGoogleScript + "?id=" + id)

	if err != nil {
		c.JSON(500, gin.H{"status": false, "pesan": "Gagal Unduh"})
		return
	}

	hasilBody, _ := ioutil.ReadAll(res.Body)
	hasilString := string(hasilBody)

	var hasilJson map[string]interface{}
	json.Unmarshal([]byte(hasilString), &hasilJson)

	fileBase64 := hasilJson["file"].(string)
	mimeType := hasilJson["mimeType"].(string)

	fileData, _ := base64.StdEncoding.DecodeString(fileBase64)

	c.Writer.Header().Set("Content-Type", mimeType)
	c.Writer.Write(fileData)
}
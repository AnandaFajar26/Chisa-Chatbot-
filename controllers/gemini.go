package controllers

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Struct untuk menangkap input JSON dari Postman
type RequestGemini struct {
	Pesan string `json:"pesan" binding:"required"`
}

func TanyaGeminiAPI(c *gin.Context) {
	var req RequestGemini
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON salah. Pastikan menggunakan key 'pesan'."})
		return
	}

	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API Key Gemini belum diatur di file .env"})
		return
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Close()

	// Pakai alias 'latest' yang selalu dialokasikan untuk free tier
model := client.GenerativeModel("gemini-flash-latest")

	// SYSTEM PROMPT: AI bertindak sebagai pakar hardware & optimizer game
	systemPrompt := "Kamu adalah Coach Esports dan teknisi hardware profesional. " +
		"Tugasmu khusus menganalisis spek perangkat (HP/PC/Laptop) pengguna dan memberikan " +
		"rekomendasi settingan grafik optimal (Frame Rate, Resolusi, Shadows, Anti-Aliasing) " +
		"untuk game berat seperti Wuthering Waves, PUBG Mobile, dll. Berikan jawaban yang ringkas, " +
		"terstruktur dengan poin-poin, serta berikan satu paragraf catatan optimasi suhu/hardware."

	// Gabungkan instruksi dengan pertanyaan asli dari Postman
	pesanLengkap := systemPrompt + "\n\nSpesifikasi & Pertanyaan Pengguna: " + req.Pesan

	resp, err := model.GenerateContent(ctx, genai.Text(pesanLengkap))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan respon: " + err.Error()})
		return
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		jawabanAI := resp.Candidates[0].Content.Parts[0]
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"reply":  jawabanAI,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Respon AI kosong"})
	}
}
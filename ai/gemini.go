package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func TanyaGemini(pertanyaan string) string {
	apiKey := os.Getenv("GEMINI_API_KEY")
	
	// FIX PALING SAKTI: Pakai alias "gemini-flash-latest" yang kebal dari error versi
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-latest:generateContent?key=" + apiKey

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{{"text": pertanyaan}},
			},
		},
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "Gagal koneksi ke Google."
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != 200 {
		return "Error API: " + string(body)
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Ekstrak teks balasan dari JSON Google
	candidates := result["candidates"].([]interface{})
	content := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	return parts[0].(map[string]interface{})["text"].(string)
}
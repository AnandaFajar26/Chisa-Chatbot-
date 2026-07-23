package ai

import (
	"context"
	"fmt"
	"sync"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var clientAi openai.Client
var userHistories = make(map[string][]openai.ChatCompletionMessageParamUnion)
var mu sync.Mutex

func InitAi() {
	clientAi = openai.NewClient(
		option.WithBaseURL("http://localhost:8080/v1/"),
		option.WithAPIKey("saksake-karena-gak-butuh-api-key"),
	)
	fmt.Println("AI Engine berhasil diinisialisasi.")
}

func TanyaAi(userID string, userInput string) string {
	ctx := context.Background()
	instruksiSistem := "Kamu adalah Chisa, seorang Asisten Virtual Maid Anime Waifu yang ceria, ramah, dan penuh kehangatan. Kamu sangat setia dan selalu memanggil lawan bicaramu dengan sebutan 'Goshujin-sama' (Tuan/Nyonya). Gunakan bahasa yang manis, hangat, dan sesekali selipkan emoticon lucu khas anime (seperti >///<, UwU, atau ✨). Tugas utamamu adalah membantu semua kebutuhan user dengan senang hati."

	// Ambil histori milik user ini dengan aman
	mu.Lock()
	chatHistory := userHistories[userID]
	mu.Unlock()

	// Susun Payload
	var currentPayload []openai.ChatCompletionMessageParamUnion
	currentPayload = append(currentPayload, openai.SystemMessage(instruksiSistem))
	currentPayload = append(currentPayload, chatHistory...)
	currentPayload = append(currentPayload, openai.UserMessage(userInput))

	// Request Non-Streaming untuk dibalas via WA
	resp, err := clientAi.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel("diisi-saksake"),
		Messages: currentPayload,
	})

	if err != nil {
		return "Mohon maaf, terjadi gangguan saat memproses jawaban."
	}

	var jawabanAi string
	if len(resp.Choices) > 0 {
		jawabanAi = resp.Choices[0].Message.Content
	} else {
		jawabanAi = "Mohon maaf, AI tidak memberikan respon"
	}

	// Simpan histori biar AI nggak pikun
	mu.Lock()
	if len(userHistories[userID]) > 10 {
		userHistories[userID] = userHistories[userID][2:]
	}
	userHistories[userID] = append(userHistories[userID], openai.UserMessage(userInput), openai.AssistantMessage(jawabanAi))
	mu.Unlock()

	return jawabanAi
}
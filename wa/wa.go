package wa

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"main/models"

	"gorm.io/gorm"

	"main/ai" // TAMBAHAN 1: Import folder AI lu

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var clientWa *whatsmeow.Client
var DB *gorm.DB // Variabel global buat nyimpen koneksi database

// Fungsi untuk kirim balasan dari database ke WA
func kirimPesanDatabase(IDPenerima types.JID, kode string) {
	var pesan models.Pesan
	// Nyari data di database MySQL yang kodenya sama dengan chat user
	result := DB.Where("kode = ?", kode).First(&pesan)

	// Kalau datanya ketemu (gak ada error), kirim balasannya
	if result.Error == nil {
		clientWa.SendMessage(context.Background(), IDPenerima, &waE2E.Message{
			Conversation: proto.String(pesan.Balasan),
		})
	}
}

// Fungsi bawaan lu buat kirim pesan "halo"
func kirimPesan(IDPenerima types.JID) {
	clientWa.SendMessage(context.Background(), IDPenerima, &waE2E.Message{
		Conversation: proto.String("halo, ini pesan dari bot wa"),
	})
}

// TAMBAHAN 2: Fungsi baru buat kirim pesan teks fleksibel (khusus AI)
func kirimPesanText(IDPenerima types.JID, isiPesan string) {
	clientWa.SendMessage(
		context.Background(),
		IDPenerima,
		&waE2E.Message{
			Conversation: proto.String(isiPesan),
		},
	)
}

// Fungsi buat nangkep pesan masuk
func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:

		fmt.Println("Received a message!", v.Message.GetConversation())
		fmt.Println(" => dari kita sendiri = ", v.Info.IsFromMe)
		fmt.Println(" => server = ", v.Info.MessageSource.Chat.Server)
		fmt.Println(" => apakah group = ", v.Info.IsGroup)
		fmt.Println(" => apakah broadcast = ", v.Info.IsIncomingBroadcast())

		// Filter biar cuma merespon chat pribadi (bukan grup/status)
		if !v.Info.IsFromMe &&
			(v.Info.MessageSource.Chat.Server == "s.whatsapp.net" || v.Info.MessageSource.Chat.Server == "lid") &&
			!v.Info.IsGroup &&
			!v.Info.IsIncomingBroadcast() {

			fmt.Println("PENGIRIM = ", v.Info.Sender.User)
			pesanAsli := v.Message.GetConversation() // Simpan huruf aslinya
			fmt.Println("PESAN = ", pesanAsli)

			var id_wa []types.MessageID
			id_wa = append(id_wa, v.Info.ID)

			clientWa.MarkRead(context.Background(), id_wa, time.Now(), v.Info.Chat, v.Info.Sender)
			clientWa.SubscribePresence(context.Background(), v.Info.Sender)

			// Pura-pura ngetik
			clientWa.SendPresence(context.Background(), types.PresenceAvailable)
			time.Sleep(2 * time.Second)
			clientWa.SendChatPresence(context.Background(), v.Info.Sender, types.ChatPresenceComposing, types.ChatPresenceMediaText)
			time.Sleep(2 * time.Second)
			clientWa.SendChatPresence(context.Background(), v.Info.Sender, types.ChatPresencePaused, types.ChatPresenceMediaText)

			// --- TAMBAHAN 3: LOGIKA AUTO REPLY GABUNGAN AI & DATABASE ---
			pesanChat := strings.ToLower(pesanAsli)

			// Kalau pesan diawali dengan "[ai]"
			if strings.HasPrefix(pesanChat, ".ai") {
				pertanyaan := strings.TrimSpace(pesanAsli[4:]) // Potong tulisan "[ai] "-nya
				if pertanyaan != "" {
					jawabanAi := ai.TanyaGemini(pertanyaan)
					kirimPesanText(v.Info.Sender, jawabanAi)
				} else {
					kirimPesanText(v.Info.Sender, "Masukkan pertanyaan setelah prefiks [ai]. Contoh: [ai] Selamat pagi")
				}
			} else if pesanChat == "tes" {
				kirimPesan(v.Info.Sender) // Kalau ngetik "tes", balas halo biasa
			} else {
				kirimPesanDatabase(v.Info.Sender, pesanChat) // Kalau cuma ketik info/prodi, cari di database
			}
		}
	}
}

// InitWa sekarang WAJIB nerima (db *gorm.DB) dari main.go
func InitWa(db *gorm.DB) {

	DB = db // Masukin koneksi DB ke variabel global

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	// TAMBAHAN 4: Trik macOS biar nggak lelet load history WA
	if deviceStore != nil {
		deviceStore.Platform = "macOS"
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)
	clientWa = client

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

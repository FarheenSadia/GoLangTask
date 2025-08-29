package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var (
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID   = os.Getenv("TELEGRAM_CHAT_ID")
)

func sendToTelegram(message string) {
	url := "https://api.telegram.org/bot" + botToken + "/sendMessage"

	body, _ := json.Marshal(map[string]string{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "Markdown",
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Telegram error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Telegram API failed with code:", resp.StatusCode)
	}
}

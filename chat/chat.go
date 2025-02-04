package chat

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Message struct {
	ChatID  string `json:"chat_id"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

var clients = make(map[*websocket.Conn]string)     // Клиенты с их chat_id
var admins = make(map[*websocket.Conn]bool)        // Администраторы
var clientRooms = make(map[string]*websocket.Conn) // Комнаты (chat_id -> клиент)
var broadcast = make(chan Message)                 // Канал для сообщений
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var mutex = &sync.Mutex{}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat_id")
	isAdmin := r.URL.Query().Get("admin") == "true"

	if chatID == "" {
		http.Error(w, "Chat ID required", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка при апгрейде соединения:", err)
		return
	}
	defer ws.Close()

	mutex.Lock()
	if isAdmin {
		admins[ws] = true
	} else {
		clients[ws] = chatID
		if _, ok := clientRooms[chatID]; !ok {
			clientRooms[chatID] = ws // Запоминаем клиента
		}
	}
	mutex.Unlock()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Ошибка при чтении JSON:", err)
			mutex.Lock()
			delete(clients, ws)
			delete(admins, ws)
			if clientRooms[chatID] == ws {
				delete(clientRooms, chatID)
			}
			mutex.Unlock()
			break
		}
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for admin := range admins {
			_ = admin.WriteJSON(msg) // Отправляем сообщение всем админам
		}
		if client, exists := clientRooms[msg.ChatID]; exists {
			_ = client.WriteJSON(msg) // Отправляем клиенту
		}
		mutex.Unlock()
	}
}

func GetActiveChats(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	chatIDs := make([]string, 0, len(clientRooms))
	for chatID := range clientRooms {
		chatIDs = append(chatIDs, chatID)
	}
	json.NewEncoder(w).Encode(chatIDs)
}

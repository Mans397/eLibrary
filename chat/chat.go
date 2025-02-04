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
var chatHistory = make(map[string][]Message)       // История сообщений для каждого чата
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

	// Отправляем историю сообщений при подключении
	if history, ok := chatHistory[chatID]; ok {
		for _, msg := range history {
			_ = ws.WriteJSON(msg) // Отправляем все сообщения в чат
		}
	}

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

		// Сохраняем сообщение в истории
		mutex.Lock()
		chatHistory[chatID] = append(chatHistory[chatID], msg)
		mutex.Unlock()
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

	// Собираем список активных чатов
	activeChats := make([]map[string]interface{}, 0, len(clientRooms))
	for chatID := range clientRooms {
		// Получаем историю сообщений для каждого чата
		history := chatHistory[chatID]

		// Формируем структуру для чата с историей сообщений
		chatData := map[string]interface{}{
			"chat_id":  chatID,
			"messages": history,
		}

		activeChats = append(activeChats, chatData)
	}

	// Отправляем активные чаты с историей сообщений в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activeChats)
}

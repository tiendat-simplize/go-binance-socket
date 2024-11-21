package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ntdat104/go-binance-socket/clients"
	"github.com/ntdat104/go-binance-socket/models"
)

func handleWebSocket(bc *clients.BinanceClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}
		defer conn.Close()

		subscriber := bc.AddSubscriber(conn)

		// Goroutine to handle Binance client messages
		subscriber.Subscribe(func(message string) {
			log.Println("message", message)
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Error writing message to client:", err)
			}
		})

		// Handle incoming client messages
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				// Handle WebSocket close errors
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Println("WebSocket closed by client:", err)
				} else {
					log.Println("Error reading message from client:", err)
				}
				break
			}

			var binanceRequest models.BinanceRequest
			if err := json.Unmarshal(p, &binanceRequest); err != nil {
				log.Println("Error decoding JSON from client:", err)
				continue
			}
			subscriber.SendRequest(binanceRequest)
		}
	}
}

func main() {
	const PORT = "8080"

	bc, err := clients.NewBinanceClient()
	if err != nil {
		log.Println("Fail to connect Binance")
	}
	defer bc.Close()

	http.HandleFunc("/ws", handleWebSocket(bc))
	log.Printf("Websocket server started on: http://localhost:%v/ws", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

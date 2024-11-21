package clients

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ntdat104/go-binance-socket/models"
	"github.com/ntdat104/go-binance-socket/utils"
)

type Subscriber struct {
	Conn        *websocket.Conn
	Id          int64
	Params      *utils.HashSet
	SendRequest func(models.BinanceRequest)
	Handler     func(string)
	Unsubscribe func()
}

type BinanceClient struct {
	Conn       *websocket.Conn
	Params     *utils.HashSet
	Subscriber map[int64]Subscriber
}

func NewBinanceClient() (*BinanceClient, error) {
	const connection_url = "wss://stream.binance.com/stream"
	conn, _, err := websocket.DefaultDialer.Dial(connection_url, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Binance socket is connected!")
	bc := &BinanceClient{
		Conn:       conn,
		Params:     utils.NewHashSet(),
		Subscriber: make(map[int64]Subscriber),
	}
	go bc.ReadMessage()
	return bc, nil
}

func (bc *BinanceClient) Close() error {
	log.Println("Closing Binance socket connection.")
	return bc.Conn.Close()
}

func (bc *BinanceClient) AddSubscriber(conn *websocket.Conn) Subscriber {
	id := time.Now().UnixMilli()
	bc.Subscriber[id] = Subscriber{
		Conn:        conn,
		Id:          id,
		Params:      utils.NewHashSet(),
		SendRequest: bc.WriteMessage,
		Handler:     nil,
		Unsubscribe: func() {
			bc.RemoveSubscriber(id)
		},
	}
	return bc.Subscriber[id]
}

func (bc *BinanceClient) RemoveSubscriber(id int64) {
	delete(bc.Subscriber, id)
}

func (s *Subscriber) Subscribe(callback func(string)) {
	s.Handler = callback
}

func (bc *BinanceClient) ReadMessage() {
	for {
		_, payload, err := bc.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v\n", err)
		}
		// log.Println("Received message:", string(payload))
		for id := range bc.Subscriber {
			if bc.Subscriber[id].Handler != nil {
				bc.Subscriber[id].Handler(string(payload))
			}
		}
	}
}

func (bc *BinanceClient) WriteMessage(request models.BinanceRequest) {
	data, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshalling request: %v\n", err)
		return
	}
	err = bc.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("Error writing WebSocket message: %v\n", err)
		return
	}
	switch request.Method {
	case "SUBSCRIBE":
		bc.Params.AddList(request.Params)
	case "UNSUBSCRIBE":
		bc.Params.RemoveList(request.Params)
	default:
		log.Printf("Unknown method: %s\n", request.Method)
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

var RC *redis.Client

func init() {
	RC = newClient()
}

func main() {
	subscription()
}

func subscription() {
	//can also use websocket.DefaultDialer
	dialer := websocket.Dialer{}
	//request
	connect, _, err := dialer.Dial("wss://stream.yshyqxx.com/stream?streams=btcusdt@aggTrade", nil)
	if err != nil {
		log.Println(err)
		return
	}
	//close connect when leave scope
	defer connect.Close()

	// loop for get message from server
	for {
		//read data from websocket
		//messageType is the type of message
		//messageData is the data of message
		messageType, messageData, err := connect.ReadMessage()
		if nil != err {
			log.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage: //text
			fmt.Println(string(messageData))

			err := RC.Set("streams=btcusdt@aggTrade", string(messageData), 0).Err()
			if err != nil {
				panic(err)
			}

		case websocket.BinaryMessage: //binary
			fmt.Println(messageData)
		case websocket.CloseMessage:
		case websocket.PingMessage:
		case websocket.PongMessage:
		default:
		}
	}

}

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	log.Println(pong)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

package model

import "github.com/gorilla/websocket"

type User struct {
	Username string `json:"username"`
	// ... other user fields
}

type Message struct {
	Sender_User User   `json:"sender_user"`
	Content     string `json:"content"`
	Timestamp   int64  `json:"timestamp"`
}

type Connection struct {
	Ws   *websocket.Conn
	User User
	Room string
}

var Connections = make(map[string][]Connection)
var UserConnection = make(map[string]*Connection)

func Init_Conn(ws *websocket.Conn, user User, room string) *Connection {
	return &Connection{
		Ws:   ws,
		User: user,
		Room: room,
	}
}

func MessageFunc() map[string][]Message {
	chatMessage := make(map[string][]Message)
	return chatMessage
}

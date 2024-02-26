package main

import (
	"chat/model"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

	g.GET("/chat/:roomID", func(c *gin.Context) {

		roomId := c.Param("roomID")
		ws, err := model.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		user := model.User{Username: "ybb"}
		conn := model.Init_Conn(ws, user, roomId)
		model.Connections[roomId] = append(model.Connections[roomId], *conn)
		go handleConnection(c, conn, roomId)
	})

	g.GET("/chat/user/:username", func(ctx *gin.Context) {
		userR := ctx.Param("username")
		ws, err := model.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// if !ok {
		// 	ctx.JSON(404, gin.H{"error": "user can't receive"})
		// 	return
		// }

		//this is who sender
		userS := model.User{Username: ctx.GetHeader("Username")}
		//userS := model.User{Username: "ybb"}
		//usermodel := model.User{Username: userS.Username}

		//ready to conn
		Sconn := model.Init_Conn(ws, userS, "")
		model.UserConnection[userS.Username] = Sconn

		//start chat
		//content := fmt.Sprintf("Welcome to %s's chat room!", userR)
		//welcomeMessage := model.Message{Sender_User: model.User{Username: userR}, Content: content}
		//err = conn.Ws.WriteJSON(welcomeMessage)
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		go handleChatWithSomeone(Sconn, userR, ctx)
	})

	g.Run(":8080")
}

func handleChatWithSomeone(Sconn *model.Connection, userR string, ctx *gin.Context) {
	// Sconn.Ws.WriteJSON("s")
	var Rconn *model.Connection
	defer func() {
		delete(model.UserConnection, Sconn.User.Username)
		Sconn.Ws.Close()
		delete(model.UserConnection, Rconn.User.Username)
		Rconn.Ws.Close()
	}()
	for {
		var msg model.Message
		err := Sconn.Ws.ReadJSON(&msg)
		if err != nil {
			log.Println(",", err)
			break
		}

		// Add timestamp to message
		msg.Timestamp = time.Now().UnixMilli()

		Rconn, ok := model.UserConnection[userR]
		if !ok {
			log.Println("not online")
			continue
		}
		// recipientConn, ok := model.UserConnection[msg.Sender_User.Username]
		// if !ok {
		// 	log.Println("not online")
		// }
		Rconn.Ws.WriteJSON(msg.Content)
	}
}

func handleConnection(c *gin.Context, conn *model.Connection, roomId string) {
	s := fmt.Sprintf("Welcome to %s's chat room!", conn.User.Username)
	conn.Ws.WriteJSON(s)
	defer conn.Ws.Close()
	defer func() {
		// Remove connection from map
		// var msg model.Message
		// err := conn.Ws.ReadJSON(&msg)
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		for i, otherConn := range model.Connections[conn.Room] {
			if otherConn.Ws == conn.Ws {
				// otherConn.Ws.WriteJSON(msg)
				model.Connections[conn.Room] = append(model.Connections[conn.Room][:i], model.Connections[conn.Room][i+1:]...)
				break
			}
		}
	}()

	// ... handle incoming messages from this connection

	// Read and process incoming messages

	// Send messages to other connected users
	for {
		var msg model.Message
		err := conn.Ws.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			break
		}

		for _, otherConn := range model.Connections[roomId] {
			if otherConn.Ws != conn.Ws {
				err = otherConn.Ws.WriteJSON(msg.Content)
				if err != nil {
					log.Println("to other err", err)
				}
			}
		}
	}

}

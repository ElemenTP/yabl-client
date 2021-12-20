package lib

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type MsgStruct struct {
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content"`
}

//client serving function.
func ClientRoutine(conn *websocket.Conn) {
	//detect ctrl + c action so that can close websocket connection normally.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//a channel to connct all goroutines, use to close the program.
	done := make(chan int)

	//a goroutine to receive and show server-sent message.
	go func() {
		defer func() {
			done <- 1
		}()
		for {
			conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(600))) //close connection after 600s
			recvmsg := new(MsgStruct)
			err := conn.ReadJSON(&recvmsg)
			if err != nil {
				if netErr, ok := err.(net.Error); ok {
					if netErr.Timeout() {
						fmt.Print("websocket receive message from " + conn.RemoteAddr().String() + " timeout")
						return
					}
				}
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					fmt.Printf("websocket receive message from %v error: %v \n", conn.RemoteAddr(), err)
					return
				}
				return
			}
			fmt.Println(time.Unix(recvmsg.Timestamp, 0).Format(time.RFC3339), recvmsg.Content)
		}
	}()

	//a goroutine to send messages read from stdin.
	go func() {
		defer func() {
			done <- 1
		}()
		for {
			conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(600))) //close connection after 600s
			sendmsg := new(MsgStruct)
			sendmsg.Timestamp = time.Now().Unix()
			fmt.Scanln(&sendmsg.Content)
			err := conn.WriteJSON(&sendmsg)
			if err != nil {
				if netErr, ok := err.(net.Error); ok {
					if netErr.Timeout() {
						fmt.Print("websocket send message to " + conn.RemoteAddr().String() + " timeout")
						return
					}
				}
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					fmt.Printf("websocket send message to %v error: %v \n", conn.RemoteAddr(), err)
					return
				}
				return
			}
		}
	}()

	//wait for the done channel and ctrl + c action, then return the program.
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("Got interrupt message, existing...")
			return
		}
	}
}

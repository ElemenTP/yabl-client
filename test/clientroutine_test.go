package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"
	"yabl-client/lib"

	"github.com/gorilla/websocket"
)

var addr string

func init() {
	addr = "127.0.0.1:34567"
	http.HandleFunc("/ws", ServeHttp)
	go http.ListenAndServe(addr, nil)
}

func ServeHttp(rw http.ResponseWriter, r *http.Request) {
	Upgrade := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			if r.Method != "GET" {
				return false
			}
			if r.URL.Path != "/ws" {
				return false
			}
			return true
		}}
	conn, err := Upgrade.Upgrade(rw, r, nil)
	if err != nil {
		return
	}
	defer func() {
		recover()
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		<-time.After(time.Second)
	}()

	msg := new(lib.MsgStruct)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(600))) //close connection after 600s
		conn.ReadJSON(&msg)
		<-time.After(time.Second)
		conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(600))) //close connection after 600s
		conn.WriteJSON(&msg)
	}
}

func TestClient(t *testing.T) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+addr+"/ws", nil)
	if err != nil {
		fmt.Println("can not connect")
		t.Fail()
	}
	defer func() {
		recover()
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		<-time.After(time.Second)
	}()
	fmt.Println("got connection")
	lib.ClientRoutine(conn)
}

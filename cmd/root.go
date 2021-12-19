/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

type MsgStruct struct {
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yabl-client",
	Short: "a cli client for yabl interpreter server.",
	Long:  "A yabl chat client in go, using cli as interface.",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		laddr, err := flags.GetString("address")
		if err != nil {
			log.Fatalln(err)
		}
		port, err := flags.GetString("port")
		if err != nil {
			log.Fatalln(err)
		}
		u := url.URL{Scheme: "ws", Host: laddr + ":" + port, Path: "/ws"}
		fmt.Println("Connecting to server", u.String())

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatalln("dial:", err)
		}
		defer func() {
			recover()
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			<-time.After(time.Second)
			fmt.Println("Disconnected from server", u.String())
		}()

		done := make(chan struct{})

		go func() {
			defer close(done)
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

		go func() {
			defer close(done)
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

		for {
			select {
			case <-done:
				return
			case <-interrupt:
				fmt.Println("interrupt")
				return
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	flags := rootCmd.Flags()
	flags.StringP("address", "a", "127.0.0.1", "connect address")
	flags.StringP("port", "p", "8080", "connect port")
}

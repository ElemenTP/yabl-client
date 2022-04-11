/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"net/url"
	"time"
	"yabl-client/lib"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yabl-client",
	Short: "a cli client for yabl interpreter server.",
	Long:  "A yabl chat client in go, using cli as interface.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		//parse adress and port from flags
		flags := cmd.Flags()
		laddr, err := flags.GetString("address")
		if err != nil {
			log.Fatalln(err)
		}
		port, err := flags.GetString("port")
		if err != nil {
			log.Fatalln(err)
		}

		//generate url from address and port.
		u := url.URL{Scheme: "ws", Host: laddr + ":" + port, Path: "/ws"}
		fmt.Println("Connecting to server", u.String())

		//connect yabl server
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatalln("dial:", err)
		}
		defer func() {
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			<-time.After(time.Second)
			fmt.Println("Disconnected from server", u.String())
		}()

		//run client function
		lib.ClientRoutine(conn)
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

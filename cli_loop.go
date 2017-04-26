package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wsxiaoys/terminal/color"
)

func cliLoop(c *websocket.Conn) {
	for {
		var cmd, arg1, arg2 string

		fmt.Printf("> ")
		fmt.Scanln(&cmd, &arg1, &arg2)

		switch cmd {
		case "print":
			msg, err := mc.Message(arg1, env)
			if err != nil {
				color.Printf("@r%s\n@|> ", err)
			} else {
				color.Printf("%s\n> ", msg)
			}
		case "send":
			if err := mc.SendMessage(arg1, env, c); err != nil {
				color.Printf("@r%s\n@|> ", err)
			} else {
				color.Printf("@gSent!\n> ")
			}
		case "env":
			if arg2 == "" {
				color.Printf("%s", env.Get(arg1))
			} else {
				env.Set(arg1, arg2)
			}
		}
	}
}

func listenInterrupt(c *websocket.Conn, interrupt chan os.Signal, done chan struct{}) {
	defer os.Exit(0)
	defer c.Close()

	for {
		select {
		case <-interrupt:
			color.Printf("\n@rInterrupt")

			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}

			return
		}
	}
}

func listenMessages(c *websocket.Conn, done chan struct{}) {
	defer c.Close()
	defer close(done)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			color.Printf("\n@r%s\n", err)
			return
		}
		color.Printf("\n@y< recv:\n%s\n@|> ", message)
	}
}

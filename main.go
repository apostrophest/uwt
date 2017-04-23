package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mikesss/uwt/uwt"
	"github.com/wsxiaoys/terminal/color"
)

var (
	addr = flag.String("addr", "localhost:8080", "http service address")
	dir  = flag.String("dir", "messages/", "messages directory")
	env  = uwt.NewEnvironment()
	mc   = uwt.NewMessageCollection()
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	done := make(chan struct{})

	go func() {
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
	}()

	go func() {
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
	}()

	mc.Load(*dir)

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

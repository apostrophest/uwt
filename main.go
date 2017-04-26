package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/mikesss/uwt/uwt"
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

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	done := make(chan struct{})

	go listenMessages(c, done)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go listenInterrupt(c, interrupt, done)

	mc.Load(*dir)

	cliLoop(c)
}

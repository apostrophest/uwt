package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wsxiaoys/terminal/color"
)

var (
	addr = flag.String("addr", "localhost:8080", "http service address")
	dir  = flag.String("dir", "messages/", "messages directory")
	env  = make(map[string]string)
	t    = template.New("")
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
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	go func() {
		for {
			select {
			case <-interrupt:
				log.Println("interrupt")
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
				c.Close()
				os.Exit(0)
				return
			}
		}
	}()

	loadTemplates()

	for {
		var cmd, arg1, arg2 string
		fmt.Printf("> ")
		fmt.Scanln(&cmd, &arg1, &arg2)

		switch cmd {
		case "print":
			PrintMessage(arg1)
		case "send":
			SendMessage(arg1, c)
		case "env":
			if arg2 == "" {
				PrintEnv(arg1)
			} else {
				SetEnv(arg1, arg2)
			}
		}
	}

}

func loadTemplates() {
	pattern := filepath.Join(*dir, "*")

	t.ParseGlob(pattern)

	fmt.Printf("Loaded %d templates\n", len(t.Templates()))
}

func PrintEnv(varname string) {
	fmt.Printf("%s: %s\n", varname, env[varname])
}

func SetEnv(varname string, value string) {
	env[varname] = value
}

func PrintMessage(name string) {
	err := t.ExecuteTemplate(os.Stdout, name, env)
	if err != nil {
		color.Printf("@r%s\n", err)
	}
}

func SendMessage(name string, c *websocket.Conn) {
	w, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		color.Printf("@r%s\n", err)
		return
	}

	err = t.ExecuteTemplate(w, name, env)
	if err != nil {
		color.Printf("@r%s\n", err)
		return
	}

	if err := w.Close(); err != nil {
		color.Printf("@r%s\n", err)
	}

	color.Println("@{g}Sent!")
}

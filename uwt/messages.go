package uwt

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/gorilla/websocket"
)

// MessageCollection is a collection of message templates that can be sent to a websocket connection
type MessageCollection struct {
	t *template.Template
}

// NewMessageCollection returns a new, empty collection
func NewMessageCollection() *MessageCollection {
	return &MessageCollection{template.New("")}
}

// Load loads the messages in the given path
func (m *MessageCollection) Load(path string) {
	pattern := filepath.Join(path, "*")
	m.t.ParseGlob(pattern)
}

// Message renders the message against the given environment and returns the resulting string
func (m *MessageCollection) Message(name string, env *Environment) (string, error) {
	b := new(bytes.Buffer)

	err := m.t.ExecuteTemplate(b, name, env)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// SendMessage sends a message to a websocket connection
func (m *MessageCollection) SendMessage(name string, env *Environment, c *websocket.Conn) error {
	w, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	err = m.t.ExecuteTemplate(w, name, env)
	if err != nil {
		return err
	}

	return w.Close()
}

// Num returns the number of messages in the collection
func (m *MessageCollection) Num() int {
	return len(m.t.Templates())
}

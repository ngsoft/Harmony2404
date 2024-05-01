package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func IsSocket(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal("ERR: ", err)
	}
	return fileInfo.Mode().Type() == fs.ModeSocket
}

func hexToInt(hexString string) int {
	i, err := strconv.ParseInt("0x"+hexString, 0, 64)
	if err != nil {
		return 0
	}
	return int(i)
}

func generateUid() string {
	return uuid.NewV4().String()
}

func (h *Handler) readMessage() (string, error) {
	var (
		mt    int
		msg   string
		input []byte
		err   error = nil
	)
	err = fmt.Errorf("no websocket connection defined")
	if h.WebSocket.Open {
		h.Mutex.Read.Lock()
		defer h.Mutex.Read.Unlock()
		if mt, input, err = h.WebSocket.Conn.ReadMessage(); err == nil {
			if mt == websocket.BinaryMessage {
				err = fmt.Errorf("binary message received")
			} else {
				msg = strings.TrimRight(string(input), "\r\n")
			}
		}

	}
	return msg, err
}

func (h *Handler) writeMessage(message string) error {
	if !h.WebSocket.Open {
		return fmt.Errorf("no websocket connection defined")
	}
	h.Mutex.Write.Lock()
	defer h.Mutex.Write.Unlock()
	return h.WebSocket.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (h *Handler) writeJson(v interface{}) error {
	if !h.WebSocket.Open {
		return fmt.Errorf("no websocket connection defined")
	}
	h.Mutex.Write.Lock()
	defer h.Mutex.Write.Unlock()
	return h.WebSocket.Conn.WriteJSON(v)
}

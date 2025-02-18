package platformapiv1

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *handler) HandleMessagePub(ctx context.Context, data []byte) error {
	message := string(data)

	if err := h.realtimeService.PublishMessage(ctx, message); err != nil {
		log.Printf("error to publish the msg: %v\n", err)
		return err
	}

	return nil
}

func (h *handler) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	connectionID := uuid.New().String()

	h.mu.Lock()
	h.connections[connectionID] = conn
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.connections, connectionID)
		h.mu.Unlock()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
	}

}

func (h *handler) broadcastMessages() error {
	ctx := context.Background() // todo use app context

	err := h.realtimeService.SubscribeToMessages(ctx, func(message string) error {
		h.mu.Lock()
		defer h.mu.Unlock()

		for connectionID, conn := range h.connections {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				delete(h.connections, connectionID)
				return err
			}
		}

		return nil
	})

	return err
}

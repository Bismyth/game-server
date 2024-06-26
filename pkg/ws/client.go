package ws

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Bismyth/game-server/pkg/api"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     AllowedOriginCheck,
}

// TODO: add array to config
var allowedOrigins = []string{
	"localhost:8080",
	"bismyth.github.io",
	"met4000.github.io",
}

func AllowedOriginCheck(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}

	if strings.EqualFold(u.Host, r.Host) {
		return true
	}

	for _, origin := range allowedOrigins {
		if strings.EqualFold(origin, u.Host) {
			return true
		}
	}

	return false
}

type CloseMessage struct {
	Code    int
	Message string
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	sessionId uuid.UUID

	verified bool

	leaveMessage *CloseMessage
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if !c.verified {
			verifiedError := (func() error {
				claims, err := api.VerifyRoomToken(string(message))
				if err != nil {
					return fmt.Errorf("invalid initilization packet received")
				}
				err = api.HandleSessionInit(c.hub, claims, c.sessionId)
				if err != nil {
					return fmt.Errorf("failed to initialize connection")
				}
				return nil
			})()
			if verifiedError != nil {
				errorPacket := api.CreateErrorPacket(verifiedError)
				c.conn.WriteMessage(websocket.TextMessage, api.MarshalPacket(&errorPacket))
				c.leaveMessage = &roomLeave
				c.hub.unregister <- c
				continue
			}

			c.verified = true

			// TODO: Check for duplicate connection
			// c.send <- api.UserInitPacket(c.id)
			continue
		}

		iPacket := api.IRawMessage{
			Message:  message,
			SessonId: c.sessionId,
		}

		c.hub.broadcast <- &iPacket
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				message := []byte{}
				if c.leaveMessage != nil {
					message = websocket.FormatCloseMessage(c.leaveMessage.Code, c.leaveMessage.Message)
				}

				c.conn.WriteMessage(websocket.CloseMessage, message)
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), sessionId: id}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

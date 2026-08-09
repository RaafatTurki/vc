package signal

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/raafat/vivid/internal/config"
)

const (
	readLimit   = 64 << 10
	writeWait   = 10 * time.Second
	pongWait    = 60 * time.Second
	pingPeriod  = 50 * time.Second
	messageSize = 256
)

type Server struct {
	cfg      config.Config
	hub      *Hub
	logger   *slog.Logger
	upgrader websocket.Upgrader
	now      func() time.Time
}

func NewServer(cfg config.Config, hub *Hub, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	s := &Server{cfg: cfg, hub: hub, logger: logger, now: time.Now}
	s.upgrader = websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin:      s.originAllowed,
	}
	return s
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /v1/ws", s.websocket)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	rooms, peers := s.hub.Counts()
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"rooms":  rooms,
		"peers":  peers,
	})
}

func (s *Server) websocket(w http.ResponseWriter, r *http.Request) {
	roomID := strings.TrimSpace(r.URL.Query().Get("room"))
	if !validRoomID(roomID) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "room must be exactly 6 alphanumeric characters",
		})
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Debug("websocket upgrade rejected", "error", err)
		return
	}

	peerID, err := randomID()
	if err != nil {
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "could not create peer ID"),
			time.Now().Add(writeWait))
		_ = conn.Close()
		return
	}

	c := &client{
		peerID: peerID,
		roomID: roomID,
		conn:   conn,
		send:   make(chan ServerMessage, messageSize),
	}
	iceServers := makeICEServers(s.cfg, peerID, s.now())
	err = s.hub.Join(c, ServerMessage{
		Type:       "welcome",
		RoomID:     roomID,
		PeerID:     peerID,
		ICEServers: iceServers,
	})
	if err != nil {
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, err.Error()),
			time.Now().Add(writeWait))
		_ = conn.Close()
		return
	}
	s.logger.Info("peer connected", "room", roomID, "peer", peerID)

	go c.writePump(s.logger)
	c.readPump(s.hub)

	s.hub.Leave(c)
	close(c.send)
	_ = conn.Close()
	s.logger.Info("peer disconnected", "room", roomID, "peer", peerID)
}

func (s *Server) originAllowed(r *http.Request) bool {
	origin := strings.TrimSuffix(r.Header.Get("Origin"), "/")
	if origin == "" {
		return true
	}
	for _, allowed := range s.cfg.AllowedOrigins {
		if allowed == "*" || strings.TrimSuffix(allowed, "/") == origin {
			return true
		}
	}
	if len(s.cfg.AllowedOrigins) != 0 {
		return false
	}
	parsed, err := url.Parse(origin)
	return err == nil && strings.EqualFold(parsed.Host, r.Host)
}

func validRoomID(value string) bool {
	if len(value) != 6 {
		return false
	}
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			continue
		}
		return false
	}
	return true
}

func randomID() (string, error) {
	buffer := make([]byte, 12)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

type client struct {
	peerID string
	roomID string
	conn   *websocket.Conn
	send   chan ServerMessage
}

func (c *client) trySend(message ServerMessage) bool {
	select {
	case c.send <- message:
		return true
	default:
		return false
	}
}

func (c *client) readPump(hub *Hub) {
	c.conn.SetReadLimit(readLimit)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		var message ClientMessage
		if err := c.conn.ReadJSON(&message); err != nil {
			return
		}
		if !message.validRelay() {
			c.trySend(ServerMessage{Type: "error", Code: "invalid-message", Message: "expected peer metadata or a WebRTC signaling message with a target and payload"})
			continue
		}
		if message.To == c.peerID {
			c.trySend(ServerMessage{Type: "error", Code: "invalid-target", Message: "cannot relay a message to yourself"})
			continue
		}
		if err := hub.Relay(c, message); err != nil {
			code := "relay-failed"
			if errors.Is(err, ErrPeerMissing) {
				code = "peer-not-found"
			}
			c.trySend(ServerMessage{Type: "error", Code: code, Message: err.Error()})
		}
	}
}

func (c *client) writePump(logger *slog.Logger) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				_ = c.conn.Close()
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
				logger.Debug("websocket ping failed", "peer", c.peerID, "error", err)
				_ = c.conn.Close()
				return
			}
		}
	}
}

package signal

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
)

const (
	messageOffer     = "offer"
	messageAnswer    = "answer"
	messageCandidate = "ice-candidate"
	messagePeerReady = "peer-ready"
	messagePeerState = "peer-state"
	messageChat      = "chat-message"
	maxChatLength    = 4000
	maxChatBytes     = 16 << 10
	maxChatHistory   = 500
	maxChatName      = 80
)

type ClientMessage struct {
	Type    string          `json:"type"`
	To      string          `json:"to"`
	Payload json.RawMessage `json:"payload"`
}

type ServerMessage struct {
	Type        string          `json:"type"`
	RoomID      string          `json:"roomId,omitempty"`
	PeerID      string          `json:"peerId,omitempty"`
	From        string          `json:"from,omitempty"`
	Peers       []string        `json:"peers,omitempty"`
	Payload     json.RawMessage `json:"payload,omitempty"`
	ICEServers  []ICEServer     `json:"iceServers,omitempty"`
	ChatHistory []ChatRecord    `json:"chatHistory,omitempty"`
	Code        string          `json:"code,omitempty"`
	Message     string          `json:"message,omitempty"`
}

type ICEServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

func (m ClientMessage) validRelay() bool {
	if len(m.Payload) == 0 || !json.Valid(m.Payload) {
		return false
	}
	if m.Type == messageChat {
		return true
	}
	if m.To == "" {
		return false
	}
	switch m.Type {
	case messageOffer, messageAnswer, messageCandidate, messagePeerReady, messagePeerState:
		return true
	default:
		return false
	}
}

type ChatRecord struct {
	From    string          `json:"from"`
	Payload json.RawMessage `json:"payload"`
}

type chatPayload struct {
	Text       string `json:"text"`
	SenderName string `json:"senderName"`
	Timestamp  int64  `json:"timestamp"`
}

func validChatPayload(payload json.RawMessage, timestamp int64) (json.RawMessage, bool) {
	var message chatPayload
	if json.Unmarshal(payload, &message) != nil || message.Text == "" || !utf8.ValidString(message.Text) {
		return nil, false
	}
	message.SenderName = strings.TrimSpace(message.SenderName)
	message.Timestamp = timestamp
	if message.SenderName == "" || !utf8.ValidString(message.SenderName) || utf8.RuneCountInString(message.SenderName) > maxChatName {
		return nil, false
	}
	if utf8.RuneCountInString(message.Text) > maxChatLength || len(message.Text) > maxChatBytes {
		return nil, false
	}
	canonical, err := json.Marshal(message)
	return canonical, err == nil
}

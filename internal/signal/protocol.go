package signal

import "encoding/json"

const (
	messageOffer     = "offer"
	messageAnswer    = "answer"
	messageCandidate = "ice-candidate"
	messagePeerReady = "peer-ready"
	messagePeerState = "peer-state"
)

type ClientMessage struct {
	Type    string          `json:"type"`
	To      string          `json:"to"`
	Payload json.RawMessage `json:"payload"`
}

type ServerMessage struct {
	Type       string          `json:"type"`
	RoomID     string          `json:"roomId,omitempty"`
	PeerID     string          `json:"peerId,omitempty"`
	From       string          `json:"from,omitempty"`
	Peers      []string        `json:"peers,omitempty"`
	Payload    json.RawMessage `json:"payload,omitempty"`
	ICEServers []ICEServer     `json:"iceServers,omitempty"`
	Code       string          `json:"code,omitempty"`
	Message    string          `json:"message,omitempty"`
}

type ICEServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

func (m ClientMessage) validRelay() bool {
	if m.To == "" || len(m.Payload) == 0 || !json.Valid(m.Payload) {
		return false
	}
	switch m.Type {
	case messageOffer, messageAnswer, messageCandidate, messagePeerReady, messagePeerState:
		return true
	default:
		return false
	}
}

package signal

import (
	"errors"
	"sort"
	"sync"
)

var (
	ErrRoomFull    = errors.New("room is full")
	ErrPeerMissing = errors.New("target peer is not in the room")
	ErrPeerSlow    = errors.New("target peer is not accepting messages")
)

type Hub struct {
	mu           sync.RWMutex
	rooms        map[string]map[string]*client
	maxRoomPeers int
}

func NewHub(maxRoomPeers int) *Hub {
	return &Hub{rooms: make(map[string]map[string]*client), maxRoomPeers: maxRoomPeers}
}

func (h *Hub) Join(c *client, welcome ServerMessage) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	room := h.rooms[c.roomID]
	if room == nil {
		room = make(map[string]*client)
		h.rooms[c.roomID] = room
	}
	if len(room) >= h.maxRoomPeers {
		return ErrRoomFull
	}

	peers := make([]string, 0, len(room))
	for id := range room {
		peers = append(peers, id)
	}
	sort.Strings(peers)
	welcome.Peers = peers
	c.send <- welcome
	room[c.peerID] = c
	for _, peer := range room {
		if peer != c {
			peer.trySend(ServerMessage{Type: "peer-joined", PeerID: c.peerID})
		}
	}
	return nil
}

func (h *Hub) Leave(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room := h.rooms[c.roomID]
	if room == nil || room[c.peerID] != c {
		return
	}
	delete(room, c.peerID)
	for _, peer := range room {
		peer.trySend(ServerMessage{Type: "peer-left", PeerID: c.peerID})
	}
	if len(room) == 0 {
		delete(h.rooms, c.roomID)
	}
}

func (h *Hub) Relay(from *client, message ClientMessage) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	target := h.rooms[from.roomID][message.To]
	if target == nil {
		return ErrPeerMissing
	}
	if !target.trySend(ServerMessage{
		Type:    message.Type,
		From:    from.peerID,
		Payload: message.Payload,
	}) {
		return ErrPeerSlow
	}
	return nil
}

func (h *Hub) Counts() (rooms, peers int) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, room := range h.rooms {
		rooms++
		peers += len(room)
	}
	return rooms, peers
}

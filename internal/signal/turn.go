package signal

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/raafat/vivid/internal/config"
)

func makeICEServers(cfg config.Config, peerID string, now time.Time) []ICEServer {
	servers := make([]ICEServer, 0, 2)
	if len(cfg.StunURLs) > 0 {
		servers = append(servers, ICEServer{URLs: cfg.StunURLs})
	}
	if len(cfg.TurnURLs) == 0 {
		return servers
	}

	username := strconv.FormatInt(now.Add(cfg.TurnTTL).Unix(), 10) + ":" + peerID
	mac := hmac.New(sha1.New, []byte(cfg.TurnSharedSecret))
	_, _ = mac.Write([]byte(username))
	credential := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return append(servers, ICEServer{
		URLs:       cfg.TurnURLs,
		Username:   username,
		Credential: credential,
	})
}

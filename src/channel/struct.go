package channel

import (
	"strings"
	"sync"
	"time"

	"rabbitsky/src/player"

	cmap "github.com/orcaman/concurrent-map"
)

type Channel struct {
	ServerTick int
	MaxPlayers int
	Players    cmap.ConcurrentMap
	LastID     int
	SkyColor   string
	Bots       []*player.Player

	// Helper
	Ticker           *time.Ticker
	BroadcastMessage strings.Builder
	BroadcastMutex   sync.Mutex
}

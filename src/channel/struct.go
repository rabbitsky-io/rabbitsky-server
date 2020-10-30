package channel

import (
	"strings"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

type Channel struct {
	ServerTick int
	MaxPlayers int
	Players    cmap.ConcurrentMap
	LastID     int

	// Helper
	Ticker           *time.Ticker
	BroadcastMessage strings.Builder
	BroadcastMutex   sync.Mutex
}

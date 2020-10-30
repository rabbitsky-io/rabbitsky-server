package httphandler

import (
	"encoding/json"
	"net/http"
)

func (h *HTTPHandler) GetChannelPlayers(w http.ResponseWriter, r *http.Request) {
	players := ChannelPlayersJSON{
		MaxPlayers: h.Channel.GetMaxPlayers(),
		Players:    h.Channel.GetCurrentPlayers(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", h.Origin)
	w.Header().Set("Access-Control-Request-Method", "GET, OPTIONS")
	json.NewEncoder(w).Encode(players)
}

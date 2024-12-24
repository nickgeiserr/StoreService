package models

type Player struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"name"`
	Clubcoin   int    `json:"clubcoin"`
	Gems       int    `json:"gems"`
}

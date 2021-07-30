package config

/* Config struct contains possible default configuration if no parameter was given to the binary */
type Config struct {
	ServerPort      int    `json:"serverPort"`
	ServerTick      int    `json:"serverTick"`
	MaxPlayers      int    `json:"maxPlayers"`
	Origin          string `json:"origin"`
	LimitPosMin     string `json:"limitPositionMin"`
	LimitPosMax     string `json:"limitPositionMax"`
	DefaultBotCount int    `json:"defaultBotCount"`

	/* later to be admin credentials */
	ServerPassword string `json:"adminPassword"`
}

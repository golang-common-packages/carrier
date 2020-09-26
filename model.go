package carrier

// Config model for carrier services
type Config struct {
	Websocket Websocket `json:"like,omitempty"`
}

// Websocket config model
type Websocket struct {
	Scheme  string `json:"scheme"`
	URL     string `json:"url"`
	Channel string `json:"channel"`
}

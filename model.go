package carrier

import "time"

// Config model for carrier services
type Config struct {
	Websocket Websocket `json:"like,omitempty"`
}

// Websocket config model
type Websocket struct {
	Scheme            string        `json:"scheme"`
	URL               string        `json:"url"`
	Channel           string        `json:"channel"`
	RecIntervalMin    time.Duration `json:"recIntervalMin"`    // RecIntervalMin specifies the initial reconnecting interval, example: 2 * time.Second (2 seconds)
	RecIntervalMax    time.Duration `json:"recIntervalMax"`    // RecIntervalMax specifies the maximum reconnecting interval, example: 30 * time.Second (30 seconds)
	RecIntervalFactor float64       `json:"recIntervalFactor"` // RecIntervalFactor specifies the rate of increase of the reconnection interval, example: 0.5 * time.Second (0.5 seconds)
}

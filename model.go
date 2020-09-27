package carrier

import (
	"time"
)

// Config model for carrier services
type Config struct {
	Gorilla    GorillaConfigModel    `json:"gorilla,omitempty"`
	Eclipse    EclipseConfigModel    `json:"eclipse,omitempty"`
	CustomMail CustomMailConfigModel `json:"customMail,omitempty"`
}

// GorillaConfigModel connection config model
type GorillaConfigModel struct {
	Scheme            string        `json:"scheme"`
	URL               string        `json:"url"`
	Channel           string        `json:"channel"`
	RecIntervalMin    time.Duration `json:"recIntervalMin"`    // RecIntervalMin specifies the initial reconnecting interval, example: 2 * time.Second (2 seconds)
	RecIntervalMax    time.Duration `json:"recIntervalMax"`    // RecIntervalMax specifies the maximum reconnecting interval, example: 30 * time.Second (30 seconds)
	RecIntervalFactor float64       `json:"recIntervalFactor"` // RecIntervalFactor specifies the rate of increase of the reconnection interval, example: 0.5 * time.Second (0.5 seconds)
}

// EclipseConfigModel connection config model
type EclipseConfigModel struct {
	ClientID string `json:"clientID"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CustomMailConfigModel connection config model
type CustomMailConfigModel struct {
	PoolSize int    `json:"poolSize"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Identity string `json:"identity"`
	Username string `json:"username"`
	Password string `json:"password"`
}

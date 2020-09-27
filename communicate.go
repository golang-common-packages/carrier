package carrier

// ICOMMUNICATE factory pattern interface
type ICOMMUNICATE interface {
	Write(message string) error
	Read() (interface{}, error)
	IsConnected() bool
	End() error
}

const (
	// WEBSOCKET services
	WEBSOCKET = iota
)

// NewCommunication Factory Pattern
func NewCommunication(carrierCompany int, config *Config) interface{} {

	switch carrierCompany {
	case WEBSOCKET:
		return NewWebsocketClient(&config.Websocket)
	}

	return nil
}

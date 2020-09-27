package carrier

const (
	// WEBSOCKET services
	WEBSOCKET = iota
	// MQTT services
	MQTT
)

// New carrier based on types
func New(carrierType int) func(carrierCompany int, config *Config) interface{} {
	switch carrierType {
	case WEBSOCKET:
		return NewWebsocket
	case MQTT:
		return NewMQTT
	default:
		return nil
	}
}

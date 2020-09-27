package carrier

const (
	// WEBSOCKET services
	WEBSOCKET = iota
	// MQTT services
	MQTT
	// MAIL services
	MAIL
)

// New carrier based on types
func New(carrierType int) func(carrierCompany int, config *Config) interface{} {
	switch carrierType {
	case WEBSOCKET:
		return NewWebsocket
	case MQTT:
		return NewMQTT
	case MAIL:
		return NewMail
	default:
		return nil
	}
}

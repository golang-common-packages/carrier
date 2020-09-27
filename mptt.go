package carrier

import mqtt "github.com/eclipse/paho.mqtt.golang"

// IMQTT factory pattern interface
type IMQTT interface {
	Publish(topic, message string) error
	Subscribe(topic string, messageHandler mqtt.MessageHandler) error
	IsConnected() bool
	End()
}

const (
	// ECLIPSE services
	ECLIPSE = iota
)

// NewMQTT Factory Pattern
func NewMQTT(carrierCompany int, config *Config) interface{} {

	switch carrierCompany {
	case ECLIPSE:
		return NewEclipseClient(&config.Eclipse)
	}

	return nil
}

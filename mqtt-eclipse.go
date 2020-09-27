package carrier

import (
	"log"

	"github.com/eclipse/paho.mqtt.golang"
)

// EclipseClient manage all MQTT client actions
type EclipseClient struct {
	client mqtt.Client
	config *MQTTConnectionModel
}

var (
	// EclipseClientSessionMapping singleton pattern
	EclipseClientSessionMapping = make(map[string]*EclipseClient)
)

// NewEclipseClient init new instance
func NewEclipseClient(config *MQTTConnectionModel) IMQTT {
	configHashed := hashObject(config)

	currentEclipseClientSession := EclipseClientSessionMapping[configHashed]
	if currentEclipseClientSession == nil {
		currentEclipseClientSession = &EclipseClient{nil, nil}

		clientOptions := mqtt.NewClientOptions().AddBroker(config.URL).SetClientID(config.ClientID).SetUsername(config.Username).SetPassword(config.Password)
		client := mqtt.NewClient(clientOptions)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("MQTT client: can't connect to broker of %v\n", token.Error())
		}

		currentEclipseClientSession.client = client
		currentEclipseClientSession.config = config
		EclipseClientSessionMapping[configHashed] = currentEclipseClientSession
		log.Println("MQTT client: connected")
	}

	return currentEclipseClientSession
}

// Publish message to channel
func (ec *EclipseClient) Publish(topic, message string) error {
	if token := ec.client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Subscribe message from channel
func (ec *EclipseClient) Subscribe(topic string, messageHandler mqtt.MessageHandler) error {
	if token := ec.client.Subscribe(topic, 0, messageHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// IsConnected return connection state
func (ec *EclipseClient) IsConnected() bool {
	return ec.client.IsConnected()
}

// End this communication
func (ec *EclipseClient) End() {
	ec.client.Disconnect(1000) // 1 second
}

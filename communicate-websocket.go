package carrier

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/golang-common-packages/hash"
	"github.com/gorilla/websocket"
)

// WebsocketClient manage all websocket client actions
type WebsocketClient struct {
	Connection *websocket.Conn
}

var (
	// websocketClientSessionMapping singleton pattern
	websocketClientSessionMapping = make(map[string]*WebsocketClient)
)

// NewWebsocketClient init new instance
func NewWebsocketClient(config *Websocket) ICOMMUNICATE {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentWebsocketClientSession := websocketClientSessionMapping[configAsString]
	if currentWebsocketClientSession == nil {
		currentWebsocketClientSession = &WebsocketClient{nil}

		url := url.URL{Scheme: config.Scheme, Host: config.URL, Path: config.Channel}

		connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
		if err != nil {
			panic(err)
		}

		currentWebsocketClientSession.Connection = connection
		websocketClientSessionMapping[configAsString] = currentWebsocketClientSession
		log.Println("Connected to Websocket Server")
	}

	return currentWebsocketClientSession
}

// Write message to channel
func (wc *WebsocketClient) Write(message string) error {
	return wc.Connection.WriteMessage(websocket.TextMessage, []byte(message))
}

// Read message from channel
func (wc *WebsocketClient) Read() (interface{}, error) {
	_, message, err := wc.Connection.ReadMessage()
	return message, err
}

// End this communication
func (wc *WebsocketClient) End() error {
	return wc.Connection.Close()
}

package carrier

import (
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
)

// WebsocketClient manage all websocket client actions
type WebsocketClient struct {
	connection       *websocket.Conn
	connectionStatus bool
	config           *Websocket
}

var (
	// websocketClientSessionMapping singleton pattern
	websocketClientSessionMapping = make(map[string]*WebsocketClient)

	// WebsocketClientMutex mutex for this service only
	WebsocketClientMutex sync.Mutex
)

// NewWebsocketClient init new instance
func NewWebsocketClient(config *Websocket) ICOMMUNICATE {
	configHashed := hashObject(config)

	currentWebsocketClientSession := websocketClientSessionMapping[configHashed]
	if currentWebsocketClientSession == nil {
		currentWebsocketClientSession = &WebsocketClient{nil, false, nil}

		url := url.URL{Scheme: config.Scheme, Host: config.URL, Path: config.Channel}
		connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
		if err != nil {
			currentWebsocketClientSession.closeAndReConect()
		} else {
			currentWebsocketClientSession.connection = connection
			currentWebsocketClientSession.connectionStatus = true
			currentWebsocketClientSession.config = config
			websocketClientSessionMapping[configHashed] = currentWebsocketClientSession
			log.Println("Connected to Websocket server")
		}
	}

	return currentWebsocketClientSession
}

// CloseAndRecconect will try to reconnect
func (wc *WebsocketClient) closeAndReConect() {
	wc.connection.Close()
	go func() {
		wc.reConnect()
	}()
}

func (wc *WebsocketClient) reConnect() {
	bo := &backoff.Backoff{
		Min:    wc.config.RecIntervalMin,
		Max:    wc.config.RecIntervalMax,
		Factor: wc.config.RecIntervalFactor,
		Jitter: true,
	}

	rand.Seed(time.Now().UTC().UnixNano())

	for {
		nextInterval := bo.Duration()

		url := url.URL{Scheme: wc.config.Scheme, Host: wc.config.URL, Path: wc.config.Channel}
		newConnection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
		if err != nil {
			log.Printf("Dial: will try again in %v seconds. Because of %v\n", nextInterval, err)
		} else {
			configHashed := hashObject(wc.config)

			WebsocketClientMutex.Lock()
			currentWebsocketClientSession := websocketClientSessionMapping[configHashed]
			currentWebsocketClientSession.connection = newConnection
			currentWebsocketClientSession.connectionStatus = err == nil
			websocketClientSessionMapping[configHashed] = currentWebsocketClientSession
			WebsocketClientMutex.Unlock()

			log.Println("Successfully reconnected to Websocket server")
			break
		}

		time.Sleep(nextInterval)
	}
}

// Write message to channel
func (wc *WebsocketClient) Write(message string) error {
	return wc.connection.WriteMessage(websocket.TextMessage, []byte(message))
}

// Read message from channel
func (wc *WebsocketClient) Read() (interface{}, error) {
	_, message, err := wc.connection.ReadMessage()
	return message, err
}

// IsConnected returns the WebSocket connection status
func (wc *WebsocketClient) IsConnected() bool {
	WebsocketClientMutex.Lock()
	defer WebsocketClientMutex.Unlock()

	return wc.connectionStatus
}

// End this communication
func (wc *WebsocketClient) End() error {
	return wc.connection.Close()
}

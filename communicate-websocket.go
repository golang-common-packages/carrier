package carrier

import (
	"errors"
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
	connection *websocket.Conn
	config     *Websocket
}

var (
	// websocketClientSessionMapping singleton pattern
	websocketClientSessionMapping = make(map[string]*WebsocketClient)

	// websocketClientMutex mutex for this service only
	websocketClientMutex sync.Mutex
)

// NewWebsocketClient init new instance
func NewWebsocketClient(config *Websocket) ICOMMUNICATE {
	configHashed := hashObject(config)

	currentWebsocketClientSession := websocketClientSessionMapping[configHashed]
	if currentWebsocketClientSession == nil {
		currentWebsocketClientSession = &WebsocketClient{nil, nil}

		url := url.URL{Scheme: config.Scheme, Host: config.URL, Path: config.Channel}
		connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
		if err != nil {
			currentWebsocketClientSession.closeAndReConect()
		} else {
			currentWebsocketClientSession.connection = connection
			currentWebsocketClientSession.config = config
			websocketClientSessionMapping[configHashed] = currentWebsocketClientSession
			log.Println("Websocket client: connected")
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
			log.Printf("Websocket client: will try again in %v seconds. Because of %v\n", nextInterval, err)
		} else {
			configHashed := hashObject(wc.config)

			websocketClientMutex.Lock()
			currentWebsocketClientSession := websocketClientSessionMapping[configHashed]
			currentWebsocketClientSession.connection = newConnection
			websocketClientSessionMapping[configHashed] = currentWebsocketClientSession
			websocketClientMutex.Unlock()

			log.Println("Websocket client: reconnected")
			break
		}

		time.Sleep(nextInterval)
	}
}

// Write message to channel
func (wc *WebsocketClient) Write(message string) error {
	if err := wc.connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		wc.closeAndReConect()
		return err
	}

	return nil
}

// Read message from channel
func (wc *WebsocketClient) Read() (interface{}, error) {
	_, message, err := wc.connection.ReadMessage()
	if err != nil {
		wc.closeAndReConect()
		return nil, err
	}

	return bytesToString(message), nil
}

// End this communication
func (wc *WebsocketClient) End() error {
	return wc.connection.Close()
}

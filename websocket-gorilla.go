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

// GorillaClient manage all websocket client actions
type GorillaClient struct {
	connection *websocket.Conn
	config     *WebsocketConnectionModel
}

var (
	// gorrilaClientSessionMapping singleton pattern
	gorrilaClientSessionMapping = make(map[string]*GorillaClient)

	// gorrilaClientMutex mutex for this service only
	gorrilaClientMutex sync.Mutex
)

// NewGorillaClient init new instance
func NewGorillaClient(config *WebsocketConnectionModel) IWebsocket {
	configHashed := hashObject(config)

	currentGorrilaClientSession := gorrilaClientSessionMapping[configHashed]
	if currentGorrilaClientSession == nil {
		currentGorrilaClientSession = &GorillaClient{nil, nil}

		url := url.URL{Scheme: config.Scheme, Host: config.URL, Path: config.Channel}
		connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
		if err != nil {
			currentGorrilaClientSession.closeAndReConect()
		} else {
			currentGorrilaClientSession.connection = connection
			currentGorrilaClientSession.config = config
			gorrilaClientSessionMapping[configHashed] = currentGorrilaClientSession
			log.Println("Websocket client: connected")
		}
	}

	return currentGorrilaClientSession
}

// CloseAndRecconect will try to reconnect
func (gc *GorillaClient) closeAndReConect() {
	gc.connection.Close()
	go func() {
		gc.reConnect()
	}()
}

func (gc *GorillaClient) reConnect() {
	bo := &backoff.Backoff{
		Min:    gc.config.RecIntervalMin,
		Max:    gc.config.RecIntervalMax,
		Factor: gc.config.RecIntervalFactor,
		Jitter: true,
	}

	rand.Seed(time.Now().UTC().UnixNano())

	for {
		nextInterval := bo.Duration()

		url := url.URL{Scheme: gc.config.Scheme, Host: gc.config.URL, Path: gc.config.Channel}
		newConnection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
		if err != nil {
			log.Printf("Websocket client: will try again in %v seconds. Because of %v\n", nextInterval, err)
		} else {
			configHashed := hashObject(gc.config)

			gorrilaClientMutex.Lock()
			currentWebsocketClientSession := gorrilaClientSessionMapping[configHashed]
			currentWebsocketClientSession.connection = newConnection
			gorrilaClientSessionMapping[configHashed] = currentWebsocketClientSession
			gorrilaClientMutex.Unlock()

			log.Println("Websocket client: reconnected")
			break
		}

		time.Sleep(nextInterval)
	}
}

// Write message to channel
func (gc *GorillaClient) Write(message string) error {
	if err := gc.connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		gc.closeAndReConect()
		return err
	}

	return nil
}

// Read message from channel
func (gc *GorillaClient) Read() (interface{}, error) {
	_, message, err := gc.connection.ReadMessage()
	if err != nil {
		gc.closeAndReConect()
		return nil, err
	}

	return bytesToString(message), nil
}

// End this communication
func (gc *GorillaClient) End() error {
	return gc.connection.Close()
}

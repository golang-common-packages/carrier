package carrier

import (
	"log"
	"net/smtp"
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/hashicorp/go-multierror"
)

// CustomMailClient manage all Gomail client actions
type CustomMailClient struct {
	client smtp.Auth
	config *CustomMailConfigModel
}

var (
	// customMailClientSessionMapping singleton pattern
	customMailClientSessionMapping = make(map[string]*CustomMailClient)

	// customMailClientMutex mutex for this service only
	customMailClientMutex sync.Mutex
)

// NewCustomMailClient init new instance
func NewCustomMailClient(config *CustomMailConfigModel) IMail {
	configHashed := hashObject(config)

	currentCustomClientSession := customMailClientSessionMapping[configHashed]
	if currentCustomClientSession == nil {
		currentCustomClientSession = &CustomMailClient{nil, nil}

		currentCustomClientSession.client = smtp.PlainAuth(config.Identity, config.Username, config.Password, config.Host)
		currentCustomClientSession.config = config
		customMailClientSessionMapping[configHashed] = currentCustomClientSession
		log.Println("Custom mail client: Initialization")
	}

	return currentCustomClientSession
}

// Send message to recipient
func (cc *CustomMailClient) Send(subject, body string, recipients ...string) error {
	var errs *multierror.Error
	wp := workerpool.New(cc.config.PoolSize)

	for _, recipient := range recipients {
		recipient := recipient
		wp.Submit(func() {
			if err := smtp.SendMail(cc.config.Host+":"+cc.config.Port, cc.client, cc.config.Username, []string{recipient}, []byte(body)); err != nil {
				customMailClientMutex.Lock()
				errs = multierror.Append(errs, err)
				customMailClientMutex.Unlock()
			}
		})
	}

	wp.StopWait()

	// Return an error if any failed
	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	return nil
}

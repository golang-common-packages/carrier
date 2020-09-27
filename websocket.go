package carrier

// IWebsocket factory pattern interface
type IWebsocket interface {
	Write(message string) error
	Read() (interface{}, error)
	End() error
}

const (
	// GORILLA services
	GORILLA = iota
)

// NewWebsocket Factory Pattern
func NewWebsocket(carrierCompany int, config *Config) interface{} {

	switch carrierCompany {
	case GORILLA:
		return NewGorillaClient(&config.Gorilla)
	}

	return nil
}

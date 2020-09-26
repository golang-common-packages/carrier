package carrier

const (
	// COMMUNICATE services
	COMMUNICATE = iota
)

// New database by abstract factory pattern
func New(carrierType int) func(carrierCompany int, config *Config) interface{} {
	switch carrierType {
	case COMMUNICATE:
		return NewCommunication
	default:
		return nil
	}
}

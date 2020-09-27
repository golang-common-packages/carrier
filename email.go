package carrier

// IMail factory pattern interface
type IMail interface {
	Send(subject, body string, recipients ...string) error
}

const (
	// CUSTOMMAIL services
	CUSTOMMAIL = iota
)

// NewMail Factory Pattern
func NewMail(carrierCompany int, config *Config) interface{} {
	switch carrierCompany {
	case CUSTOMMAIL:
		return NewCustomMailClient(&config.CustomMail)
	}

	return nil
}

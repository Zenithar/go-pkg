package application

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Config defines application metric collector configuration
type Config struct {
	Name     string
	ID       string
	Version  string
	Revision string
}

// Validate configuration
func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		// Name must not be empty and a valid ascii value
		validation.Field(&c.Name, validation.Required),
	)
}

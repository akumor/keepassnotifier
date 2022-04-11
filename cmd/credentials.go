package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Credentials struct {
	Keepass struct {
		// password to use for accessing the keepass database
		DatabasePassword string `json:"database_password"`
	} `yaml:"keepass"`
	SendGrid struct {
		// API key to use for sending a message with sendgrid
		ApiKey string `json:"api_key"`
	} `yaml:"sendgrid"`
}

// NewCredentialsFromFile returns a new decoded Credentials struct
func (c *Credentials) LoadCredentialsFromEnv() error {
	x := "${KEEPASSNOTIFIER_KEEPASS_DATABASE_PASSWORD}"
	val := os.ExpandEnv(x)
	if val != "" {
		c.Keepass.DatabasePassword = val
	}
	x = "${KEEPASSNOTIFIER_SENDGRID_API_KEY}"
	val = os.ExpandEnv(x)
	if val != "" {
		c.SendGrid.ApiKey = val
	}
	return nil
}

// NewCredentialsFromFile returns a new decoded Credentials struct
func (c *Credentials) LoadCredentialsFromFile(credPath string) error {
	// validate credentials file path
	s, err := os.Stat(credPath)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", credPath)
	}

	// Open credentials file
	file, err := os.Open(credPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Init new json decode
	d := json.NewDecoder(file)

	// Start JSON decoding from file
	if err := d.Decode(&c); err != nil {
		return err
	}

	return nil
}

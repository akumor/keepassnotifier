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
}

// NewCredentialsFromFile returns a new decoded Credentials struct
func (c *Credentials) LoadCredentialsFromEnv() error {
	x := "${KEEPASSNOTIFIER_KEEPASS_DATABASE_PASSWORD}"
	val := os.ExpandEnv(x)
	if val != "" {
		c.Keepass.DatabasePassword = val
	}
	return nil
}

// NewCredentialsFromFile returns a new decoded Credentials struct
func (c *Credentials) LoadCredentialsFromFile(credPath string) error {
	// Open config file
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

// ValidateCredPath just makes sure, that the path provided is a file,
// that can be read
func ValidateCredPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

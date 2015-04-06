// Package keyring provides functions for reading and writing passwords
// securely.
package keyring

import (
	"errors"
	"runtime"
)

var ErrUnsupportedOS = errors.New("not supported")

var keyring = map[string]Keyring{}

// Keyring defines the Keyring client available on the platform.
type Keyring interface {
	// Set sets a password for the service and username.
	Set(service, username string, password []byte) error

	// Get returns the password for the service and username if exists.
	Get(service, username string) ([]byte, error)

	// Delete deletes a password belongs to the service and username.
	Delete(service, username string) error
}

// New returns the Keyring client available on the platform.
func New() (Keyring, error) {
	switch runtime.GOOS {
	case "freebsd", "linux", "netbsd", "openbsd":
		return keyring["factotum"], nil
	case "darwin":
		return keyring["darwin"], nil
	}
	return nil, ErrUnsupportedOS
}

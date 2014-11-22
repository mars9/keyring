// Package keyring provides functions for reading and writing passwords
// securely.
package keyring

// Keyring defines the Keyring client available on the platform.
type Keyring interface {
	// Get gets the password for the service and username if exists.
	Get(service, username string) (string, error)

	// Set sets a password for the service and username.
	Set(service, username, password string) error

	// Delete deletes a password belongs to the service and username.
	Delete(service, username string) error
}

// NewKeyring returns the Keyring client available on the platform.
// Currently only supports Factotum
//		http://plan9.bell-labs.com/magic/man2html/4/factotum
func NewKeyring() Keyring { return new(factotum) }

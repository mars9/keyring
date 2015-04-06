// +build freebsd linux netbsd openbsd

package keyring

import (
	"bytes"
	"testing"
)

var testFactotum = false

func TestKeyringFactotum(t *testing.T) {
	if !testFactotum {
		return
	}

	f := &factotum{}
	err := f.Set("test.service", "test.user", []byte("test.password"))
	if err != nil {
		t.Fatalf("set: %v", err)
	}

	password, err := f.Get("test.service", "test.user")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if bytes.Compare(password, []byte("test.password")) != 0 {
		t.Fatalf("get: expected test.password, got %s", password)
	}

	err = f.Delete("test.service", "test.user")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
}

// +build freebsd linux netbsd openbsd

package keyring

import (
	"bytes"
	"fmt"

	"9fans.net/go/plan9"
	"9fans.net/go/plan9/client"
)

func init() { keyring["factotum"] = &factotum{} }

type rpc struct {
	fid *client.Fid
}

func newRPC(name string) (*rpc, error) {
	fsys, err := client.MountService("factotum")
	if err != nil {
		return nil, err
	}
	fid, err := fsys.Open(name, plan9.ORDWR)
	if err != nil {
		return nil, err
	}
	return &rpc{fid: fid}, nil
}

func (c *rpc) Close() error { return c.fid.Close() }

type factotum struct{}

func (f *factotum) Set(service, username string, password []byte) error {
	params := fmt.Sprintf("dom=%s proto=pass role=client", service)
	key := params + fmt.Sprintf(" user=%s !password=%s", username, password)

	ctl, err := newRPC("ctl")
	if err != nil {
		return err
	}
	defer ctl.Close()

	_, err = ctl.fid.Write([]byte("key " + key))
	return err
}

func (f *factotum) Get(service, username string) ([]byte, error) {
	params := fmt.Sprintf("dom=%s proto=pass role=client", service)

	ctl, err := newRPC("rpc")
	if err != nil {
		return nil, err
	}
	defer ctl.Close()

	if _, err = ctl.fid.Write([]byte("start " + params)); err != nil {
		return nil, err
	}

	buf := make([]byte, 4096)
	n, err := ctl.fid.Read(buf)
	if err != nil {
		return nil, err
	}
	if !bytes.HasPrefix(buf, []byte("ok")) {
		return nil, fmt.Errorf("start failed: %s", buf[:n])
	}

	if _, err = ctl.fid.Write([]byte("read")); err != nil {
		return nil, err
	}

	if n, err = ctl.fid.Read(buf); err != nil {
		return nil, err
	}
	if !bytes.HasPrefix(buf, []byte("ok")) {
		return nil, fmt.Errorf("read failed: %s", buf[:n])
	}

	elems := bytes.Split(buf[:n], []byte(" "))
	if len(elems) != 3 {
		return nil, fmt.Errorf("split response failed")
	}
	return elems[2], nil
}

func (f *factotum) Delete(service, username string) error {
	params := fmt.Sprintf("dom=%s proto=pass role=client user=%s",
		service, username)

	ctl, err := newRPC("ctl")
	if err != nil {
		return err
	}
	defer ctl.Close()

	_, err = ctl.fid.Write([]byte("delkey " + params))
	return err
}

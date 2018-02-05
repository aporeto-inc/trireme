package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/aporeto-inc/trireme-lib/common"
)

// Client is an api client structure.
type Client struct {
	addr *net.UnixAddr
}

// NewClient creates a new client.
func NewClient(path string) (*Client, error) {
	addr, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %s", err)
	}

	return &Client{addr: addr}, nil
}

// SendRequest sends a request to the remote.
// TODO: Add retries
func (c *Client) SendRequest(event *common.EventInfo) error {

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.DialUnix("unix", nil, c.addr)
			},
		},
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(event)
	resp, err := httpc.Post("http://unix", "application/json", b)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	errorBuffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Invalid request: %s", err)
	}

	return fmt.Errorf("Invalid request : %s", string(errorBuffer))
}

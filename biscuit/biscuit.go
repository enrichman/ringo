package biscuit

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Client struct {
	filename string
}

func NewClient(filename string) (*Client, error) {
	if _, err := os.Open(filename); err != nil {
		return nil, err
	}
	return &Client{filename: filename}, nil
}

func (c *Client) List() ([]string, error) {
	out, err := exec.Command("biscuit", "list", "-f", c.filename).Output()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

func (c *Client) Get(secret string) (string, error) {
	cmd := exec.Command("biscuit", "get", "-f", c.filename, secret)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return stderr.String(), err
	}

	return string(out), nil
}

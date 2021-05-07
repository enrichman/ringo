package biscuit

import (
	"bytes"
	"os/exec"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

func KmsCallerIdentity() error {
	cmd := exec.Command("biscuit", "kms", "get-caller-identity")
	_, err := handleCommand(cmd)
	return err
}

func List(filename string) ([]string, error) {
	out, err := handleCommand(exec.Command("biscuit", "list", "-f", filename))
	if err != nil {
		return nil, err
	}

	secrets := strings.Split(out, "\n")
	sort.Strings(secrets)

	return secrets, nil
}

func Get(filename, secret string) (string, error) {
	cmd := exec.Command("biscuit", "get", "-f", filename, secret)
	return handleCommand(cmd)
}

func handleCommand(cmd *exec.Cmd) (string, error) {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	return strings.TrimSpace(string(out)), nil
}

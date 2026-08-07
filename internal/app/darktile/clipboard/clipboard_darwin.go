//go:build darwin

package clipboard

import (
	"bytes"
	"os/exec"
)

const (
	pbcopyPath  = "/usr/bin/pbcopy"
	pbpastePath = "/usr/bin/pbpaste"
)

var runCommand = func(path string, input []byte) ([]byte, error) {
	cmd := exec.Command(path)
	if input != nil {
		cmd.Stdin = bytes.NewReader(input)
	}
	return cmd.Output()
}

func Set(text string) error {
	_, err := runCommand(pbcopyPath, []byte(text))
	return err
}

func Get() (string, error) {
	output, err := runCommand(pbpastePath, nil)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

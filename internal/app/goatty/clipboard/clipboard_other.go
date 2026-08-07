//go:build !darwin

package clipboard

import platformclipboard "github.com/d-tsuji/clipboard"

func Set(text string) error {
	return platformclipboard.Set(text)
}

func Get() (string, error) {
	return platformclipboard.Get()
}

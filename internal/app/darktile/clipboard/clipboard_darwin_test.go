//go:build darwin

package clipboard

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDarwinClipboardCommands(t *testing.T) {
	originalRunCommand := runCommand
	t.Cleanup(func() { runCommand = originalRunCommand })

	var calls []struct {
		path  string
		input string
	}
	runCommand = func(path string, input []byte) ([]byte, error) {
		calls = append(calls, struct {
			path  string
			input string
		}{path: path, input: string(input)})
		if path == pbpastePath {
			return []byte("from pasteboard"), nil
		}
		return nil, nil
	}

	require.NoError(t, Set("to pasteboard"))
	pasted, err := Get()
	require.NoError(t, err)
	require.Equal(t, "from pasteboard", pasted)
	require.Equal(t, []struct {
		path  string
		input string
	}{
		{path: pbcopyPath, input: "to pasteboard"},
		{path: pbpastePath},
	}, calls)
}

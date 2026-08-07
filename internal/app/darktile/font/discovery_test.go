package font

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyamel/goatty/internal/app/darktile/packed"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

func TestNormalizeStyleName(t *testing.T) {
	tests := map[string]StyleName{
		"Regular":      StyleRegular,
		"Roman":        StyleRegular,
		"Bold":         StyleBold,
		"Italic":       StyleItalic,
		"Oblique":      StyleItalic,
		"Bold Italic":  StyleBoldItalic,
		"Bold Oblique": StyleBoldItalic,
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			require.Equal(t, expected, normalizeStyleName(input))
		})
	}
}

func TestDiscoverFontsReadsOpenTypeMetadata(t *testing.T) {
	collection, err := opentype.ParseCollection(packed.MesloLGSNFRegularTTF)
	require.NoError(t, err)
	parsed, err := collection.Font(0)
	require.NoError(t, err)
	family, err := preferredFontName(parsed, sfnt.NameIDTypographicFamily, sfnt.NameIDFamily)
	require.NoError(t, err)

	dir := t.TempDir()
	fontPath := filepath.Join(dir, "regular.ttf")
	require.NoError(t, os.WriteFile(fontPath, packed.MesloLGSNFRegularTTF, 0o600))

	candidates, err := discoverFonts(family, []string{dir})
	require.NoError(t, err)
	require.Equal(t, []fontCandidate{{
		Family: family,
		Style:  StyleRegular,
		Path:   fontPath,
		Index:  0,
	}}, candidates)
}

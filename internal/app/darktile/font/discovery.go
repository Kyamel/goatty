package font

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/liamg/fontinfo"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type fontCandidate struct {
	Family string
	Style  StyleName
	Path   string
	Index  int
}

func findFontCandidates(family string) ([]fontCandidate, error) {
	matches, err := fontinfo.Match(fontinfo.MatchFamily(family))
	if err != nil {
		return nil, err
	}

	candidates := make([]fontCandidate, 0, len(matches))
	seen := make(map[string]struct{})
	for _, match := range matches {
		style := normalizeStyleName(match.Style)
		candidate := fontCandidate{Family: match.Family, Style: style, Path: match.Path}
		candidates = appendUniqueFont(candidates, seen, candidate)
	}

	discovered, err := discoverFonts(family, platformFontDirs())
	if err != nil {
		return nil, err
	}
	for _, candidate := range discovered {
		candidates = appendUniqueFont(candidates, seen, candidate)
	}

	return candidates, nil
}

// ListFamilies returns the installed regular font families that Darktile can
// load. On macOS this includes fonts in the user, local, network and system
// Font Book locations.
func ListFamilies() ([]string, error) {
	matches, err := fontinfo.List()
	if err != nil {
		return nil, err
	}

	families := make(map[string]string)
	for _, match := range matches {
		if normalizeStyleName(match.Style) == StyleRegular {
			families[strings.ToLower(match.Family)] = match.Family
		}
	}

	discovered, err := discoverFonts("", platformFontDirs())
	if err != nil {
		return nil, err
	}
	for _, candidate := range discovered {
		if candidate.Style == StyleRegular {
			families[strings.ToLower(candidate.Family)] = candidate.Family
		}
	}

	result := make([]string, 0, len(families))
	for _, family := range families {
		result = append(result, family)
	}
	sort.Strings(result)
	return result, nil
}

func appendUniqueFont(candidates []fontCandidate, seen map[string]struct{}, candidate fontCandidate) []fontCandidate {
	key := candidate.Path + "\x00" + strconv.Itoa(candidate.Index)
	if _, ok := seen[key]; ok {
		return candidates
	}
	seen[key] = struct{}{}
	return append(candidates, candidate)
}

func normalizeStyleName(style string) StyleName {
	normalized := strings.ToLower(strings.TrimSpace(style))
	bold := strings.Contains(normalized, "bold")
	italic := strings.Contains(normalized, "italic") || strings.Contains(normalized, "oblique")

	switch {
	case bold && italic:
		return StyleBoldItalic
	case bold:
		return StyleBold
	case italic:
		return StyleItalic
	default:
		return StyleRegular
	}
}

func platformFontDirs() []string {
	if runtime.GOOS != "darwin" {
		return nil
	}

	home, _ := os.UserHomeDir()
	dirs := []string{
		"/Library/Fonts",
		"/Network/Library/Fonts",
		"/System/Library/Fonts",
		"/System/Library/Fonts/Supplemental",
	}
	if home != "" {
		dirs = append([]string{filepath.Join(home, "Library", "Fonts")}, dirs...)
	}
	return dirs
}

func discoverFonts(family string, dirs []string) ([]fontCandidate, error) {
	var candidates []fontCandidate
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}

		err = filepath.WalkDir(dir, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				if entry != nil && entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if entry.IsDir() || !isFontFile(path) {
				return nil
			}

			matches, err := candidatesFromFile(path, family)
			if err != nil {
				return nil
			}
			candidates = append(candidates, matches...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return candidates, nil
}

func isFontFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".dfont", ".otc", ".otf", ".ttc", ".ttf":
		return true
	default:
		return false
	}
}

func candidatesFromFile(path, wantedFamily string) ([]fontCandidate, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	collection, err := opentype.ParseCollectionReaderAt(f)
	if err != nil {
		return nil, err
	}

	var candidates []fontCandidate
	for index := 0; index < collection.NumFonts(); index++ {
		font, err := collection.Font(index)
		if err != nil {
			continue
		}
		family, err := preferredFontName(font, sfnt.NameIDTypographicFamily, sfnt.NameIDFamily)
		if err != nil || (wantedFamily != "" && !strings.EqualFold(family, wantedFamily)) {
			continue
		}
		style, _ := preferredFontName(font, sfnt.NameIDTypographicSubfamily, sfnt.NameIDSubfamily)
		candidates = append(candidates, fontCandidate{
			Family: family,
			Style:  normalizeStyleName(style),
			Path:   path,
			Index:  index,
		})
	}
	return candidates, nil
}

func preferredFontName(font *sfnt.Font, ids ...sfnt.NameID) (string, error) {
	var lastErr error
	for _, id := range ids {
		name, err := font.Name(nil, id)
		if err == nil && name != "" {
			return name, nil
		}
		lastErr = err
	}
	return "", lastErr
}

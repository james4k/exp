package ui

type FontStyle struct {
	Family string
	Size   FontSize
}

const (
	Pt FontPt = 1
	Px FontPx = 1

	DefaultFontSize = 16 * Pt
)

type FontSize interface {
	FontSize(dpi int) int
}

type FontPt int

func (f FontPt) FontSize(dpi int) int {
	return int(f) * dpi
}

type FontPx int

func (f FontPx) FontSize(dpi int) int {
	return int(f)
}

// TODO: FontEm/FontRootEm

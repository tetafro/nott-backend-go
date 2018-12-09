package markdown

import (
	"github.com/Depado/bfchroma"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

// theme is a color theme for rendering markdown to HTML.
const theme = "monokailight"

// Render renders markdown to HTML.
func Render(markdown string) (html string) {
	r := bfchroma.NewRenderer(bfchroma.Style(theme))
	b := blackfriday.Run(
		[]byte(markdown),
		blackfriday.WithRenderer(r),
	)
	return string(b)
}

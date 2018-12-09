package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	t.Run("plain text", func(t *testing.T) {
		md := "Hello, world"
		html := "<p>Hello, world</p>\n"
		assert.Equal(t, html, Render(md))
	})

	t.Run("text with special markdown characters", func(t *testing.T) {
		md := "# Hello, world\n" +
			"* one\n" +
			"* two"
		html := "<h1>Hello, world</h1>\n\n" +
			"<ul>\n" +
			"<li>one</li>\n" +
			"<li>two</li>\n" +
			"</ul>\n"
		assert.Equal(t, html, Render(md))
	})

	t.Run("code", func(t *testing.T) {
		md := "```go\n" +
			"import \"fmt\"\n" +
			"func main() {\n" +
			"\t fmt.Println(\"Hi!\")\n" +
			"}\n" +
			"```"
		// Some text inside styled "<pre>" tags
		re := `<pre style=".+">(.|\s)+<\/pre>`
		assert.Regexp(t, re, Render(md))
	})
}

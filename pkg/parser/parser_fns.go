package parser

import (
	"bytes"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.New(
				meta.WithStoresInDocument(),
			),
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			mathjax.MathJax,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	p := md.Parser()
	r := md.Renderer()

	return &Parser{
		md:       md,
		parser:   p,
		renderer: r,
	}
}

// parse file and sets the parsed node, return the metadata
func (p *Parser) Parse(file *File) *Frontmatter {
	// ctx := parser.NewContext()
	node := p.parser.Parse(text.NewReader(file.Body))

	file.Node = node

	meta := node.OwnerDocument().Meta()

	frontmatter := NewFrontmatter()

	for key, value := range meta {
		frontmatter.Set(key, value)
	}

	return frontmatter
}

// render file and return the rendered bytes
func (p *Parser) Render(file *File) (bytes.Buffer, error) {
	var b bytes.Buffer

	err := p.renderer.Render(&b, file.Body, file.Node)

	return b, err
}

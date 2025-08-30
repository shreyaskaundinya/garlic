package parser

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
)

type Parser struct {
	md       goldmark.Markdown
	parser   parser.Parser
	renderer renderer.Renderer
}

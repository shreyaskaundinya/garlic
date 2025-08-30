package parser

import "github.com/yuin/goldmark/ast"

const (
	FILE_TYPE_TEMPLATE  = "FILE_TYPE_TEMPLATE"
	FILE_TYPE_COMPONENT = "FILE_TYPE_COMPONENT"
	FILE_TYPE_MARKDOWN  = "FILE_TYPE_MARKDOWN"
)

type File struct {
	// Type
	Type string

	// Path
	Path string

	// Body
	Body []byte

	// Ast Node
	Node ast.Node
}

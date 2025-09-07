package server

import (
	"github.com/fsnotify/fsnotify"
	"github.com/shreyaskaundinya/garlic/pkg/parser"
)

/*
FOLDER -> FILES -> PARSE (each)
-> setup meta
-> render html
-> build dest
-> serve if required
*/
type Server struct {
	// TODO : http server

	// source folder path
	SrcPath string

	// destination folder path
	DestPath string

	// metadata
	MD *parser.Metadata

	// template metadata
	TemplateMD *parser.Metadata

	// components metadata
	ComponentsMD *parser.Metadata

	// Parser
	Parser *parser.Parser

	// file chan
	fileCh chan *parser.File

	// parse chan
	parseCh chan *parser.File

	// render chan
	renderCh chan *parser.File
}

type RenderEvent struct {
	Event               fsnotify.Event
	RenderAll           bool
	ProcessAssets       bool
	ProcessDependencies bool
	ProcessContent      bool
	ProcessTags         bool
}

package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/aarol/reload"
	"github.com/fsnotify/fsnotify"
	"github.com/shreyaskaundinya/garlic/pkg/parser"
	"github.com/shreyaskaundinya/garlic/pkg/utils"
)

func NewServer(SrcPath string, DestPath string) (*Server, error) {
	// check if source folder exists
	isSrcExist, srcErr := utils.PathExists(SrcPath)

	if !isSrcExist {
		return nil, srcErr
	}

	// check if blogs folder exists
	blogsPath := filepath.Join(SrcPath, "content")
	isBlogsExist, blogsErr := utils.PathExists(blogsPath)

	if !isBlogsExist {
		return nil, blogsErr
	}

	// check if templates folder exists
	templatePath := filepath.Join(SrcPath, "templates")
	isTemplatesExist, templatesErr := utils.PathExists(templatePath)

	if !isTemplatesExist {
		return nil, templatesErr
	}

	// check if destination folder exists
	isDestExist, destErr := utils.PathExists(DestPath)

	if !isDestExist {
		return nil, destErr
	}

	return &Server{
		SrcPath:      filepath.FromSlash(SrcPath),
		DestPath:     filepath.FromSlash(DestPath),
		MD:           parser.NewMetadataMap(),
		TemplateMD:   parser.NewMetadataMap(),
		ComponentsMD: parser.NewMetadataMap(),
		Parser:       parser.NewParser(),
		fileCh:       make(chan *parser.File, 4),
		parseCh:      make(chan *parser.File, 4),
		renderCh:     make(chan *parser.File, 4),
	}, nil
}

func (s *Server) Start() {
	log := utils.NewLogger()

	// need to add a watcher to check for changes in the source folder
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalw("Error creating watcher", "error", err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				log.Infow("event:", "event", event)

				log.Infow("modified file:", "event", event.Name)

				err = s.render()
				if err != nil {
					log.Errorw("Error rendering", "error", err)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Infow("error:", "error", err)
			}
		}
	}()

	// walk all the paths inside source and add all the directories to the watcher
	err = filepath.WalkDir(s.SrcPath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				log.Errorw("Error adding watcher", "error", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Errorw("Error walking source", "error", err)
		return
	}

	// read templates and render once on init
	err = s.render()
	if err != nil {
		log.Errorw("Error rendering", "error", err)
		return
	}

	// serve the html files
	go s.serve()

	// Block main goroutine forever.
	<-make(chan struct{})
}

func (s *Server) serve() {
	log := utils.NewLogger()

	fs := http.FileServer(http.Dir(s.DestPath))

	var handler http.Handler = fs

	// Call `New()` with a list of directories to recursively watch
	reloader := reload.New(s.DestPath)

	// Optionally, define a callback to
	// invalidate any caches
	reloader.OnReload = func() {
		log.Infow("Reloading http server...")
	}

	// Use the Handle() method as a middleware
	handler = reloader.Handle(handler)

	http.ListenAndServe(":8084", handler)
}

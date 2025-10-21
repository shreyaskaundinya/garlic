package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/aarol/reload"
	"github.com/fsnotify/fsnotify"
	"github.com/shreyaskaundinya/garlic/models"
	"github.com/shreyaskaundinya/garlic/pkg/parser"
	"github.com/shreyaskaundinya/garlic/pkg/utils"
)

func checkAndCreateFolder(path, aliasName string) error {
	log := utils.NewLogger()

	isExist, err := utils.PathExists(path)

	if err != nil {
		return err
	}

	if !isExist {
		log.Infow(aliasName + " does not exist, creating it...")

		// create the folder
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkAndSeedFile(srcPath, contentPath, aliasName string) error {
	log := utils.NewLogger()

	srcFileExist, srcFileErr := utils.PathExists(srcPath)

	if srcFileErr != nil {
		return srcFileErr
	}

	if !srcFileExist {
		log.Infow(aliasName + " does not exist, creating it...")

		// read the content file
		content, err := os.ReadFile(contentPath)
		if err != nil {
			return err
		}

		// create the index MD file
		err = os.WriteFile(srcPath, content, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedSrc(config *models.Config) error {
	SrcPath := config.SrcPath
	return checkAndCreateFolder(SrcPath, "src")
}

func seedContent(config *models.Config) error {
	SrcPath := config.SrcPath
	blogsPath := filepath.Join(SrcPath, "content")

	err := checkAndCreateFolder(blogsPath, "src/content")
	if err != nil {
		return err
	}

	// create the default content files
	if config.ShouldSeedFiles {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		indexMDPath := filepath.Join(blogsPath, "index.md")

		defaultContentPath := filepath.Join(cwd, "default-assets", "content", "index.md")

		err = checkAndSeedFile(indexMDPath, defaultContentPath, "Index MD")
		if err != nil {
			return err
		}
	}

	return nil
}

func seedTemplates(config *models.Config) error {
	SrcPath := config.SrcPath

	templatesPath := filepath.Join(SrcPath, "templates")

	err := checkAndCreateFolder(templatesPath, "src/templates")
	if err != nil {
		return err
	}

	// seed the default templates
	if config.ShouldSeedFiles {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// templates/index.html
		indexHTMLPath := filepath.Join(templatesPath, "index.html")
		defaultTemplatePath := filepath.Join(cwd, "default-assets", "templates", "index.html")
		err = checkAndSeedFile(indexHTMLPath, defaultTemplatePath, "Index HTML")
		if err != nil {
			return err
		}

		// templates/_tags.html
		tagsHTMLPath := filepath.Join(templatesPath, "_tags.html")
		defaultTemplatePath = filepath.Join(cwd, "default-assets", "templates", "_tags.html")
		err = checkAndSeedFile(tagsHTMLPath, defaultTemplatePath, "Tags HTML")
		if err != nil {
			return err
		}

		// templates/_individual_tag.html
		individualTagHTMLPath := filepath.Join(templatesPath, "_individual_tag.html")
		defaultTemplatePath = filepath.Join(cwd, "default-assets", "templates", "_individual_tag.html")

		err = checkAndSeedFile(individualTagHTMLPath, defaultTemplatePath, "Individual Tag HTML")
		if err != nil {
			return err
		}
	}

	return nil
}

func seedComponents(config *models.Config) error {
	SrcPath := config.SrcPath

	componentsPath := filepath.Join(SrcPath, "components")

	err := checkAndCreateFolder(componentsPath, "src/components")
	if err != nil {
		return err
	}

	// seed the default components
	if config.ShouldSeedFiles {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// components/Footerbar.html
		footerbarHTMLPath := filepath.Join(componentsPath, "Footerbar.html")
		defaultComponentPath := filepath.Join(cwd, "default-assets", "components", "Footerbar.html")
		err = checkAndSeedFile(footerbarHTMLPath, defaultComponentPath, "Footerbar HTML")
		if err != nil {
			return err
		}

		// components/Navbar.html
		navbarHTMLPath := filepath.Join(componentsPath, "Navbar.html")
		defaultComponentPath = filepath.Join(cwd, "default-assets", "components", "Navbar.html")
		err = checkAndSeedFile(navbarHTMLPath, defaultComponentPath, "Navbar HTML")
		if err != nil {
			return err
		}

		// components/Tags.html
		tagsHTMLPath := filepath.Join(componentsPath, "Tags.html")
		defaultComponentPath = filepath.Join(cwd, "default-assets", "components", "Tags.html")
		err = checkAndSeedFile(tagsHTMLPath, defaultComponentPath, "Tags HTML")
		if err != nil {
			return err
		}
	}

	return nil
}

func seedAssets(config *models.Config) error {
	SrcPath := config.SrcPath

	assetsPath := filepath.Join(SrcPath, "assets")

	err := checkAndCreateFolder(assetsPath, "src/assets")
	if err != nil {
		return err
	}

	// seed the default assets
	if config.ShouldSeedFiles {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// styles folder
		stylesFolderPath := filepath.Join(assetsPath, "styles")

		err = checkAndCreateFolder(stylesFolderPath, "src/assets/styles")
		if err != nil {
			return err
		}

		// assets/styles/global.css
		globalCSSPath := filepath.Join(assetsPath, "styles", "global.css")
		defaultAssetPath := filepath.Join(cwd, "default-assets", "assets", "styles", "global.css")
		err = checkAndSeedFile(globalCSSPath, defaultAssetPath, "Global CSS")
		if err != nil {
			return err
		}
	}

	return nil
}

func seedDest(config *models.Config) error {
	return checkAndCreateFolder(config.DestPath, "Destination")
}

func seed(config *models.Config) error {
	err := seedSrc(config)
	if err != nil {
		return err
	}

	err = seedContent(config)
	if err != nil {
		return err
	}

	err = seedTemplates(config)
	if err != nil {
		return err
	}

	err = seedComponents(config)
	if err != nil {
		return err
	}

	err = seedAssets(config)
	if err != nil {
		return err
	}

	err = seedDest(config)
	if err != nil {
		return err
	}

	return nil
}

func NewServer(config *models.Config) (*Server, error) {
	log := utils.NewLogger()

	err := seed(config)
	if err != nil {
		log.Errorw("Error seeding", "error", err)
		return nil, err
	}

	return &Server{
		SrcPath:      filepath.FromSlash(config.SrcPath),
		DestPath:     filepath.FromSlash(config.DestPath),
		MD:           parser.NewMetadataMap(),
		TemplateMD:   parser.NewMetadataMap(),
		ComponentsMD: parser.NewMetadataMap(),
		Parser:       parser.NewParser(),
		fileCh:       make(chan *parser.File, 4),
		parseCh:      make(chan *parser.File, 4),
		renderCh:     make(chan *parser.File, 4),
	}, nil
}

func (s *Server) Start(config *models.Config) {
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

				log.Infow("modified file:", "event", event.Name)

				err = s.render(&RenderEvent{Event: event, RenderAll: false})
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
	err = s.render(&RenderEvent{RenderAll: true})
	if err != nil {
		log.Errorw("Error rendering", "error", err)
		return
	}

	// serve the html files
	if config.ShouldServe {
		go s.serve()
	}

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

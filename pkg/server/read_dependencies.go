package server

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shreyaskaundinya/garlic/pkg/parser"
	"github.com/shreyaskaundinya/garlic/pkg/utils"
)

func (s *Server) readDependencies() error {
	err := s.readTemplates()
	if err != nil {
		return err
	}

	err = s.readComponents()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) readComponents() error {
	log := utils.NewLogger()

	start := time.Now()

	defer func() {
		log.Debugw("Time taken to read components", "time", time.Since(start))
	}()

	// collect components from the directory
	err := filepath.WalkDir(filepath.Join(s.SrcPath, "components"), func(
		path string,
		info os.DirEntry,
		err error,
	) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// strip the .html extension
		componentName := strings.ToLower(
			strings.TrimSuffix(info.Name(), ".html"),
		)

		f := parser.NewFile(path, parser.FILE_TYPE_COMPONENT)
		err = f.ReadFile()
		if err != nil {
			return err
		}

		s.ComponentsMD.Set(componentName, &parser.Meta{
			Title:       componentName,
			Description: "",
			F:           f,
			Tags:        make([]string, 0),
		})

		log.Infow("Component File: ", "componentName", componentName)

		return nil
	})

	if err != nil {
		log.Errorw("Error reading components", "error", err)
		return err
	}

	return nil
}

func (s *Server) readTemplates() error {
	log := utils.NewLogger()

	start := time.Now()

	defer func() {
		log.Debugw("Time taken to read templates", "time", time.Since(start))
	}()

	// read templates
	err := filepath.Walk(filepath.Join(s.SrcPath, "templates"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		log.Infow("Template File: ", "path", path)

		f := parser.NewFile(path, parser.FILE_TYPE_TEMPLATE)

		s.TemplateMD.Set(f.Path, &parser.Meta{
			Title:       path,
			Description: "",
			F:           f,
			Tags:        make([]string, 0),
		})

		err = f.ReadFile()

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorw("Error reading templates", "error", err)
		return err
	}

	return nil
}

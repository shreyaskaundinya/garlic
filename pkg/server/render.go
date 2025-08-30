package server

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shreyaskaundinya/garlic/pkg/parser"
	"github.com/shreyaskaundinya/garlic/pkg/utils"
)

func (s *Server) runBeforeRenderProcess() error {
	log := utils.NewLogger()

	err := s.readDependencies()
	if err != nil {
		log.Errorw("Error reading dependencies", "error", err)
		return err
	}

	return nil
}

func (s *Server) runAfterRenderProcess() error {
	log := utils.NewLogger()

	err := s.processTags()
	if err != nil {
		log.Errorw("Error processing tags", "error", err)
		return err
	}

	return nil
}

func (s *Server) setupMarkdown(path string) (*parser.Meta, error) {
	log := utils.NewLogger()

	f := parser.NewFile(path, parser.FILE_TYPE_MARKDOWN)

	// TODO : could optimize by reading only till end of frontmatter
	err := f.ReadFile()

	if err != nil {
		return nil, err
	}

	frontmatter := s.Parser.Parse(f)
	if frontmatter == nil {
		log.Errorw("Error parsing frontmatter", "path", path)
		return nil, errors.New("error parsing frontmatter")
	}

	title, _ := frontmatter.Get("title")

	description, _ := frontmatter.Get("description")

	sitepath := strings.TrimPrefix(
		strings.TrimSuffix(path, filepath.Ext(path)), s.SrcPath,
	)

	// if sitepath ends with index, remove it
	sitepath = strings.TrimSuffix(sitepath, "index")

	markdownMeta := &parser.Meta{
		Title:       utils.GetSafeValue[string](title),
		Sitepath:    sitepath,
		Description: utils.GetSafeValue[string](description),
		F:           f,
		Tags:        frontmatter.GetTags(),
		Frontmatter: frontmatter,
	}

	s.MD.Set(f.Path, markdownMeta)

	return markdownMeta, nil
}

func (s *Server) render() error {
	log := utils.NewLogger()

	err := s.runBeforeRenderProcess()
	if err != nil {
		return err
	}

	start := time.Now()

	defer func() {
		log.Debugw("Time taken to render", "time", time.Since(start))
	}()

	// currentDir := ""
	// depth := 1
	// read blogs
	err = filepath.WalkDir(filepath.Join(s.SrcPath, "content"), func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath := strings.Split(path, s.SrcPath)[1]

		depth := strings.Count(relativePath, "\\")
		log.Infow(relativePath, "depth", depth)

		if info.IsDir() {
			log.Infow("--- DIR : ", "isDir", info.IsDir())
			return nil
		}

		log.Infow("Blog File: ", "path", path)

		markdownMeta, err := s.setupMarkdown(path)
		if err != nil {
			return err
		}

		publishIf, ok := markdownMeta.Frontmatter.Get("publish")

		if !ok {
			log.Infow("%s missing publish attribute in metadata (front matter)", "name", info.Name())
			return nil
		}

		publish := publishIf.(bool)

		if !publish {
			return nil
		}

		html, err := s.Parser.Render(markdownMeta.F)

		if err != nil {
			return err
		}

		// inject html into template
		content, err := s.injectHTML(markdownMeta, html)
		if err != nil {
			return err
		}

		fileName := utils.FileNameWithoutExtension(
			strings.TrimPrefix(
				relativePath, filepath.Dir(relativePath),
			),
		)

		log.Infow("File Name: ", "fileName", fileName)

		// make dirs if not already made
		// filepath.Dir(relativePath[1]) => content/projects/
		var renderFolderPath string
		if strings.HasSuffix(fileName, "index") {
			renderFolderPath = filepath.Join(
				s.DestPath,
				filepath.Dir(relativePath),
			)
		} else {
			renderFolderPath = filepath.Join(
				s.DestPath,
				filepath.Dir(relativePath),
				fileName,
			)
		}

		doesDestPathExist, err := utils.PathExists(renderFolderPath)

		if !doesDestPathExist || err != nil {
			err = os.MkdirAll(renderFolderPath, os.ModeDir)

			if err != nil {
				panic(err)
			}
		}

		err = markdownMeta.F.WriteToDest(
			renderFolderPath,
			"index.html",
			[]byte(content),
		)

		if err != nil {
			log.Errorw("Error writing to dest", "error", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorw("Error rendering", "error", err)
		return err
	}

	// run after render process
	err = s.runAfterRenderProcess()
	if err != nil {
		return err
	}

	return nil
}

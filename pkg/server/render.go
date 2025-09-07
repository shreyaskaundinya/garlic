package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shreyaskaundinya/garlic/pkg/parser"
	"github.com/shreyaskaundinya/garlic/pkg/utils"
)

func (s *Server) runBeforeRenderProcess(event *RenderEvent) error {
	log := utils.NewLogger()

	if event.ProcessAssets {
		err := s.readAndCopyAssets()
		if err != nil {
			log.Errorw("Error reading and copying assets", "error", err)
			return err
		}
	}

	if event.ProcessDependencies {
		err := s.readDependencies()
		if err != nil {
			log.Errorw("Error reading dependencies", "error", err)
			return err
		}
	}

	return nil
}

func (s *Server) readAndCopyAssets() error {
	log := utils.NewLogger()

	assetsSrcPath := filepath.Join(s.SrcPath, "assets")
	assetsDestPath := filepath.Join(s.DestPath, "assets")

	// Check if assets directory exists
	if _, err := os.Stat(assetsSrcPath); os.IsNotExist(err) {
		log.Infow("Assets directory does not exist, skipping asset copy", "path", assetsSrcPath)
		return nil
	}

	// copy assets from source to destination
	err := filepath.WalkDir(assetsSrcPath, func(srcPath string, info os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}

		// Get relative path from assets source directory
		relPath, err := filepath.Rel(assetsSrcPath, srcPath)
		if err != nil {
			return fmt.Errorf("error getting relative path: %w", err)
		}

		destPath := filepath.Join(assetsDestPath, relPath)

		if info.IsDir() {
			// Create directory in destination
			err = os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directory: %w", err)
			}
			return nil
		}

		// Read source file content
		content, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		// Create destination directory if it doesn't exist
		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}

		// Remove existing file if it exists to avoid permission issues
		if _, err := os.Stat(destPath); err == nil {
			err = os.Remove(destPath)
			if err != nil {
				return fmt.Errorf("error removing existing file: %w", err)
			}
		}

		// Write file to destination
		err = os.WriteFile(destPath, content, 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}

		log.Debugw("Copied asset file", "src", srcPath, "dest", destPath)
		return nil
	})
	if err != nil {
		log.Errorw("Error copying assets", "error", err)
		return err
	}

	log.Infow("Successfully copied all assets", "src", assetsSrcPath, "dest", assetsDestPath)
	return nil
}

func (s *Server) runAfterRenderProcess(event *RenderEvent) error {
	log := utils.NewLogger()

	if event.ProcessTags {
		err := s.processTags()
		if err != nil {
			log.Errorw("Error processing tags", "error", err)
			return err
		}
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

	contentPath := strings.TrimPrefix(path, s.SrcPath)
	pathSplits := strings.Split(contentPath, string(os.PathSeparator))

	sitePath := filepath.Join(pathSplits[2:]...)

	sitepath := strings.TrimPrefix(
		strings.TrimSuffix(sitePath, filepath.Ext(sitePath)),
		s.SrcPath,
	)

	// if sitepath ends with index, remove it
	sitepath = strings.TrimSuffix(sitepath, "index")

	sitepath = "/" + sitepath

	log.Infow("[Sitepath]",
		"contentPath", contentPath,
		"pathSplits", pathSplits,
		"sitePath", sitePath,
		"sitepath", sitepath,
	)

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

func (s *Server) processEvent(event *RenderEvent) {
	if event.RenderAll {
		event.ProcessAssets = true
		event.ProcessDependencies = true
		event.ProcessTags = true
		event.ProcessContent = true
		return
	}

	log := utils.NewLogger()

	// read the path of changed file
	path := event.Event.Name
	relativePath, err := filepath.Rel(s.SrcPath, path)
	if err != nil {
		log.Errorw("Error getting relative path", "error", err)
		event.RenderAll = true
		s.processEvent(event)
		return
	}

	// if assets are changed, process assets
	if strings.HasPrefix(relativePath, "assets") {
		event.ProcessAssets = true
	}

	// if dependencies are changed, process dependencies and tags
	if strings.HasPrefix(relativePath, "components") {
		event.ProcessDependencies = true
		event.ProcessContent = true
		event.ProcessTags = true
	}

	if strings.HasPrefix(relativePath, "templates") {
		event.ProcessDependencies = true
		event.ProcessContent = true
		event.ProcessTags = true
	}

	if strings.HasPrefix(relativePath, "content") {
		event.ProcessDependencies = true
		event.ProcessContent = true
		event.ProcessTags = true
	}
}

func (s *Server) render(event *RenderEvent) error {
	start := time.Now()

	log := utils.NewLogger()

	s.processEvent(event)

	err := s.runBeforeRenderProcess(event)
	if err != nil {
		return err
	}

	defer func() {
		log.Debugw("Time taken to render", "time", time.Since(start))
	}()

	if event.ProcessContent {
		// currentDir := ""
		// depth := 1
		// read blogs
		err = filepath.WalkDir(filepath.Join(s.SrcPath, "content"), func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			relativePath := strings.Split(path, s.SrcPath)[1]

			// depth := strings.Count(relativePath, "\\")
			// log.Infow(relativePath, "depth", depth)

			if info.IsDir() {
				// log.Infow("--- DIR : ", "isDir", info.IsDir())
				return nil
			}

			// log.Infow("Blog File: ", "path", path)

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

			log.Infow("[debug] relativePath", "relativePath", relativePath)

			// remove content/ from the relative path
			splits := strings.Split(relativePath, string(os.PathSeparator))

			if len(splits) == 0 {
				return nil
			}

			relativePath = filepath.Join(splits[2:]...)

			log.Infow("[debug] relativePath", "relativePath", relativePath, "splits", splits)

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
	}

	// run after render process
	err = s.runAfterRenderProcess(event)
	if err != nil {
		return err
	}

	return nil
}

package server

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/shreyaskaundinya/garlic/pkg/parser"
	"github.com/shreyaskaundinya/garlic/pkg/utils"
)

func (s *Server) processTags() error {
	log := utils.NewLogger()

	log.Debugw("Processing tags")

	start := time.Now()

	defer func() {
		log.Debugw("Time taken to process tags", "time", time.Since(start))
	}()

	// need to generate a tags page with all tags
	tagToFilesMap := map[string][]*parser.Meta{}
	tagsMapset := mapset.NewSet[string]()

	s.MD.Store.Range(func(_ string, value *parser.Meta) bool {
		publishIf, ok := value.Frontmatter.Get("publish")
		if !ok {
			return true
		}
		publish := publishIf.(bool)
		if !publish {
			return true
		}

		if !publish {
			return true
		}

		tags := value.Tags
		for _, tag := range tags {
			if _, ok := tagToFilesMap[tag]; !ok {
				tagToFilesMap[tag] = []*parser.Meta{}
				tagsMapset.Add(tag)
			}

			tagToFilesMap[tag] = append(tagToFilesMap[tag], value)
			tagsMapset.Add(tag)
		}
		return true
	})

	tags := tagsMapset.ToSlice()

	log.Debugw("Tags", "tags", tags)

	sort.Strings(tags)

	template := `
			{{ $content }}
		`

	// find the template for the tag page
	tagsTemplatePath := filepath.Join(
		s.SrcPath,
		"templates",
		"_tags.html",
	)

	log.Infow("Tags Template Path: ", "tagsTemplatePath", tagsTemplatePath)

	t, ok := s.TemplateMD.Get(
		tagsTemplatePath,
	)
	if ok {
		log.Infow("Found tags template")

		err := t.F.ReadFile()
		if err != nil {
			return err
		}

		template = string(t.F.Body)

		log.Infow("Template: ", "template", template)
	}

	tagsHTML := `<ul>
	`

	for _, tag := range tags {
		count := len(tagToFilesMap[tag])
		tagsHTML += fmt.Sprintf(`
		<li><a href="/tags/%s">%s ( %d )</a></li>
		`, tag, tag, count)
	}

	tagsHTML += `</ul>`

	destPath := path.Join(s.DestPath, "tags")
	doesDestPathExist, err := utils.PathExists(destPath)
	if err != nil {
		return err
	}

	if !doesDestPathExist {
		err = os.MkdirAll(destPath, os.ModeDir)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(path.Join(destPath, "index.html"), []byte(tagsHTML), 0644)
	if err != nil {
		return err
	}

	// create a page for each tag
	for _, tag := range tags {
		tagHTML := `<ul>
		`
		for _, meta := range tagToFilesMap[tag] {
			if meta == nil {
				continue
			}

			sitepath := meta.Sitepath
			title := meta.Title

			tagHTML += fmt.Sprintf(`
			<li><a href="%s">%s</a></li>
			`, sitepath, title)
		}

		tagHTML += `</ul>`

		pageHTML := strings.ReplaceAll(template, "{{ $content }}", tagHTML)
		pageHTML = strings.ReplaceAll(pageHTML, "{{ $title }}", tag)

		content, err := s.injectComponents(&parser.Meta{}, &pageHTML)
		if err != nil {
			return err
		}

		destPath := path.Join(s.DestPath, "tags", tag)

		doesDestPathExist, err := utils.PathExists(destPath)
		if err != nil {
			return err
		}

		if !doesDestPathExist {
			err = os.MkdirAll(destPath, os.ModeDir)
			if err != nil {
				return err
			}
		}

		err = os.WriteFile(path.Join(destPath, "index.html"), []byte(content), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

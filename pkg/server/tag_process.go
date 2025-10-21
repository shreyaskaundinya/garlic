package server

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/shreyaskaundinya/garlic/pkg/parser"
	"github.com/shreyaskaundinya/garlic/pkg/utils"
	"golang.org/x/net/html"
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

	tagsTemplate, ok := s.TemplateMD.Get(
		tagsTemplatePath,
	)
	if ok {
		log.Infow("Found tags template")

		err := tagsTemplate.F.ReadFile()
		if err != nil {
			return err
		}

		template = string(tagsTemplate.F.Body)
	}

	tagsTemplateAST, err := html.Parse(bytes.NewReader([]byte(template)))
	if err != nil {
		return err
	}

	tagsList := &html.Node{
		Type: html.ElementNode,
		Data: "ul",
	}

	for _, tag := range tags {
		count := len(tagToFilesMap[tag])

		li := &html.Node{
			Type: html.ElementNode,
			Data: "li",
		}

		a := &html.Node{
			Type: html.ElementNode,
			Data: "a",
		}

		a.Attr = append(a.Attr, html.Attribute{
			Key: "href",
			Val: fmt.Sprintf("/tags/%s", tag),
		})

		a.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: fmt.Sprintf("%s ( %d )", tag, count),
		})

		li.AppendChild(a)

		tagsList.AppendChild(li)
	}

	// replace recursively
	s.recursivelyReplace(
		tagsTemplateAST,
		tagsList,
		&parser.Meta{},
		nil,
	)

	tagsHTML := bytes.NewBuffer(make([]byte, 0))
	err = html.Render(tagsHTML, tagsTemplateAST)
	if err != nil {
		return err
	}

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

	err = os.WriteFile(path.Join(destPath, "index.html"), tagsHTML.Bytes(), 0644)
	if err != nil {
		return err
	}

	// ------------------------------------------------------------------

	individualTagTemplatePath := filepath.Join(
		s.SrcPath,
		"templates",
		"_individual_tag.html",
	)

	log.Infow("Tags Template Path: ", "tagsTemplatePath", tagsTemplatePath)

	individualTagTemplate, ok := s.TemplateMD.Get(
		individualTagTemplatePath,
	)
	if ok {
		log.Infow("Found individual tag template")

		err := individualTagTemplate.F.ReadFile()
		if err != nil {
			return err
		}

		template = string(individualTagTemplate.F.Body)
	}

	// create a page for each tag
	for _, tag := range tags {
		tagsTemplateAST, err = html.Parse(bytes.NewReader([]byte(template)))
		if err != nil {
			return err
		}

		tagList := &html.Node{
			Type: html.ElementNode,
			Data: "ul",
		}

		for _, meta := range tagToFilesMap[tag] {
			if meta == nil {
				continue
			}

			li := &html.Node{
				Type: html.ElementNode,
				Data: "li",
			}

			sitepath := meta.Sitepath
			title := meta.Title

			a := &html.Node{
				Type: html.ElementNode,
				Data: "a",
			}

			a.Attr = append(a.Attr, html.Attribute{
				Key: "href",
				Val: sitepath,
			})

			a.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: title,
			})

			li.AppendChild(a)

			tagList.AppendChild(li)
		}

		// replace recursively
		s.recursivelyReplace(
			tagsTemplateAST,
			tagList,
			&parser.Meta{
				Title:    tag,
				Sitepath: fmt.Sprintf("/tags/%s", tag),
			},
			nil,
		)

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

		tagHTML := bytes.NewBuffer(make([]byte, 0))
		err = html.Render(tagHTML, tagsTemplateAST)
		if err != nil {
			return err
		}

		err = os.WriteFile(path.Join(destPath, "index.html"), tagHTML.Bytes(), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

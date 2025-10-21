package server

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"

	"github.com/shreyaskaundinya/garlic/pkg/parser"
	"github.com/shreyaskaundinya/garlic/pkg/utils"
)

func (s *Server) recursivelyReplaceTags(
	tagTemplate *html.Node,
	fileMetadata *parser.Meta,
) {
	log := utils.NewLogger()

	log.Infow("Recursively replacing tags")

	// generate li tags for each tag
	tags := fileMetadata.Tags

	log.Infow("Tags", "tags", tags)

	liTags := make([]*html.Node, 0)

	for _, tag := range tags {
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
			Data: tag,
		})

		li.AppendChild(a)

		liTags = append(liTags, li)
	}

	for child := range tagTemplate.Descendants() {
		// log.Infow("Child", "child", child.Data)
		switch child.Type {
		case html.TextNode:
			switch {
			case strings.Contains(child.Data, "{{ $tags }}"):
				log.Infow("Injecting tags into child node")
				for _, li := range liTags {
					// log.Infow("Appending li tag", "li", li.Data)
					child.Parent.AppendChild(li)
				}

				// log.Infow("Removing child node", "child", child.Data)
				child.Parent.RemoveChild(child)

			}
		}
	}

	log.Infow("Tags injected")
}

func (s *Server) recursivelyReplace(
	templateHTML *html.Node,
	contentHTML *html.Node,
	fileMetadata *parser.Meta,
	tagsTemplate *html.Node,
) {
	log := utils.NewLogger()

	log.Infow("Recursively replacing", "node", templateHTML.Data)

	for child := range templateHTML.Descendants() {
		// log.Infow("Child", "child", child.Data)

		switch child.Type {
		case html.TextNode:
			switch {
			case strings.Contains(child.Data, "{{ $content }}"):
				// log.Infow("Injecting content")

				// remove the {{ $content }} from the child.Data
				child.Parent.InsertBefore(contentHTML, child)
				child.Parent.RemoveChild(child)
			case strings.Contains(child.Data, "{{ $title }}"):
				// log.Infow("Injecting title")

				// remove the {{ $title }} from the child.Data
				child.Data = strings.ReplaceAll(
					child.Data, "{{ $title }}", utils.GetSafeValue[string](fileMetadata.Title),
				)
			}
		case html.ElementNode:
			// skip if the element is a content element
			if len(child.Attr) > 0 && child.Attr[0].Key == "x-type" && child.Attr[0].Val == "content" {
				continue
			}

			if child.Data == "tags" {
				child.Parent.InsertBefore(tagsTemplate, child)
				child.Parent.RemoveChild(child)
			} else {
				// find child.Data in s.ComponentsMD
				component, ok := s.ComponentsMD.Get(strings.ToLower(child.Data))
				if ok {
					// log.Infow("Found component", "component", child.Data)

					// parse the component
					parsedComponent, err := html.Parse(bytes.NewReader(component.F.Body))
					if err != nil {
						return
					}

					child.Parent.InsertBefore(parsedComponent, child)
					// child.Parent.RemoveChild(child)
				}
			}
		}
	}
}

func (s *Server) injectHTML(
	fileMetadata *parser.Meta,
	mdHTML bytes.Buffer,
) (string, error) {
	log := utils.NewLogger()

	template, ok := fileMetadata.Frontmatter.Get("template")
	if !ok {
		return "", fmt.Errorf("template not found")
	}

	t, ok := s.TemplateMD.Get(
		filepath.Join(
			s.SrcPath,
			"templates",
			fmt.Sprintf("%s.html", utils.GetSafeValue[string](template)),
		),
	)

	if !ok {
		return "", fmt.Errorf("template not found")
	}

	// st := string(t.F.Body)

	// parse the templates html
	parsedTemplate, err := html.Parse(bytes.NewReader(t.F.Body))
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	parsedMDHTML, err := html.Parse(bytes.NewReader(mdHTML.Bytes()))
	if err != nil {
		return "", fmt.Errorf("error parsing md html: %w", err)
	}

	parsedMDHTML.Attr = append(parsedMDHTML.Attr, html.Attribute{
		Key: "x-type",
		Val: "content",
	})

	// if we find <Tags />, we need to inject the tags from the file metadata
	// load the tags template
	tagsTemplateFile, ok := s.ComponentsMD.Get(
		"tags",
	)
	if !ok {
		return "", fmt.Errorf("tags template not found")
	}

	log.Infow("Tags Template File", "tagsTemplateBody", string(tagsTemplateFile.F.Body))

	// parse the tags template
	tagsTemplate, err := html.Parse(bytes.NewReader(tagsTemplateFile.F.Body))
	if err != nil {
		return "", fmt.Errorf("error parsing tags template: %w", err)
	}

	// replace tags in template
	s.recursivelyReplaceTags(tagsTemplate, fileMetadata)

	s.recursivelyReplace(
		parsedTemplate,
		parsedMDHTML,
		fileMetadata,
		tagsTemplate,
	)

	contentBuffer := bytes.NewBuffer(make([]byte, 0))
	html.Render(contentBuffer, parsedTemplate)

	content := contentBuffer.String()

	// content := strings.ReplaceAll(st, "{{ $content }}", mdHTML.String())
	// content = strings.ReplaceAll(content, "{{ $title }}", utils.GetSafeValue[string](fileMetadata.Title))

	return content, nil
}

func (s *Server) injectComponents(
	fileMetadata *parser.Meta,
	html *string,
) (string, error) {
	if html == nil {
		return "", fmt.Errorf("html is nil")
	}

	content := *html

	// inject components like <Navbar /> with the content of the component
	stoppedMidway := false
	s.ComponentsMD.Store.Range(func(key string, value *parser.Meta) bool {
		componentContent := value.F.Body

		// TODO: see if component actually exists
		if componentContent == nil {
			return true
		}

		// FIXME: this is a hack to inject the component content
		content = strings.ReplaceAll(
			content,
			fmt.Sprintf("<%s />", key),
			string(componentContent),
		)

		content = strings.ReplaceAll(
			content,
			fmt.Sprintf("<%s/>", key),
			string(componentContent),
		)

		content = strings.ReplaceAll(
			content,
			fmt.Sprintf("<%s><%s/>", key, key),
			string(componentContent),
		)

		content = strings.ReplaceAll(
			content,
			fmt.Sprintf("<%s><%s />", key, key),
			string(componentContent),
		)

		return true
	})

	if stoppedMidway {
		return "", fmt.Errorf("error injecting components")
	}

	return content, nil
}

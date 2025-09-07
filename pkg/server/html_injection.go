package server

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shreyaskaundinya/garlic/pkg/parser"
	"github.com/shreyaskaundinya/garlic/pkg/utils"
)

func (s *Server) injectHTML(
	fileMetadata *parser.Meta,
	html bytes.Buffer,
) (string, error) {
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

	st := string(t.F.Body)

	content := strings.ReplaceAll(st, "{{ $content }}", html.String())
	content = strings.ReplaceAll(content, "{{ $title }}", utils.GetSafeValue[string](fileMetadata.Title))

	return s.injectComponents(fileMetadata, &content)
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

		// if we find <Tags />, we need to inject the tags from the file metadata
		if key == "Tags" {
			tags := fileMetadata.Tags

			tagsString := ""
			for _, tag := range tags {
				tagsString += fmt.Sprintf(
					"<li><a href=\"/tags/%s\">%s</a></li>\n",
					tag, tag,
				)
			}

			componentContent = []byte(strings.ReplaceAll(
				string(componentContent), "{{ $tags }}",
				tagsString,
			))
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

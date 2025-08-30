package parser

import "github.com/shreyaskaundinya/garlic/pkg/utils"

type Frontmatter struct {
	Store map[string]any
}

func NewFrontmatter() *Frontmatter {
	return &Frontmatter{
		Store: make(map[string]any),
	}
}

func (f *Frontmatter) Get(key string) (any, bool) {
	v, ok := f.Store[key]

	if !ok {
		return nil, false
	}

	return v, ok
}

func (f *Frontmatter) Set(key string, value any) {
	f.Store[key] = value
}

func (f *Frontmatter) GetTags() []string {
	log := utils.NewLogger()

	tags, ok := f.Store["tags"]
	if !ok {
		log.Errorw("Tags not found in frontmatter", "frontmatter", *f)
		return []string{}
	}

	tagsArray, ok := tags.([]any)
	if !ok {
		log.Errorw("Tags not found in frontmatter", "frontmatter", *f)
		return []string{}
	}

	tagsStrings := make([]string, len(tagsArray))
	for i, tag := range tagsArray {
		tagsStrings[i] = utils.GetSafeValue[string](tag)
	}

	return tagsStrings
}

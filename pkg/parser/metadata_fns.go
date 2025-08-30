package parser

import "github.com/puzpuzpuz/xsync/v3"

func NewMetadataMap() *Metadata {
	return &Metadata{
		Store: xsync.NewMapOf[string, *Meta](),
	}
}

func (m *Metadata) Set(key string, value *Meta) {
	m.Store.Store(key, value)
}

func (m *Metadata) Get(key string) (*Meta, bool) {
	v, ok := m.Store.Load(key)

	if !ok {
		return nil, false
	}

	return v, true
}

func (m *Metadata) Delete(key string) {
	m.Store.Delete(key)
}

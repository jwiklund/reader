package reader

import (
	"github.com/jwiklund/reader/types"
)

type reader struct {
	store types.Store
	rss   types.Rss
}

func (r reader) GetStore() types.Store {
	return r.store
}

func (r reader) GetRss() types.Rss {
	return r.rss
}

func (r reader) Close() {
	r.rss.Close()
	r.store.Close()
}

func NewReader(dataDir string) types.Reader {
	store := NewStore(dataDir)
	rss := NewRss(store)
	return reader{store, rss}
}

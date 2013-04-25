package reader

import (
	"github.com/jwiklund/reader/rss"
	"github.com/jwiklund/reader/types"
)

func NewRss(s types.Store) types.Rss {
	return rss.NewRss(s)
}

package types

const (
	DateFormat = "2006-01-02 15:04:05"
)

type Item struct {
	Id          string
	Title       string
	Description string
	Published   string
	Content     string
	Type        string
	Url         string
}

type Feed struct {
	Id          string
	Title       string
	Type        string
	Url         string
	LastFetched string
	LastError   string
	Items       []Item
}

type Store interface {
	Get(id string) (*Feed, error)
	GetByUser(user string) ([]Item, error)
	GetByType(feedType string) ([]string, error)
	GetAllInfo() ([]Feed, error)
	Put(feed *Feed) error
	Close()
}

type Rss interface {
	Fetch(id string) error
	FetchAll() []error
	Close()
}

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

type UserFeed struct {
	FeedId     string
	ReadItemId string
}

type UserFeedGroup struct {
	GroupName string
	Feeds     []UserFeed
}

type User struct {
	Id    string
	Feeds []UserFeedGroup
}

type Store interface {
	GetAllFeedsInfo() ([]Feed, error)
	GetFeed(id string) (*Feed, error)
	GetFeedByUser(user, group string) ([]Item, error)
	GetFeedByType(feedType string) ([]string, error)
	PutFeed(feed *Feed) error
	GetUser(user string) (*User, error)
	Close()
}

type Rss interface {
	Fetch(id string) error
	FetchAll() []error
	Close()
}

reader
======

feed
Id, Type, Json { Id, Title, Type, Url, LastFetched, LastError }

item
Id, Json [ {Id, Title, Description, Content, Type, Url }]

user
Id, Json { feed: item, read: id|id }

rss.update -> store
user.update <- store -> store?
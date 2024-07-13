package models

type Document struct {
	URL            string
	PubDate        uint64
	FetchTime      uint64
	Text           string
	FirstFetchTime uint64
}

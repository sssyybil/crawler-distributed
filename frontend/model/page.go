package model

import "crawler-distributed/model"

type SearchResult struct {
	Hits     float64
	Start    int
	Query    string
	PrevFrom int
	NextFrom int
	Items    []*model.Item
}

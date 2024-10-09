package collapse

import "github.com/KingSolvewer/elasticsearch-query-builder/esearch"

type Collapse struct {
	Field                      string `json:"field"`
	esearch.ExpandInnerHits    `json:"inner_hits,omitempty"`
	MaxConcurrentGroupSearches int `json:"max_concurrent_group_searches,omitempty"`
}

type InnerHits struct {
	Name string            `json:"name"`
	Size esearch.Paginator `json:"size"`
	From esearch.Paginator `json:"from,omitempty"`
	Sort esearch.Sort      `json:"sort,omitempty"`
}

type MultiInnerHits []InnerHits

func (hit InnerHits) String() {

}

func (m MultiInnerHits) String() {

}

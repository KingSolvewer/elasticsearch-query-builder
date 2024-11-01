package collapse

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
)

type ParamsFunc func() CollapsedParams

type Collapser struct {
	Field string `json:"field"`
	CollapsedParams
}

func (c Collapser) Collapse() {

}

type CollapsedParams struct {
	InnerHits                  esearch.ExpandInnerHits `json:"inner_hits,omitempty"`
	MaxConcurrentGroupSearches int                     `json:"max_concurrent_group_searches,omitempty"`
}

type InnerHits struct {
	Name string            `json:"name"`
	Size esearch.Paginator `json:"size"`
	From esearch.Paginator `json:"from,omitempty"`
	Sort []esearch.SortMap `json:"sort,omitempty"`
}

type MultiInnerHits []InnerHits

func (hit InnerHits) ExpandHits() {

}

func (m MultiInnerHits) ExpandHits() {

}

package aggs

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/es"
)

type TopHitsFunc func() TopHitsParam

type TermsAggs struct {
	Terms `json:"terms"`
	Aggs  map[string]TopHitsAggs `json:"aggs,omitempty"`
}

type Terms struct {
	Field string `json:"field"`
	TermsParam
}

type TermsParam struct {
	Size  int                     `json:"size,omitempty"`
	Order map[string]es.OrderType `json:"order,omitempty"`
}

func (t TermsAggs) Aggregate(field string) es.Aggregator {
	return t
}

type HistogramAggs struct {
	Histogram `json:"histogram"`
	Aggs      map[string]TopHitsAggs `json:"aggs,omitempty"`
}

type Histogram struct {
	Field string `json:"field"`
	HistogramParam
}

type HistogramParam struct {
	Interval       any                     `json:"interval"`
	MinDocCount    int                     `json:"min_doc_count,omitempty"`
	ExtendedBounds map[string]int          `json:"extended_bounds,omitempty"`
	Order          map[string]es.OrderType `json:"order,omitempty"`
	Offset         int                     `json:"offset,omitempty"`
	Format         string                  `json:"format,omitempty"`
}

func (h HistogramAggs) Aggregate(field string) es.Aggregator {
	return h
}

type RangeAggs struct {
	Range `json:"range"`
	Aggs  map[string]TopHitsAggs `json:"aggs,omitempty"`
}

type Range struct {
	Field string `json:"field"`
	RangeParam
}

type RangeParam struct {
	Keyed  bool     `json:"keyed,omitempty"`
	Ranges []Ranges `json:"ranges"`
	Format string   `json:"format,omitempty"`
}

type Ranges struct {
	From any    `json:"from,omitempty"`
	To   any    `json:"to,omitempty"`
	Key  string `json:"key,omitempty"`
}

func (r RangeAggs) Aggregate(field string) es.Aggregator {
	return r
}

type AvgAggs struct {
	Metric `json:"avg"`
}

type Metric struct {
	Field string `json:"field"`
	MetricParam
}

type MetricParam struct {
	Missing int `json:"missing,omitempty"`
}

func (a MetricParam) Aggregate(field string) es.Aggregator {
	return Metric{
		Field:       field,
		MetricParam: a,
	}
}

type MaxAggs struct {
	Metric `json:"max"`
}
type MinAggs struct {
	Metric `json:"min"`
}
type SumAggs struct {
	Metric `json:"sum"`
}
type StatsAggs struct {
	Metric `json:"stats"`
}
type ExtendedStatsAggs struct {
	Metric `json:"extended_stats"`
}

type CardinalityAggs struct {
	Cardinality `json:"cardinality"`
}

type Cardinality struct {
	Field string
	CardinalityParam
}

type CardinalityParam struct {
	Missing            int `json:"missing,omitempty"`
	PrecisionThreshold int `json:"precision_threshold,omitempty"`
}

func (p CardinalityParam) Aggregate(field string) es.Aggregator {
	return Cardinality{
		Field:            field,
		CardinalityParam: p,
	}
}

type TopHitsAggs struct {
	TopHits `json:"top_hits"`
}

type TopHits struct {
	From   es.Paginator `json:"from,omitempty"`
	Size   es.Paginator `json:"size,omitempty"`
	Sort   []es.Sort    `json:"sort,omitempty"`
	Source []string     `json:"_source,omitempty"`
}

type TopHitsParam struct {
	From   uint                    `json:"from,omitempty"`
	Size   uint                    `json:"size,omitempty"`
	Sort   map[string]es.OrderType `json:"sort,omitempty"`
	Source []string                `json:"_source,omitempty"`
}

func (t TopHits) Aggregate(field string) es.Aggregator {
	return t
}

func (p TopHitsParam) TopHits() TopHits {
	var (
		size es.Uint
		from es.Uint
	)
	if p.Size >= 0 {
		size = es.Uint(p.Size)
	} else {
		size = es.Uint(10)
	}

	if p.From > 0 {
		from = es.Uint(p.From)
	}

	sorts := make([]es.Sort, 0)
	if p.Sort != nil {
		for field, order := range p.Sort {
			sort := make(es.Sort)
			sort[field] = es.Order{
				Order: order,
			}

			sorts = append(sorts, sort)
		}
	}

	topHits := TopHits{
		From:   from,
		Size:   size,
		Sort:   sorts,
		Source: p.Source,
	}

	return topHits
}

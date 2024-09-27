package aggs

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/es"
)

type Aggregator interface {
	Aggregate(args string) Aggregator
}

type Terms struct {
	Field string `json:"field"`
	TermsParam
}

type TermsParam struct {
	Size  int                     `json:"size,omitempty"`
	Order map[string]es.OrderType `json:"order,omitempty"`
}

func (t TermsParam) Aggregate(field string) Aggregator {
	return Terms{
		Field:      field,
		TermsParam: t,
	}
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

func (h HistogramParam) Aggregate(field string) Aggregator {
	return Histogram{
		Field:          field,
		HistogramParam: h,
	}
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

func (r RangeParam) Aggregate(field string) Aggregator {
	return Range{
		Field:      field,
		RangeParam: r,
	}
}

type Avg struct {
	Field string `json:"field"`
	MetricParam
}

type MetricParam struct {
	Missing int `json:"missing,omitempty"`
}

func (a MetricParam) Aggregate(field string) Aggregator {
	return Avg{
		Field:       field,
		MetricParam: a,
	}
}

type Max = Avg
type Min = Avg
type Sum = Avg
type Count = Avg
type Stats = Avg
type ExtendedStats = Avg

type Cardinality struct {
	Field string
	CardinalityParam
}

type CardinalityParam struct {
	Missing            int `json:"missing,omitempty"`
	PrecisionThreshold int `json:"precision_threshold,omitempty"`
}

func (p CardinalityParam) Aggregate(field string) Aggregator {
	return Cardinality{
		Field:            field,
		CardinalityParam: p,
	}
}

package aggs

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
)

type TopHitsFunc func() TopHitsParam

type TermsAggs struct {
	Terms `json:"terms"`
	Aggs  map[string]esearch.Aggregator `json:"aggs,omitempty"`
}

type Terms struct {
	Field string `json:"field"`
	TermsParam
}

type TermsParam struct {
	Size  int             `json:"size,omitempty"`
	Order esearch.SortMap `json:"order,omitempty"`
}

func (agg *TermsAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
	agg.Aggs = subAgg
}

type HistogramAggs struct {
	Histogram `json:"histogram"`
	Aggs      map[string]esearch.Aggregator `json:"aggs,omitempty"`
}

type Histogram struct {
	Field string `json:"field"`
	HistogramParam
}

type HistogramParam struct {
	Interval       any                          `json:"interval"`
	MinDocCount    int                          `json:"min_doc_count,omitempty"`
	ExtendedBounds map[string]any               `json:"extended_bounds,omitempty"`
	Order          map[string]esearch.OrderType `json:"order,omitempty"`
	Offset         any                          `json:"offset,omitempty"`
	Format         string                       `json:"format,omitempty"`
}

func (agg *HistogramAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
	agg.Aggs = subAgg
}

type DateHistogramAggs struct {
	Histogram `json:"date_histogram"`
	Aggs      map[string]esearch.Aggregator `json:"aggs,omitempty"`
}

func (agg *DateHistogramAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
	agg.Aggs = subAgg
}

type RangeAggs struct {
	Range `json:"range"`
	Aggs  map[string]esearch.Aggregator `json:"aggs,omitempty"`
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

func (agg *RangeAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
	agg.Aggs = subAgg
}

type DateRangeAggs struct {
	Range `json:"date_range"`
	Aggs  map[string]esearch.Aggregator `json:"aggs,omitempty"`
}

func (agg *DateRangeAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
	agg.Aggs = subAgg
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

func (metric *AvgAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
}

type MaxAggs struct {
	Metric `json:"max"`
}

func (metric *MaxAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
}

type MinAggs struct {
	Metric `json:"min"`
}

func (metric *MinAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
}

type SumAggs struct {
	Metric `json:"sum"`
}

func (metric *SumAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
}

type StatsAggs struct {
	Metric `json:"stats"`
}

func (metric *StatsAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
}

type ExtendedStatsAggs struct {
	Metric `json:"extended_stats"`
}

func (metric *ExtendedStatsAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
}

type ValueCount struct {
	Metric `json:"value_count"`
}

func (metric *ValueCount) Aggregate(subAgg map[string]esearch.Aggregator) {
}

type CardinalityAggs struct {
	Cardinality `json:"cardinality"`
}

type Cardinality struct {
	Field string `json:"field"`
	CardinalityParam
}

type CardinalityParam struct {
	Missing            int `json:"missing,omitempty"`
	PrecisionThreshold int `json:"precision_threshold,omitempty"`
}

type CardinalityFunc func() CardinalityParam

func (metric *CardinalityAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
}

type TopHitsAggs struct {
	TopHits `json:"top_hits"`
}

func (metric *TopHitsAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
}

type TopHits struct {
	From   esearch.Paginator `json:"from,omitempty"`
	Size   esearch.Paginator `json:"size,omitempty"`
	Sort   []esearch.Sorter  `json:"sort,omitempty"`
	Source []string          `json:"_source,omitempty"`
}

type TopHitsParam struct {
	From   uint            `json:"from,omitempty"`
	Size   uint            `json:"size,omitempty"`
	Sort   esearch.SortMap `json:"sort,omitempty"`
	Source []string        `json:"_source,omitempty"`
}

func (p TopHitsParam) TopHitsAgg() *TopHitsAggs {
	var (
		size esearch.Uint
		from esearch.Uint
	)
	if p.Size >= 0 {
		size = esearch.Uint(p.Size)
	} else {
		size = esearch.Uint(10)
	}

	if p.From > 0 {
		from = esearch.Uint(p.From)
	}

	sorts := make([]esearch.Sorter, 0)
	if p.Sort != nil {
		for field, order := range p.Sort {
			sort := make(esearch.Sort)
			sort[field] = esearch.Order{
				Order: order,
			}

			sorts = append(sorts, sort)
		}
	}

	topHitsAgg := &TopHitsAggs{
		TopHits: TopHits{
			From:   from,
			Size:   size,
			Sort:   sorts,
			Source: p.Source,
		},
	}

	return topHitsAgg
}

type FilterAggs struct {
	Filter esearch.Query                 `json:"filter"`
	Aggs   map[string]esearch.Aggregator `json:"aggs,omitempty"`
}

func (agg *FilterAggs) Aggregate(subAgg map[string]esearch.Aggregator) {
	agg.Aggs = subAgg
}

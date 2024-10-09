package elastic

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/aggs"
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
)

func (b *Builder) GroupBy(field string, param aggs.TermsParam, hitsFunc aggs.TopHitsFunc) *Builder {
	terms := aggs.TermsAggs{
		Terms: aggs.Terms{
			Field:      field,
			TermsParam: param,
		},
		Aggs: make(map[string]aggs.TopHitsAggs),
	}

	if hitsFunc != nil {
		terms.Aggs[field] = aggs.TopHitsAggs{TopHits: hitsFunc().TopHits()}
	}

	return b.Aggs(field+"_"+esearch.Terms, terms)
}

func (b *Builder) Histogram(field string, param aggs.HistogramParam, hitsFunc aggs.TopHitsFunc) *Builder {
	histogram := aggs.HistogramAggs{
		Histogram: aggs.Histogram{
			Field:          field,
			HistogramParam: param,
		},
		Aggs: make(map[string]aggs.TopHitsAggs),
	}

	if hitsFunc != nil {
		histogram.Aggs[field] = aggs.TopHitsAggs{TopHits: hitsFunc().TopHits()}
	}

	return b.Aggs(field+"_"+esearch.Histogram, histogram)
}

func (b *Builder) DateGroupBy(field string, param aggs.HistogramParam, hitsFunc aggs.TopHitsFunc) *Builder {
	return b.Histogram(field, param, hitsFunc)
}

func (b *Builder) Range(field string, param aggs.RangeParam, hitsFunc aggs.TopHitsFunc) *Builder {
	if !checkAggsRangeType(param.Ranges) {
		return b
	}

	rangeAggs := aggs.RangeAggs{
		Range: aggs.Range{
			Field: field,
		},
		Aggs: make(map[string]aggs.TopHitsAggs),
	}

	if hitsFunc != nil {
		rangeAggs.Aggs[field] = aggs.TopHitsAggs{TopHits: hitsFunc().TopHits()}
	}

	return b.Aggs(field+"_"+esearch.Range, rangeAggs)
}

func (b *Builder) DateRange(field string, param aggs.RangeParam, hitsFunc aggs.TopHitsFunc) *Builder {
	return b.Range(field, param, hitsFunc)
}

func (b *Builder) Aggs(aggsField string, aggregator esearch.Aggregator) *Builder {
	b.aggs[aggsField] = aggregator
	return b
}

func (b *Builder) Avg(field string, param aggs.MetricParam) *Builder {
	avgAggs := aggs.AvgAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+esearch.Avg, avgAggs)
}

func (b *Builder) Max(field string, param aggs.MetricParam) *Builder {
	maxAggs := aggs.MaxAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}

	return b.Aggs(field+"_"+esearch.Max, maxAggs)
}

func (b *Builder) Min(field string, param aggs.MetricParam) *Builder {
	minAggs := aggs.MinAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+esearch.Min, minAggs)
}

func (b *Builder) Sum(field string, param aggs.MetricParam) *Builder {
	sumAggs := aggs.SumAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+esearch.Sum, sumAggs)
}

func (b *Builder) Stats(field string, param aggs.MetricParam) *Builder {
	statsAggs := aggs.StatsAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+esearch.Stats, statsAggs)
}

func (b *Builder) ExtendedStats(field string, param aggs.MetricParam) *Builder {
	statsAggs := aggs.ExtendedStatsAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+esearch.ExtendedStats, statsAggs)
}

func (b *Builder) TopHits(hits aggs.TopHitsParam) *Builder {
	hitsAggs := hits.TopHits()

	return b.Aggs(esearch.TopHits, hitsAggs)
}

func (b *Builder) TopHitsFunc(fn NestedFunc) *Builder {
	newBuilder := NewBuilder()
	fn(newBuilder)
	hits := newBuilder.topHits()
	return b.Aggs(esearch.TopHits, hits)
}

func (b *Builder) topHits() aggs.TopHits {
	var (
		size esearch.Uint
		from esearch.Uint
	)
	if b.size >= 0 {
		size = esearch.Uint(b.size)
	} else {
		size = esearch.Uint(10)
	}

	if b.from > 0 {
		from = esearch.Uint(b.from)
	}

	topHits := aggs.TopHits{
		From:   from,
		Size:   size,
		Sort:   b.sort,
		Source: b.fields,
	}

	return topHits
}

func (b *Builder) Cardinality(field string, fn aggs.CardinalityFunc) *Builder {
	cardinality := aggs.CardinalityAggs{
		Cardinality: aggs.Cardinality{
			Field: field,
		},
	}
	if fn != nil {
		cardinality.Cardinality.CardinalityParam = fn()
	}

	b.Aggs(field+"_"+esearch.Cardinality, cardinality)
	return b
}

func checkAggsRangeType(ranges []aggs.Ranges) bool {
	for _, r := range ranges {
		switch r.From.(type) {
		case int, string:
		default:
			return false
		}
		switch r.To.(type) {
		case int, string:
		default:
			return false
		}
	}

	return true
}

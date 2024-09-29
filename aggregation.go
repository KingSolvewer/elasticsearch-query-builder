package elastic

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/aggs"
	"github.com/KingSolvewer/elasticsearch-query-builder/es"
)

func GroupBy(field string, param aggs.TermsParam, hitsFunc aggs.TopHitsFunc) *Builder {
	return builder.GroupBy(field, param, hitsFunc)
}

func (b *Builder) GroupBy(field string, param aggs.TermsParam, hitsFunc aggs.TopHitsFunc) *Builder {
	terms := aggs.TermsAggs{
		Terms: aggs.Terms{
			Field:      field,
			TermsParam: param,
		},
		Aggs: make(map[string]aggs.TopHitsAggs),
	}

	if hitsFunc != nil {
		terms.Aggs[field] = aggs.TopHitsAggs{TopHits: hitsFunc()}
	}

	return b.Aggs(field+"_"+es.Terms, terms)
}

func Histogram(field string, param aggs.HistogramParam, hitsFunc aggs.TopHitsFunc) *Builder {
	return builder.Histogram(field, param, hitsFunc)
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
		histogram.Aggs[field] = aggs.TopHitsAggs{TopHits: hitsFunc()}
	}

	return b.Aggs(field+"_"+es.Histogram, histogram)
}

func DateGroupBy(field string, param aggs.HistogramParam, hitsFunc aggs.TopHitsFunc) *Builder {
	return builder.DateGroupBy(field, param, hitsFunc)
}

func (b *Builder) DateGroupBy(field string, param aggs.HistogramParam, hitsFunc aggs.TopHitsFunc) *Builder {
	return b.Histogram(field, param, hitsFunc)
}

func Range(field string, param aggs.RangeParam, hitsFunc aggs.TopHitsFunc) *Builder {
	return builder.Range(field, param, hitsFunc)
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
		rangeAggs.Aggs[field] = aggs.TopHitsAggs{TopHits: hitsFunc()}
	}

	return b.Aggs(field+"_"+es.Range, rangeAggs)
}

func DateRange(field string, param aggs.RangeParam, hitsFunc aggs.TopHitsFunc) *Builder {
	return builder.DateRange(field, param, hitsFunc)
}

func (b *Builder) DateRange(field string, param aggs.RangeParam, hitsFunc aggs.TopHitsFunc) *Builder {
	return b.Range(field, param, hitsFunc)
}

func Aggs(aggsField string, aggregator es.Aggregator) *Builder {
	return builder.Aggs(aggsField, aggregator)
}

func (b *Builder) Aggs(aggsField string, aggregator es.Aggregator) *Builder {
	b.aggs[aggsField] = aggregator
	return b
}

func Avg(field string, param aggs.MetricParam) *Builder {
	return builder.Avg(field, param)
}

func (b *Builder) Avg(field string, param aggs.MetricParam) *Builder {
	avgAggs := aggs.AvgAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+es.Avg, avgAggs)
}

func Max(field string, param aggs.MetricParam) *Builder {
	return builder.Max(field, param)
}

func (b *Builder) Max(field string, param aggs.MetricParam) *Builder {
	maxAggs := aggs.MaxAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}

	return b.Aggs(field+"_"+es.Max, maxAggs)
}

func Min(field string, param aggs.MetricParam) *Builder {
	return builder.Max(field, param)
}

func (b *Builder) Min(field string, param aggs.MetricParam) *Builder {
	minAggs := aggs.MinAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+es.Min, minAggs)
}

func Sum(field string, param aggs.MetricParam) *Builder {
	return builder.Sum(field, param)
}

func (b *Builder) Sum(field string, param aggs.MetricParam) *Builder {
	sumAggs := aggs.SumAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+es.Sum, sumAggs)
}

func Stats(field string, param aggs.MetricParam) *Builder {
	return builder.Stats(field, param)
}

func (b *Builder) Stats(field string, param aggs.MetricParam) *Builder {
	statsAggs := aggs.StatsAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+es.Stats, statsAggs)
}

func ExtendedStats(field string, param aggs.MetricParam) *Builder {
	return builder.ExtendedStats(field, param)
}

func (b *Builder) ExtendedStats(field string, param aggs.MetricParam) *Builder {
	statsAggs := aggs.ExtendedStatsAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+es.ExtendedStats, statsAggs)
}

func TopHits(hits aggs.TopHits) *Builder {
	return builder.TopHits(hits)
}

func (b *Builder) TopHits(hits aggs.TopHits) *Builder {
	return b.Aggs(es.TopHits, hits)
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

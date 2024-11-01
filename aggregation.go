package elastic

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/aggs"
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
)

type SubAggFunc func(b *Builder)

type Aggregation struct {
	Params  esearch.Aggregator
	SubAggs []SubAggFunc
}

func (b *Builder) GroupBy(field string, param aggs.TermsParam, subAggFuncSet ...SubAggFunc) *Builder {
	termsAgg := &aggs.TermsAggs{
		Terms: aggs.Terms{
			Field:      field,
			TermsParam: param,
		},
	}

	return b.Aggs(field+"_"+esearch.Terms, termsAgg, subAggFuncSet...)
}

func (b *Builder) Histogram(field string, param aggs.HistogramParam, subAggFuncSet ...SubAggFunc) *Builder {
	histogram := &aggs.HistogramAggs{
		Histogram: aggs.Histogram{
			Field:          field,
			HistogramParam: param,
		},
	}

	return b.Aggs(field+"_"+esearch.Histogram, histogram, subAggFuncSet...)
}

func (b *Builder) DateGroupBy(field string, param aggs.HistogramParam, subAggFuncSet ...SubAggFunc) *Builder {
	histogram := &aggs.DateHistogramAggs{
		Histogram: aggs.Histogram{
			Field:          field,
			HistogramParam: param,
		},
	}

	return b.Aggs(field+"_"+esearch.DateHistogram, histogram, subAggFuncSet...)
}

func (b *Builder) Range(field string, param aggs.RangeParam, subAggFuncSet ...SubAggFunc) *Builder {
	if !checkAggsRangeType(param.Ranges) {
		return b
	}

	rangeAggs := &aggs.RangeAggs{
		Range: aggs.Range{
			Field: field,
		},
	}

	return b.Aggs(field+"_"+esearch.Range, rangeAggs, subAggFuncSet...)
}

func (b *Builder) DateRange(field string, param aggs.RangeParam, subAggFuncSet ...SubAggFunc) *Builder {
	if !checkAggsRangeType(param.Ranges) {
		return b
	}

	rangeAggs := &aggs.DateRangeAggs{
		Range: aggs.Range{
			Field: field,
		},
	}

	return b.Aggs(field+"_"+esearch.DateRange, rangeAggs, subAggFuncSet...)
}

func (b *Builder) AggsFilter(field string, fn NestWhereFunc, subAggFuncs ...SubAggFunc) *Builder {

	if fn != nil {
		newBuilder := NewBuilder()
		fn(newBuilder)

		query := make(esearch.Query)
		query["bool"] = newBuilder.componentWhere()

		filterAggs := &aggs.FilterAggs{
			Filter: query,
		}

		return b.Aggs(field+"_"+esearch.AggsFilter, filterAggs, subAggFuncs...)
	}

	return b
}

func (b *Builder) Aggs(aggField string, aggregator esearch.Aggregator, subAggFuncSet ...SubAggFunc) *Builder {
	agg := &Aggregation{
		Params:  aggregator,
		SubAggs: subAggFuncSet,
	}

	b.aggregations[aggField] = agg
	return b
}

func (b *Builder) Avg(field string, param aggs.MetricParam) *Builder {
	avgAggs := &aggs.AvgAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+esearch.Avg, avgAggs)
}

func (b *Builder) Max(field string, param aggs.MetricParam) *Builder {
	maxAggs := &aggs.MaxAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}

	return b.Aggs(field+"_"+esearch.Max, maxAggs)
}

func (b *Builder) Min(field string, param aggs.MetricParam) *Builder {
	minAggs := &aggs.MinAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+esearch.Min, minAggs)
}

func (b *Builder) Sum(field string, param aggs.MetricParam) *Builder {
	sumAggs := &aggs.SumAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+esearch.Sum, sumAggs)
}

func (b *Builder) Stats(field string, param aggs.MetricParam) *Builder {
	statsAggs := &aggs.StatsAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+esearch.Stats, statsAggs)
}

func (b *Builder) ExtendedStats(field string, param aggs.MetricParam) *Builder {
	statsAggs := &aggs.ExtendedStatsAggs{
		Metric: aggs.Metric{
			Field:       field,
			MetricParam: param,
		},
	}
	return b.Aggs(field+"_"+esearch.ExtendedStats, statsAggs)
}

func (b *Builder) ValueCount(field string) *Builder {
	statsAggs := &aggs.ValueCount{
		Metric: aggs.Metric{
			Field: field,
		},
	}
	return b.Aggs(field+"_"+esearch.ValueCount, statsAggs)
}

func (b *Builder) Cardinality(field string, fn aggs.CardinalityFunc) *Builder {
	cardinality := &aggs.CardinalityAggs{
		Cardinality: aggs.Cardinality{
			Field: field,
		},
	}
	if fn != nil {
		cardinality.Cardinality.CardinalityParam = fn()
	}

	return b.Aggs(field+"_"+esearch.Cardinality, cardinality)
}

func (b *Builder) TopHits(hits aggs.TopHitsParam) *Builder {
	hitsAggs := hits.TopHitsAgg()

	return b.Aggs(esearch.TopHits, hitsAggs)
}

func (b *Builder) TopHitsFunc(fn NestWhereFunc) *Builder {
	newBuilder := NewBuilder()
	fn(newBuilder)
	hits := newBuilder.topHitsAgg()
	return b.Aggs(esearch.TopHits, hits)
}

func (b *Builder) topHitsAgg() *aggs.TopHitsAggs {
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

	topHitsAgg := &aggs.TopHitsAggs{
		TopHits: aggs.TopHits{
			From:   from,
			Size:   size,
			Sort:   b.sort,
			Source: b.fields,
		},
	}

	return topHitsAgg
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

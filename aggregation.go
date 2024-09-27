package elastic

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/aggs"
	"github.com/KingSolvewer/elasticsearch-query-builder/es"
)

func GroupBy(field string, param aggs.TermsParam) *Builder {
	return builder.GroupBy(field, param)
}

func (b *Builder) GroupBy(field string, param aggs.TermsParam) *Builder {
	terms := param.Aggregate(field)

	b.Aggs(field+"_"+es.Terms, es.Terms, terms)

	return b
}

func Histogram(field string, param aggs.HistogramParam) *Builder {
	return builder.Histogram(field, param)
}

func (b *Builder) Histogram(field string, param aggs.HistogramParam) *Builder {
	histogram := param.Aggregate(field)

	b.Aggs(field+"_"+es.Histogram, es.Histogram, histogram)

	return b
}

func DateGroupBy(field string, param aggs.HistogramParam) *Builder {
	return builder.DateGroupBy(field, param)
}

func (b *Builder) DateGroupBy(field string, param aggs.HistogramParam) *Builder {
	dateHistogram := param.Aggregate(field)

	b.Aggs(field+"_"+es.DateHistogram, es.DateHistogram, dateHistogram)

	return b
}

func DateRange(field string, param aggs.RangeParam) *Builder {
	return builder.DateRange(field, param)
}

func (b *Builder) DateRange(field string, param aggs.RangeParam) *Builder {
	if !checkAggsRangeType(param.Ranges) {
		return b
	}

	dateRange := param.Aggregate(field)

	b.Aggs(field+"_"+es.Range, es.Range, dateRange)

	return b
}

func Range(field string, param aggs.RangeParam) *Builder {
	return builder.Range(field, param)
}

func (b *Builder) Range(field string, param aggs.RangeParam) *Builder {
	if !checkAggsRangeType(param.Ranges) {
		return b
	}

	rangeAggs := param.Aggregate(field)

	b.Aggs(field+"_"+es.Range, es.Range, rangeAggs)

	return b
}

func Aggs(aggsField string, bucketType string, aggregation aggs.Aggregator) *Builder {
	return builder.Aggs(aggsField, bucketType, aggregation)
}

func (b *Builder) Aggs(aggsField string, bucketType string, aggregator aggs.Aggregator) *Builder {
	b.aggs[aggsField] = make(map[string]aggs.Aggregator)
	b.aggs[aggsField][bucketType] = aggregator
	return b
}

func Avg(field string, param aggs.MetricParam) *Builder {
	return builder.Avg(field, param)
}

func (b *Builder) Avg(field string, param aggs.MetricParam) *Builder {
	avg := param.Aggregate(field)
	b.Aggs(field+"_"+es.Avg, es.Avg, avg)
	return b
}

func Max(field string, param aggs.MetricParam) *Builder {
	return builder.Max(field, param)
}

func (b *Builder) Max(field string, param aggs.MetricParam) *Builder {
	avg := param.Aggregate(field)
	b.Aggs(field+"_"+es.Max, es.Max, avg)
	return b
}

func Min(field string, param aggs.MetricParam) *Builder {
	return builder.Max(field, param)
}

func (b *Builder) Min(field string, param aggs.MetricParam) *Builder {
	avg := param.Aggregate(field)
	b.Aggs(field+"_"+es.Min, es.Min, avg)
	return b
}

func Sum(field string, param aggs.MetricParam) *Builder {
	return builder.Sum(field, param)
}

func (b *Builder) Sum(field string, param aggs.MetricParam) *Builder {
	avg := param.Aggregate(field)
	b.Aggs(field+"_"+es.Sum, es.Sum, avg)
	return b
}

func Count(field string, param aggs.MetricParam) *Builder {
	return builder.Count(field, param)
}

func (b *Builder) Count(field string, param aggs.MetricParam) *Builder {
	avg := param.Aggregate(field)
	b.Aggs(field+"_"+es.Count, es.Count, avg)
	return b
}

func Stats(field string, param aggs.MetricParam) *Builder {
	return builder.Stats(field, param)
}

func (b *Builder) Stats(field string, param aggs.MetricParam) *Builder {
	avg := param.Aggregate(field)
	b.Aggs(field+"_"+es.Stats, es.Stats, avg)
	return b
}

func ExtendedStats(field string, param aggs.MetricParam) *Builder {
	return builder.ExtendedStats(field, param)
}

func (b *Builder) ExtendedStats(field string, param aggs.MetricParam) *Builder {
	avg := param.Aggregate(field)
	b.Aggs(field+"_"+es.ExtendedStats, es.ExtendedStats, avg)
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

package elastic

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
)

func (b *Builder) compile() *esearch.ElasticQuery {
	query := &esearch.ElasticQuery{
		Query:      make(esearch.Query),
		PostFilter: make(esearch.Query),
	}

	if b.fields != nil || len(b.fields) > 0 {
		query.Source = b.fields
	}
	if b.sort != nil || len(b.sort) > 0 {
		query.Sort = b.sort
	}

	if b.size >= 0 {
		query.Size = esearch.Uint(b.size)
	} else {
		query.Size = esearch.Uint(10)
	}

	if b.from > 0 {
		query.From = esearch.Uint(b.from)
	}

	if b.collapse != nil {
		query.Collapse = b.collapse
	}

	if b.where != nil {
		boolQuery := b.componentWhere()

		if len(boolQuery.Should) > 0 {
			boolQuery.MinimumShouldMatch = b.minimumShouldMatch
		}

		query.Query["bool"] = boolQuery
	}

	if b.postWhere != nil {
		newBuilder := NewBuilder()
		b.postWhere(newBuilder)
		postQuery := newBuilder.componentWhere()
		query.PostFilter["bool"] = postQuery
	}

	if b.aggregations != nil {
		aggSet := make(map[string]esearch.Aggregator)
		b.componentAggs(aggSet)
		query.Aggs = aggSet
	}

	return query
}

func (b *Builder) componentWhere() *esearch.BoolQuery {
	boolQuery := &esearch.BoolQuery{}

	for key, items := range b.where {
		switch key {
		case esearch.Must:
			boolQuery.Must = append(boolQuery.Must, items...)
		case esearch.MustNot:
			boolQuery.MustNot = append(boolQuery.MustNot, items...)
		case esearch.Should:
			boolQuery.Should = append(boolQuery.Should, items...)
		case esearch.FilterClause:
			boolQuery.Filter = append(boolQuery.Filter, items...)
		}
	}

	for key, fns := range b.nested {
		for _, fn := range fns {
			newBuilder := NewBuilder()
			fn(newBuilder)
			newBoolQuery := newBuilder.componentWhere()

			newQuery := make(esearch.Query)
			newQuery["bool"] = newBoolQuery

			switch key {
			case esearch.Must:
				boolQuery.Must = append(boolQuery.Must, newQuery)
			case esearch.MustNot:
				boolQuery.MustNot = append(boolQuery.MustNot, newQuery)
			case esearch.Should:
				boolQuery.Should = append(boolQuery.Should, newQuery)
			case esearch.FilterClause:
				boolQuery.Filter = append(boolQuery.Filter, newQuery)
			}
		}
	}

	return boolQuery
}

func (b *Builder) componentAggs(aggSet map[string]esearch.Aggregator) {
	for alias, aggregation := range b.aggregations {
		aggregation.subAggs()
		aggSet[alias] = aggregation.Params
	}
}

func (aggregation *Aggregation) subAggs() {
	if aggregation.SubAggs != nil {
		newAggSet := make(map[string]esearch.Aggregator)
		for _, subAggFunc := range aggregation.SubAggs {
			newBuilder := NewBuilder()
			subAggFunc(newBuilder)
			newBuilder.componentAggs(newAggSet)
		}
		aggregation.Params.Aggregate(newAggSet)
	}
}

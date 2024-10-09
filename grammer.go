package elastic

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
)

func (b *Builder) compile() {

	b.query = &esearch.ElasticQuery{
		Source: make([]string, 0),
		Sort:   make([]esearch.Sort, 0),
		Query:  make(map[string]esearch.QueryBuilder),
	}

	if b.fields != nil || len(b.fields) > 0 {
		b.query.Source = b.fields
	}
	if b.sort != nil || len(b.sort) > 0 {
		b.query.Sort = b.sort
	}

	if b.size >= 0 {
		b.query.Size = esearch.Uint(b.size)
	} else {
		b.query.Size = esearch.Uint(10)
	}

	if b.from > 0 {
		b.query.From = esearch.Uint(b.from)
	}

	if b.collapse.Field != "" {
		b.query.Collapse = b.collapse
	}

	boolQuery := b.component()

	if len(boolQuery.Should) > 0 {
		boolQuery.MinimumShouldMatch = b.minimumShouldMatch
	}

	b.query.Query["bool"] = boolQuery
	b.query.Aggs = b.aggs
}

func (b *Builder) component() esearch.BoolQuery {
	boolQuery := esearch.BoolQuery{}

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
			newBoolQuery := newBuilder.component()

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
